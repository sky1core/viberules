package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateAndRemoveSymlink(t *testing.T) {
	// Create temp directory (outside project folder)
	tempDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	sourceContent := "test content"
	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Set target path
	targetFile := filepath.Join(tempDir, "target.txt")

	// Create symlink
	if err := createSymlink(sourceFile, targetFile); err != nil {
		t.Fatalf("createSymlink() failed: %v", err)
	}

	// Check if symlink was created correctly
	if !IsSymlinkValid(targetFile, sourceFile) {
		t.Error("Created symlink is not valid")
	}

	// Read content through symlink
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Failed to read through symlink: %v", err)
	}

	if string(content) != sourceContent {
		t.Errorf("Content through symlink = %s, want %s", string(content), sourceContent)
	}

	// Remove symlink
	if err := removeSymlink(targetFile); err != nil {
		t.Fatalf("removeSymlink() failed: %v", err)
	}

	// Check if symlink was removed
	if _, err := os.Lstat(targetFile); !os.IsNotExist(err) {
		t.Error("Symlink was not removed")
	}
}

func TestCreateSymlinkWithSubdirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Target path in subdirectory
	targetFile := filepath.Join(tempDir, "subdir", "target.txt")

	// Create symlink (subdirectory should be auto-created)
	if err := createSymlink(sourceFile, targetFile); err != nil {
		t.Fatalf("createSymlink() with subdirectory failed: %v", err)
	}

	// Check symlink
	if !IsSymlinkValid(targetFile, sourceFile) {
		t.Error("Symlink in subdirectory is not valid")
	}

	// Check if subdirectory was created
	if _, err := os.Stat(filepath.Dir(targetFile)); os.IsNotExist(err) {
		t.Error("Subdirectory was not created")
	}
}

func TestCreateTargetSymlinks(t *testing.T) {
	tempDir := t.TempDir()

	// Change current directory to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create .viberules directory
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}

	// Create source file
	sourceFile := ".viberules/rules.md"
	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file %s: %v", sourceFile, err)
	}

	// Create symlink for claude target
	if err := CreateTargetSymlinks("claude"); err != nil {
		t.Fatalf("CreateTargetSymlinks(claude) failed: %v", err)
	}

	// Check created symlink
	expectedLinks := []struct {
		target string
		source string
	}{
		{"CLAUDE.md", ".viberules/rules.md"},
	}

	for _, link := range expectedLinks {
		if !IsSymlinkValid(link.target, link.source) {
			t.Errorf("Symlink %s -> %s is not valid", link.target, link.source)
		}
	}
}

func TestRemoveTargetSymlinks(t *testing.T) {
	tempDir := t.TempDir()

	// Change current directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create .viberules directory
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}

	// Create source file
	sourceFile := ".viberules/rules.md"
	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file %s: %v", sourceFile, err)
	}

	// Create symlink
	if err := CreateTargetSymlinks("claude"); err != nil {
		t.Fatalf("CreateTargetSymlinks(claude) failed: %v", err)
	}

	// Remove symlink
	if err := RemoveTargetSymlinks("claude"); err != nil {
		t.Fatalf("RemoveTargetSymlinks(claude) failed: %v", err)
	}

	// Check if symlink was removed
	targetFiles := []string{"CLAUDE.md"}
	for _, file := range targetFiles {
		if _, err := os.Lstat(file); !os.IsNotExist(err) {
			t.Errorf("Symlink %s was not removed", file)
		}
	}
}

func TestCreateAllSymlinks(t *testing.T) {
	tempDir := t.TempDir()

	// Change current directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create .viberules directory
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}

	// Create source file
	sourceFile := ".viberules/rules.md"
	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file %s: %v", sourceFile, err)
	}

	// Create all symlinks
	if err := CreateAllSymlinks(); err != nil {
		t.Fatalf("CreateAllSymlinks() failed: %v", err)
	}

	// Check if all symlinks were created correctly
	valid, missing := CheckAllSymlinks()
	if !valid {
		t.Errorf("Not all symlinks are valid. Missing: %v", missing)
	}

	// Check if required directories were created
	requiredDirs := GetRequiredDirectories()
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Required directory %s was not created", dir)
		}
	}
}

func TestIsSymlinkValid(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Target path
	targetFile := filepath.Join(tempDir, "target.txt")

	// No symlink case
	if IsSymlinkValid(targetFile, sourceFile) {
		t.Error("IsSymlinkValid() should return false for non-existent symlink")
	}

	// Create symlink
	if err := os.Symlink(sourceFile, targetFile); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Valid symlink
	if !IsSymlinkValid(targetFile, sourceFile) {
		t.Error("IsSymlinkValid() should return true for valid symlink")
	}

	// Check with wrong target
	wrongSource := filepath.Join(tempDir, "wrong.txt")
	if IsSymlinkValid(targetFile, wrongSource) {
		t.Error("IsSymlinkValid() should return false for wrong target")
	}

	// Create regular file (not symlink)
	regularFile := filepath.Join(tempDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// Regular file should return false
	if IsSymlinkValid(regularFile, sourceFile) {
		t.Error("IsSymlinkValid() should return false for regular file")
	}
}

func TestSymlinkErrorCases(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test creating symlink with valid source
	if err := os.WriteFile("source.txt", []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// createSymlink actually creates the parent directory automatically
	if err := createSymlink("source.txt", "subdir/link.txt"); err != nil {
		t.Errorf("createSymlink should succeed: %v", err)
	}

	// Test removing non-existent symlink
	if err := removeSymlink("nonexistent.txt"); err != nil {
		t.Errorf("removeSymlink should succeed for non-existent file: %v", err)
	}

	// SECURITY TEST: Ensure removeSymlink refuses to remove regular files
	regularFile := "important.txt"
	if err := os.WriteFile(regularFile, []byte("critical data"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// Try to remove regular file - should fail
	err = removeSymlink(regularFile)
	if err == nil {
		t.Fatal("SECURITY: removeSymlink should refuse to remove regular files")
	}
	if !strings.Contains(err.Error(), "not a symlink") {
		t.Errorf("Expected 'not a symlink' error, got: %v", err)
	}

	// Verify file still exists
	if _, err := os.Stat(regularFile); os.IsNotExist(err) {
		t.Fatal("SECURITY BREACH: Regular file was deleted!")
	}

	// Test with directory
	if err := os.Mkdir("testdir", 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Try to remove directory - should fail
	err = removeSymlink("testdir")
	if err == nil {
		t.Fatal("SECURITY: removeSymlink should refuse to remove directories")
	}
	if !strings.Contains(err.Error(), "not a symlink") {
		t.Errorf("Expected 'not a symlink' error for directory, got: %v", err)
	}

	// Verify directory still exists
	if _, err := os.Stat("testdir"); os.IsNotExist(err) {
		t.Fatal("SECURITY BREACH: Directory was deleted!")
	}
}

func TestCreateTargetSymlinksErrors(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test with invalid target
	if err := CreateTargetSymlinks("invalid"); err == nil {
		t.Error("CreateTargetSymlinks should fail for invalid target")
	}

	// Create .viberules directory but without source files
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}

	// Test without source files - this should actually succeed but symlinks will be broken
	// The function doesn't validate source file existence before creating symlinks
	if err := CreateTargetSymlinks("claude"); err != nil {
		t.Logf("CreateTargetSymlinks with missing source files: %v", err)
	}
}

func TestCreateSymlinkOverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// Create source files
	sourceFile1 := filepath.Join(tempDir, "source1.txt")
	sourceFile2 := filepath.Join(tempDir, "source2.txt")
	targetFile := filepath.Join(tempDir, "target.txt")

	if err := os.WriteFile(sourceFile1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create source1: %v", err)
	}
	if err := os.WriteFile(sourceFile2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create source2: %v", err)
	}

	// Create first symlink
	if err := createSymlink(sourceFile1, targetFile); err != nil {
		t.Fatalf("Failed to create first symlink: %v", err)
	}

	// Overwrite with second symlink
	if err := createSymlink(sourceFile2, targetFile); err != nil {
		t.Fatalf("Failed to overwrite symlink: %v", err)
	}

	// Check if pointing to second source
	if !IsSymlinkValid(targetFile, sourceFile2) {
		t.Error("Symlink was not properly overwritten")
	}

	// Check content
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Failed to read target file: %v", err)
	}

	if string(content) != "content2" {
		t.Errorf("Target file content = %s, want content2", string(content))
	}
}
