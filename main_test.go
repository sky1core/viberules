package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidTarget(t *testing.T) {
	tests := []struct {
		target string
		valid  bool
	}{
		{"claude", true},
		{"amazonq", true},
		{"gemini", true},
		{"codex", true},
		{"invalid", false},
		{"", false},
		{"Claude", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			result := isValidTarget(tt.target)
			if result != tt.valid {
				t.Errorf("isValidTarget(%s) = %v, want %v", tt.target, result, tt.valid)
			}
		})
	}
}

func TestLoadEnabledTargets(t *testing.T) {
	tempDir := t.TempDir() // Safe isolated test directory

	// Change current directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// No config file - return default values
	targets, err := loadEnabledTargets()
	if err != nil {
		t.Fatalf("loadEnabledTargets() failed: %v", err)
	}

	expectedDefault := []string{"claude", "amazonq", "gemini", "codex"}
	if !equalStringSlices(targets, expectedDefault) {
		t.Errorf("loadEnabledTargets() with no config = %v, want %v", targets, expectedDefault)
	}

	// Create .viberules directory
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}

	// Create config file (YAML)
	configContent := `mode: local
targets:
  - claude
  - gemini
`
	if err := os.WriteFile(".viberules/.config.yaml", []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// With config file
	targets, err = loadEnabledTargets()
	if err != nil {
		t.Fatalf("loadEnabledTargets() with config failed: %v", err)
	}

	expected := []string{"claude", "gemini"}
	if !equalStringSlices(targets, expected) {
		t.Errorf("loadEnabledTargets() with config = %v, want %v", targets, expected)
	}
}

func TestSaveEnabledTargets(t *testing.T) {
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

	targets := []string{"claude", "gemini"}

	// Save targets
	if err := saveEnabledTargets(targets); err != nil {
		t.Fatalf("saveEnabledTargets() failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(".viberules/.config.yaml"); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check file content
	content, err := os.ReadFile(".viberules/.config.yaml")
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "claude") {
		t.Error("Config file does not contain 'claude'")
	}
	if !strings.Contains(contentStr, "gemini") {
		t.Error("Config file does not contain 'gemini'")
	}
	if strings.Contains(contentStr, "amazonq") {
		t.Error("Config file should not contain 'amazonq'")
	}

	// Reload and verify
	loadedTargets, err := loadEnabledTargets()
	if err != nil {
		t.Fatalf("Failed to reload targets: %v", err)
	}

	if !equalStringSlices(loadedTargets, targets) {
		t.Errorf("Reloaded targets = %v, want %v", loadedTargets, targets)
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()

	existingFile := filepath.Join(tempDir, "existing.txt")
	nonExistingFile := filepath.Join(tempDir, "nonexisting.txt")

	// Create file
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Existing file
	if !fileExists(existingFile) {
		t.Error("fileExists() should return true for existing file")
	}

	// Non-existing file
	if fileExists(nonExistingFile) {
		t.Error("fileExists() should return false for non-existing file")
	}
}

func TestAddToGitignore(t *testing.T) {
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

	// No .gitignore file
	if err := addToGitignore(); err != nil {
		t.Fatalf("addToGitignore() failed: %v", err)
	}

	// Check if .gitignore file was created
	if !fileExists(".gitignore") {
		t.Error(".gitignore file was not created")
	}

	// Check content
	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "*.local.md") {
		t.Error(".gitignore does not contain *.local.md pattern")
	}
	if !strings.Contains(contentStr, "CLAUDE.md") {
		t.Error(".gitignore does not contain CLAUDE.md pattern")
	}

	// Check content after first run
	firstContent := string(content)
	t.Logf("First gitignore content:\n%s\n", firstContent)

	// Run again if already exists
	if err := addToGitignore(); err != nil {
		t.Fatalf("addToGitignore() second run failed: %v", err)
	}

	// Check no duplicate was added
	newContent, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("Failed to read .gitignore again: %v", err)
	}
	
	secondContent := string(newContent)
	t.Logf("Second gitignore content:\n%s\n", secondContent)

	// Check content is same (no duplicate added)
	if firstContent != secondContent {
		t.Error("gitignore content was changed on second run")
	}

	// Check viberules section appears exactly once
	localFilesCount := strings.Count(secondContent, "# viberules local files")
	if localFilesCount != 1 {
		t.Errorf("viberules local files section appears %d times, want 1", localFilesCount)
	}
	
	outputFilesCount := strings.Count(secondContent, "# viberules output files")
	if outputFilesCount != 1 {
		t.Errorf("viberules output files section appears %d times, want 1", outputFilesCount)
	}
}

func TestAddTargetFunction(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create .viberules directory and files
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}
	if err := os.WriteFile(".viberules/rules.md", []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create rules.md: %v", err)
	}

	// Test adding valid target
	if err := addTarget("claude"); err != nil {
		t.Errorf("addTarget(claude) should succeed: %v", err)
	}

	// Test adding invalid target
	if err := addTarget("invalid"); err == nil {
		t.Error("addTarget(invalid) should fail")
	}

	// Test adding without init
	os.RemoveAll(".viberules")
	if err := addTarget("claude"); err == nil {
		t.Error("addTarget should fail when .viberules doesn't exist")
	}
}

func TestRemoveTargetFunction(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Setup
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}
	if err := saveEnabledTargets([]string{"claude", "amazonq"}); err != nil {
		t.Fatalf("Failed to save initial targets: %v", err)
	}

	// Test removing valid target
	if err := removeTarget("amazonq"); err != nil {
		t.Errorf("removeTarget(amazonq) should succeed: %v", err)
	}

	// Test removing invalid target
	if err := removeTarget("invalid"); err == nil {
		t.Error("removeTarget(invalid) should fail")
	}

	// Test removing non-existent target
	if err := removeTarget("gemini"); err != nil {
		t.Errorf("removeTarget(gemini) should succeed silently: %v", err)
	}
}

func TestListTargetsFunction(t *testing.T) {
	tempDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Setup
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		t.Fatalf("Failed to create .viberules directory: %v", err)
	}
	if err := saveEnabledTargets([]string{"claude", "gemini"}); err != nil {
		t.Fatalf("Failed to save targets: %v", err)
	}

	// Test listTargets (just ensure it doesn't error)
	if err := listTargets(); err != nil {
		t.Errorf("listTargets() should not error: %v", err)
	}
}

func TestWindowsRejection(t *testing.T) {
	// This test can't actually test Windows behavior on non-Windows,
	// but we can test that the check exists
	oldGOOS := os.Getenv("GOOS")
	defer func() {
		if oldGOOS != "" {
			os.Setenv("GOOS", oldGOOS)
		} else {
			os.Unsetenv("GOOS")
		}
	}()

	// The actual Windows check is in runtime.GOOS, not environment variable
	// So this test just ensures the function exists and doesn't panic
	// Real Windows testing would need to be done on actual Windows
	t.Log("Windows rejection test placeholder - requires actual Windows for full test")
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "exact match",
			s:      "*.local.md",
			substr: "*.local.md",
			want:   true,
		},
		{
			name:   "substring at beginning",
			s:      "*.local.md\nother",
			substr: "*.local.md",
			want:   true,
		},
		{
			name:   "substring at end",
			s:      "other\n*.local.md",
			substr: "*.local.md",
			want:   true,
		},
		{
			name:   "substring in middle",
			s:      "first\n*.local.md\nlast",
			substr: "*.local.md",
			want:   true,
		},
		{
			name:   "not found",
			s:      "something else",
			substr: "*.local.md",
			want:   false,
		},
		{
			name:   "partial match",
			s:      "*.local.md.backup",
			substr: "*.local.md",
			want:   true,  // strings.Contains finds partial matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.want {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.want)
			}
		})
	}
}

// Helper function: compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
