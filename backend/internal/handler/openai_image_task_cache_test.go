package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioImageTaskReplayCompletedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	scope := "test:replay"
	requestHash := "same-request"
	task, owner, conflict := acquireStudioImageTask(scope, requestHash)
	require.True(t, owner)
	require.False(t, conflict)
	defer studioImageTasks.Delete(scope)

	ownerRecorder := httptest.NewRecorder()
	ownerCtx, _ := gin.CreateTestContext(ownerRecorder)
	ownerWriter := &studioImageTaskCaptureWriter{ResponseWriter: ownerCtx.Writer}
	ownerWriter.Header().Set("Content-Type", "application/json")
	ownerWriter.WriteHeader(http.StatusCreated)
	_, err := ownerWriter.Write([]byte(`{"ok":true}`))
	require.NoError(t, err)
	task.complete(ownerWriter)

	replayRecorder := httptest.NewRecorder()
	replayCtx, _ := gin.CreateTestContext(replayRecorder)
	replayCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
	serveStudioImageTask(replayCtx, task)

	require.Equal(t, http.StatusCreated, replayRecorder.Code)
	require.JSONEq(t, `{"ok":true}`, replayRecorder.Body.String())
	require.Equal(t, "replayed", replayRecorder.Header().Get("X-Sub2API-Image-Task-Status"))
}

func TestStudioImageTaskRejectsConflictingRequestHash(t *testing.T) {
	scope := "test:conflict"
	task, owner, conflict := acquireStudioImageTask(scope, "hash-a")
	require.NotNil(t, task)
	require.True(t, owner)
	require.False(t, conflict)
	defer studioImageTasks.Delete(scope)

	existing, owner, conflict := acquireStudioImageTask(scope, "hash-b")
	require.Same(t, task, existing)
	require.False(t, owner)
	require.True(t, conflict)
}

func TestNormalizeStudioImageTaskKey(t *testing.T) {
	require.Equal(t, "conv_1:turn_2", normalizeStudioImageTaskKey(" conv/1:turn 2 "))
	require.Empty(t, normalizeStudioImageTaskKey(" \t\n "))
}

func TestStudioImageTaskHashIgnoresMultipartBoundary(t *testing.T) {
	a := studioImageTaskRequestHash("/v1/images/edits", "multipart/form-data; boundary=a", []byte("body-a"))
	b := studioImageTaskRequestHash("/v1/images/edits", "multipart/form-data; boundary=b", []byte("body-b"))
	require.Equal(t, a, b)

	jsonA := studioImageTaskRequestHash("/v1/images/generations", "application/json", []byte(`{"prompt":"a"}`))
	jsonB := studioImageTaskRequestHash("/v1/images/generations", "application/json", []byte(`{"prompt":"b"}`))
	require.NotEqual(t, jsonA, jsonB)
}
