package media

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadSmallImage(t *testing.T) {
	// Create a temporary test image
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jpg")
	testData := []byte("fake image data")
	require.NoError(t, os.WriteFile(testFile, testData, 0644))

	initCalled := false
	appendCalled := false
	finalizeCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Parse form to get command from either query params or form data
		r.ParseMultipartForm(10 << 20) // 10MB
		command := r.URL.Query().Get("command")
		if command == "" && r.MultipartForm != nil {
			if values := r.MultipartForm.Value["command"]; len(values) > 0 {
				command = values[0]
			}
		}

		switch command {
		case "INIT":
			initCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"media_id":123456,"media_id_string":"123456"}`))
		case "APPEND":
			appendCalled = true
			w.WriteHeader(http.StatusNoContent)
		case "FINALIZE":
			finalizeCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"media_id":123456,"media_id_string":"123456","size":15}`))
		default:
			t.Logf("unexpected command: %s", command)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	service := NewService("test-token")
	service.uploadBaseURL = server.URL + "/"

	mediaID, err := service.UploadFile(context.Background(), testFile, CategoryTweetImage)

	require.NoError(t, err)
	require.Equal(t, "123456", mediaID)
	require.True(t, initCalled, "INIT should be called")
	require.True(t, appendCalled, "APPEND should be called")
	require.True(t, finalizeCalled, "FINALIZE should be called")
}

func TestDetectMediaType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.png", "image/png"},
		{"test.gif", "image/gif"},
		{"test.webp", "image/webp"},
		{"test.mp4", "video/mp4"},
		{"test.mov", "video/quicktime"},
		{"test.unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := detectMediaType(tt.path)
			require.Equal(t, tt.expected, result)
		})
	}
}
