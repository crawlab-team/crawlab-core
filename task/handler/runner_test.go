package handler

import (
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRunner struct {
	mock.Mock
	Runner
}

func (m *MockRunner) downloadFile(url string, filePath string) error {
	args := m.Called(url, filePath)
	return args.Error(0)
}

func newMockRunner() *MockRunner {
	r := &MockRunner{}
	r.s = &models.Spider{}
	workspacePath := viper.GetString("workspace")
	_ = os.MkdirAll(filepath.Join(workspacePath, r.s.GetId().Hex()), os.ModePerm)
	return r
}

func TestSyncFiles_SuccessWithDummyFiles(t *testing.T) {
	// Create a test server that responds with a list of files
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"file1.txt":{"path": "file1.txt", "hash": "hash1"}, "file2.txt":{"path": "file2.txt", "hash": "hash2"}}`))
	}))
	defer ts.Close()

	// Create a mock runner
	r := newMockRunner()
	r.On("downloadFile", mock.Anything, mock.Anything).Return(nil)

	// Create dummy files
	workspacePath := viper.GetString("workspace")
	os.Create(filepath.Join(workspacePath, r.s.GetId().Hex(), "file1.txt"))
	os.Create(filepath.Join(workspacePath, r.s.GetId().Hex(), "file2.txt"))

	// Set the master URL to the test server URL
	viper.Set("api.endpoint", ts.URL)

	localPath := os.TempDir()
	os.MkdirAll(filepath.Join(localPath, r.s.GetId().Hex()), os.ModePerm)
	defer os.RemoveAll(localPath)

	// Call the method under test
	err := r.syncFiles()

	// Assert that there was no error and the downloadFile method was called twice
	assert.NoError(t, err)
	r.AssertNumberOfCalls(t, "downloadFile", 2)

	assert.FileExists(t, filepath.Join(localPath, r.s.GetId().Hex(), "file1.txt"))
	assert.FileExists(t, filepath.Join(localPath, r.s.GetId().Hex(), "file2.txt"))

	// Clean up dummy files
	os.Remove(filepath.Join(workspacePath, r.s.GetId().Hex(), "file1.txt"))
	os.Remove(filepath.Join(workspacePath, r.s.GetId().Hex(), "file2.txt"))
}
