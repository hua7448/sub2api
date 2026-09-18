package service

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewriteOpenAIImagesMultipartRequest_GPTImage2PreservesExactSize(t *testing.T) {
	var input bytes.Buffer
	writer := multipart.NewWriter(&input)
	require.NoError(t, writer.WriteField("model", "gpt-image-2"))
	require.NoError(t, writer.WriteField("prompt", "edit"))
	require.NoError(t, writer.WriteField("size", "3840x2160"))
	require.NoError(t, writer.WriteField("quality", "hd"))
	require.NoError(t, writer.Close())

	rewritten, contentType, err := rewriteOpenAIImagesRequest(
		input.Bytes(),
		writer.FormDataContentType(),
		"gpt-image-2",
		"3840x2160",
	)
	require.NoError(t, err)

	_, params, err := mime.ParseMediaType(contentType)
	require.NoError(t, err)
	reader := multipart.NewReader(bytes.NewReader(rewritten), params["boundary"])
	fields := make(map[string]string)
	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		require.NoError(t, nextErr)
		value, readErr := io.ReadAll(part)
		require.NoError(t, readErr)
		fields[part.FormName()] = string(value)
	}

	require.Equal(t, "gpt-image-2", fields["model"])
	require.Equal(t, "3840x2160", fields["size"])
	require.Equal(t, "hd", fields["quality"])
}
