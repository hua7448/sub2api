package service

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
)

// marshalCodexASCIIJSON serializes a value as compact JSON and escapes every
// non-ASCII code point in string content.  Codex carries
// x-codex-turn-metadata through HTTP headers as well as client_metadata; an
// ordinary encoding/json marshal leaves most UTF-8 bytes unescaped, which can
// make an otherwise valid metadata blob fail header validation.  The wire
// representation remains JSON (rather than base64), matching Codex's
// to_ascii_json_string helper.
//
// SetEscapeHTML(false) keeps the serializer's output close to serde_json's
// default.  The second pass only changes Unicode inside JSON strings and
// leaves JSON punctuation and existing escapes untouched.
func marshalCodexASCIIJSON(value any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	raw := buf.Bytes()
	if len(raw) > 0 && raw[len(raw)-1] == '\n' {
		raw = raw[:len(raw)-1]
	}
	return escapeCodexJSONNonASCII(raw), nil
}

const codexJSONHex = "0123456789abcdef"

func appendCodexJSONUnicodeEscape(dst []byte, codeUnit uint16) []byte {
	dst = append(dst, '\\', 'u')
	dst = append(dst,
		codexJSONHex[(codeUnit>>12)&0x0f],
		codexJSONHex[(codeUnit>>8)&0x0f],
		codexJSONHex[(codeUnit>>4)&0x0f],
		codexJSONHex[codeUnit&0x0f],
	)
	return dst
}

func appendCodexJSONRune(dst []byte, r rune) []byte {
	if r <= 0xffff {
		return appendCodexJSONUnicodeEscape(dst, uint16(r))
	}
	// Encode the supplementary-plane rune as the UTF-16 surrogate pair used
	// by the official Codex serializer (for example 😀 → \ud83d\ude00).
	r -= 0x10000
	high := uint16(0xd800 + (r >> 10))
	low := uint16(0xdc00 + (r & 0x3ff))
	dst = appendCodexJSONUnicodeEscape(dst, high)
	return appendCodexJSONUnicodeEscape(dst, low)
}

func escapeCodexJSONNonASCII(raw []byte) []byte {
	// Most metadata is ASCII and should not incur a second allocation when no
	// Unicode is present.  We still scan strings so existing JSON escapes remain
	// byte-for-byte unchanged.
	needsEscape := false
	for i := 0; i < len(raw); i++ {
		if raw[i] >= utf8.RuneSelf {
			needsEscape = true
			break
		}
	}
	if !needsEscape {
		return raw
	}

	out := make([]byte, 0, len(raw)+16)
	inString := false
	escaped := false
	for i := 0; i < len(raw); {
		b := raw[i]
		if !inString {
			out = append(out, b)
			if b == '"' {
				inString = true
			}
			i++
			continue
		}

		if escaped {
			out = append(out, b)
			escaped = false
			i++
			continue
		}
		switch b {
		case '\\':
			out = append(out, b)
			escaped = true
			i++
		case '"':
			out = append(out, b)
			inString = false
			i++
		default:
			if b < utf8.RuneSelf {
				out = append(out, b)
				i++
				continue
			}
			r, size := utf8.DecodeRune(raw[i:])
			if r == utf8.RuneError && size == 1 {
				// encoding/json replaces invalid UTF-8 before this helper is
				// reached.  Keep the output ASCII even for a custom raw value
				// that slips an invalid byte through.
				out = appendCodexJSONUnicodeEscape(out, 0xfffd)
				i++
				continue
			}
			out = appendCodexJSONRune(out, r)
			i += size
		}
	}
	return out
}
