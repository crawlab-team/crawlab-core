package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceV2_List(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	testDir := filepath.Join(rootDir, "dir")
	os.Mkdir(testDir, 0755)
	ioutil.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("hello world"), 0644)
	ioutil.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("hello again"), 0644)
	subDir := filepath.Join(testDir, "subdir")
	os.Mkdir(subDir, 0755)
	ioutil.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("subdir file"), 0644)
	os.Mkdir(filepath.Join(testDir, "empty"), 0755) // explicitly testing empty dir inclusion

	svc := NewFsServiceV2(rootDir)

	files, err := svc.List("dir")
	if err != nil {
		t.Errorf("Failed to list files: %v", err)
	}

	// Assert correct number of items
	if assert.Len(t, files, 1) && assert.Len(t, files[0].GetChildren(), 4) {
		// Use a map to verify presence and characteristics of files/directories to avoid order issues
		items := make(map[string]bool)
		for _, item := range files[0].GetChildren() {
			items[item.GetName()] = item.GetIsDir()
		}

		_, file1Exists := items["file1.txt"]
		_, file2Exists := items["file2.txt"]
		_, subdirExists := items["subdir"]
		_, emptyExists := items["empty"]

		assert.True(t, file1Exists)
		assert.True(t, file2Exists)
		assert.True(t, subdirExists)
		assert.True(t, emptyExists) // Verify that the empty directory is included

		if subdirExists && len(files[0].GetChildren()[2].GetChildren()) > 0 {
			assert.Equal(t, "file3.txt", files[0].GetChildren()[2].GetChildren()[0].GetName())
		}
	}
}

func TestServiceV2_GetFile(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	expectedContent := []byte("hello world")
	ioutil.WriteFile(filepath.Join(rootDir, "file.txt"), expectedContent, 0644)

	svc := NewFsServiceV2(rootDir)

	content, err := svc.GetFile("file.txt")
	if err != nil {
		t.Errorf("Failed to get file: %v", err)
	}
	assert.Equal(t, expectedContent, content)
}

func TestServiceV2_Delete(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	filePath := filepath.Join(rootDir, "file.txt")
	ioutil.WriteFile(filePath, []byte("hello world"), 0644)

	svc := NewFsServiceV2(rootDir)

	// Delete the file
	err = svc.Delete("file.txt")
	if err != nil {
		t.Errorf("Failed to delete file: %v", err)
	}

	// Verify deletion
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}
