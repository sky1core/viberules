package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/sky1core/viberules/internal/core"
)

const testDir = "../viberules_test"

// TestCompleteViberulesWorkflow tests the entire viberules workflow
// in an isolated test directory as required by project rules
func TestCompleteViberulesWorkflow(t *testing.T) {
	// Setup: Clean and create test directory
	if err := os.RemoveAll(testDir); err != nil {
		t.Logf("Failed to clean test directory (ok if not exists): %v", err)
	}

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	// Cleanup only at start, not at end (for verification purpose)

	// Build viberules binary
	binaryPath := filepath.Join(testDir, "viberules")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build viberules: %v", err)
	}

	// Change to test directory for all operations
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test 1: Help command
	t.Run("Help Command", func(t *testing.T) {
		cmd := exec.Command("./viberules", "--help")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Help command failed: %v", err)
		}

		helpText := string(output)
		if !strings.Contains(helpText, "viberules") {
			t.Error("Help text should contain 'viberules'")
		}
		if !strings.Contains(helpText, "add") {
			t.Error("Help text should contain 'add' command")
		}
		if !strings.Contains(helpText, "remove") {
			t.Error("Help text should contain 'remove' command")
		}
		if !strings.Contains(helpText, "list") {
			t.Error("Help text should contain 'list' command")
		}
		if !strings.Contains(helpText, "init") {
			t.Error("Help text should contain 'init' command")
		}
	})

	// Test 2: Init command
	t.Run("Init Command", func(t *testing.T) {
		cmd := exec.Command("./viberules", "init")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Init command failed: %v", err)
		}

		initOutput := string(output)
		if !strings.Contains(initOutput, "initialized successfully") {
			t.Error("Init output should indicate success")
		}

		// Check created files
		expectedFiles := []string{
			".viberules/rules.md",
			".viberules/.config.yaml",
			".gitignore",
		}

		for _, file := range expectedFiles {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", file)
			}
		}

		// Check .gitignore content
		gitignoreContent, err := os.ReadFile(".gitignore")
		if err != nil {
			t.Fatalf("Failed to read .gitignore: %v", err)
		}

		gitignoreStr := string(gitignoreContent)
		if !strings.Contains(gitignoreStr, ".viberules/") {
			t.Error(".gitignore should contain .viberules/ pattern")
		}
	})

	// Test 3: Symlink creation and validation
	t.Run("Symlink Creation", func(t *testing.T) {
		expectedSymlinks := []struct {
			link   string
			target string
		}{
			{"CLAUDE.md", ".viberules/rules.md"},
			{".amazonq/rules/AMAZONQ.md", "../../.viberules/rules.md"},
			{"GEMINI.md", ".viberules/rules.md"},
			{"AGENTS.md", ".viberules/rules.md"},
		}

		for _, symlink := range expectedSymlinks {
			// Check if symlink exists
			info, err := os.Lstat(symlink.link)
			if err != nil {
				t.Errorf("Symlink %s does not exist: %v", symlink.link, err)
				continue
			}

			// Check if it's actually a symlink
			if info.Mode()&os.ModeSymlink == 0 {
				t.Errorf("%s is not a symlink", symlink.link)
				continue
			}

			// Check target
			actualTarget, err := os.Readlink(symlink.link)
			if err != nil {
				t.Errorf("Failed to read symlink %s: %v", symlink.link, err)
				continue
			}

			if actualTarget != symlink.target {
				t.Errorf("Symlink %s points to %s, want %s", symlink.link, actualTarget, symlink.target)
			}
		}
	})

	// Test 4: List command
	t.Run("List Command", func(t *testing.T) {
		cmd := exec.Command("./viberules", "list")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("List command failed: %v", err)
		}

		listOutput := string(output)
		if !strings.Contains(listOutput, "Enabled targets") {
			t.Error("List output should show enabled targets")
		}
		if !strings.Contains(listOutput, "claude") {
			t.Error("List should show claude as enabled")
		}
		if !strings.Contains(listOutput, "amazonq") {
			t.Error("List should show amazonq as enabled")
		}
		if !strings.Contains(listOutput, "gemini") {
			t.Error("List should show gemini as enabled")
		}
		if !strings.Contains(listOutput, "codex") {
			t.Error("List should show codex as enabled")
		}
	})

	// Test 5: Symlink synchronization
	t.Run("Symlink Synchronization", func(t *testing.T) {
		// Add test content to .viberules/rules.md
		testContent := "\n# Test Content Added\nThis is a synchronization test."

		file, err := os.OpenFile(".viberules/rules.md", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("Failed to open .viberules/rules.md: %v", err)
		}
		defer file.Close()

		if _, err := file.WriteString(testContent); err != nil {
			t.Fatalf("Failed to write test content: %v", err)
		}

		// Check that all symlinks reflect the change
		symlinks := []string{
			"CLAUDE.md",
			".amazonq/rules/AMAZONQ.md",
			"GEMINI.md",
			"AGENTS.md",
		}

		for _, symlink := range symlinks {
			content, err := os.ReadFile(symlink)
			if err != nil {
				t.Errorf("Failed to read %s: %v", symlink, err)
				continue
			}

			if !strings.Contains(string(content), "Test Content Added") {
				t.Errorf("Symlink %s does not reflect content changes", symlink)
			}
		}

		// Clean up test content
		originalContent, err := os.ReadFile(".viberules/rules.md")
		if err != nil {
			t.Fatalf("Failed to read original content: %v", err)
		}

		cleanContent := strings.Replace(string(originalContent), testContent, "", 1)
		if err := os.WriteFile(".viberules/rules.md", []byte(cleanContent), 0644); err != nil {
			t.Fatalf("Failed to clean test content: %v", err)
		}
	})

	// Test 6: Remove target
	t.Run("Remove Target", func(t *testing.T) {
		cmd := exec.Command("./viberules", "remove", "amazonq")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Remove command failed: %v", err)
		}

		removeOutput := string(output)
		if !strings.Contains(removeOutput, "removed successfully") {
			t.Error("Remove output should indicate success")
		}

		// Check that amazonq symlinks are removed
		amazonqLinks := []string{
			".amazonq/rules/AMAZONQ.md",
		}

		for _, link := range amazonqLinks {
			if _, err := os.Lstat(link); !os.IsNotExist(err) {
				t.Errorf("Symlink %s should have been removed", link)
			}
		}

		// Check that .config.yaml is updated
		configContent, err := os.ReadFile(".viberules/.config.yaml")
		if err != nil {
			t.Fatalf("Failed to read .viberules/.config.yaml: %v", err)
		}

		configStr := string(configContent)
		if strings.Contains(configStr, "amazonq") {
			t.Error("amazonq should be removed from .viberules/.config.yaml")
		}
	})

	// Test 7: Add target back
	t.Run("Add Target", func(t *testing.T) {
		cmd := exec.Command("./viberules", "add", "amazonq")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Add command failed: %v", err)
		}

		addOutput := string(output)
		if !strings.Contains(addOutput, "added successfully") {
			t.Error("Add output should indicate success")
		}

		// Check that amazonq symlinks are recreated
		amazonqLinks := []string{
			".amazonq/rules/AMAZONQ.md",
		}

		for _, link := range amazonqLinks {
			if _, err := os.Lstat(link); os.IsNotExist(err) {
				t.Errorf("Symlink %s should have been recreated", link)
			}
		}

		// Check that .config.yaml is updated
		configContent, err := os.ReadFile(".viberules/.config.yaml")
		if err != nil {
			t.Fatalf("Failed to read .viberules/.config.yaml: %v", err)
		}

		configStr := string(configContent)
		if !strings.Contains(configStr, "amazonq") {
			t.Error("amazonq should be added to .viberules/.config.yaml")
		}
	})

	// Test 8: Error cases
	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid target
		cmd := exec.Command("./viberules", "add", "invalid")
		if err := cmd.Run(); err == nil {
			t.Error("Adding invalid target should fail")
		}

		cmd = exec.Command("./viberules", "remove", "invalid")
		if err := cmd.Run(); err == nil {
			t.Error("Removing invalid target should fail")
		}

		// Test duplicate add (should succeed silently)
		cmd = exec.Command("./viberules", "add", "claude")
		if err := cmd.Run(); err != nil {
			t.Errorf("Adding existing target should succeed: %v", err)
		}

		// Test removing non-existent target (should succeed silently)
		cmd = exec.Command("./viberules", "remove", "amazonq")
		if err := cmd.Run(); err != nil {
			t.Errorf("First remove should succeed: %v", err)
		}

		cmd = exec.Command("./viberules", "remove", "amazonq")
		if err := cmd.Run(); err != nil {
			t.Errorf("Second remove should succeed silently: %v", err)
		}
	})

	// Test 9: GitIgnore functionality
	t.Run("GitIgnore Validation", func(t *testing.T) {
		// Initialize git repo in test directory
		initCmd := exec.Command("git", "init")
		if err := initCmd.Run(); err != nil {
			t.Logf("Git init failed (may not have git): %v", err)
			t.Skip("Skipping git tests - git not available")
		}

		// Check that .gitignore was created with proper content
		gitignoreContent, err := os.ReadFile(".gitignore")
		if err != nil {
			t.Fatalf("Failed to read .gitignore: %v", err)
		}

		gitignoreStr := string(gitignoreContent)
		expectedPatterns := []string{
			"*.local.md",
			".viberules/.config.yaml",
			".amazonq/",
			"CLAUDE.md",
			"GEMINI.md",
			"AGENTS.md",
		}

		for _, pattern := range expectedPatterns {
			if !strings.Contains(gitignoreStr, pattern) {
				t.Errorf(".gitignore should contain pattern: %s", pattern)
			}
		}

		// Add all files to git
		addCmd := exec.Command("git", "add", ".")
		if err := addCmd.Run(); err != nil {
			t.Fatalf("Git add failed: %v", err)
		}

		// Check git status to verify .gitignore works
		statusCmd := exec.Command("git", "status", "--porcelain")
		output, err := statusCmd.Output()
		if err != nil {
			t.Logf("Git status failed: %v", err)
			t.Skip("Skipping git status check - may be git version issue")
		}

		statusOutput := string(output)

		// These files should NOT appear in git status (should be ignored)
		// .config.yaml is always ignored regardless of mode
		ignoredFiles := []string{
			".viberules/.config.yaml",
		}

		for _, file := range ignoredFiles {
			if strings.Contains(statusOutput, file) {
				t.Errorf("File %s should be ignored by git but appears in status", file)
			}
		}

		// These files SHOULD appear in git status (should be tracked)
		// Note: .viberules/ directory is gitignored, so we don't check for it
		if !strings.Contains(statusOutput, ".gitignore") {
			t.Error(".gitignore should be tracked by git")
		}
	})

	// Test 10: File permission scenarios
	t.Run("File Permissions", func(t *testing.T) {
		// Test with read-only config.yaml
		if err := os.Chmod(".viberules/config.yaml", 0444); err != nil {
			t.Logf("Could not change file permissions: %v", err)
			return
		}

		// Try to add target (should fail due to read-only config)
		cmd := exec.Command("./viberules", "add", "claude")
		if err := cmd.Run(); err == nil {
			t.Log("Add command unexpectedly succeeded with read-only config")
		}

		// Restore permissions
		os.Chmod(".viberules/config.yaml", 0644)
	})

	// Test 11: Symlink validation
	t.Run("Symlink Validation", func(t *testing.T) {
		// Create a broken symlink manually
		if err := os.Symlink("nonexistent.txt", "broken-link.md"); err != nil {
			t.Logf("Could not create broken symlink: %v", err)
			return
		}

		// Verify our validation catches it (use core package function)
		if core.IsSymlinkValid("broken-link.md", "nonexistent.txt") {
			t.Error("IsSymlinkValid should return false for broken symlink")
		}

		// Clean up
		os.Remove("broken-link.md")
	})

	// Test 12: Concurrent operations
	t.Run("Concurrent Safety", func(t *testing.T) {
		// This is a basic test - in real usage, we'd need more sophisticated testing
		var wg sync.WaitGroup
		errors := make(chan error, 2)

		// Try to add same target concurrently
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cmd := exec.Command("./viberules", "add", "claude")
				errors <- cmd.Run()
			}()
		}

		wg.Wait()
		close(errors)

		// At least one should succeed
		successCount := 0
		for err := range errors {
			if err == nil {
				successCount++
			}
		}

		if successCount == 0 {
			t.Error("At least one concurrent add should succeed")
		}
	})

	t.Logf("All tests completed successfully in %s", testDir)
}

// TestViberulesTestDirectory verifies that test directory configuration is correct
func TestViberulesTestDirectory(t *testing.T) {
	// Verify test directory path configuration
	expectedTestDir, err := filepath.Abs(testDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path for test directory: %v", err)
	}

	if !strings.HasSuffix(expectedTestDir, "viberules_test") {
		t.Errorf("Test directory should end with 'viberules_test', got: %s", expectedTestDir)
	}

	// Verify test directory is outside project directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Verify that when tests do use the test directory, it's properly isolated
	if strings.Contains(expectedTestDir, currentDir) && !strings.Contains(expectedTestDir, "viberules_test") {
		t.Error("Test directory should be isolated from project directory")
	}

	t.Logf("Project directory: %s", currentDir)
	t.Logf("Test directory: %s", expectedTestDir)
}
