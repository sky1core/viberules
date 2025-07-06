package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateAllSymlinks creates symlinks for all AI assistant targets
func CreateAllSymlinks() error {
	targets := GetAllTargets()

	// Create required directories first
	for _, dir := range GetRequiredDirectories() {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create symlinks for each target
	for _, target := range targets {
		for _, link := range target.Links {
			if err := createSymlink(link.Source, link.Target); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %w", target.Name, err)
			}
		}
	}

	return nil
}

// RemoveAllSymlinks removes all symlinks created by viberules
func RemoveAllSymlinks() error {
	targets := GetAllTargets()

	for _, target := range targets {
		for _, link := range target.Links {
			if err := removeSymlink(link.Target); err != nil {
				return fmt.Errorf("failed to remove symlink for %s: %w", target.Name, err)
			}
		}
	}

	return nil
}

// createSymlink creates a symlink, removing existing file if necessary
func createSymlink(source, target string) error {
	// Clean paths to prevent path traversal
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	// Remove existing file/symlink if it exists
	if err := removeSymlink(target); err != nil {
		return err
	}

	// Create parent directory if needed
	targetDir := filepath.Dir(target)
	if targetDir != "." {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
	}

	// Create the symlink
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("failed to create symlink %s -> %s: %w", target, source, err)
	}

	return nil
}

// removeSymlink removes a symlink or file if it exists
func removeSymlink(path string) error {
	path = filepath.Clean(path)

	// Check if file exists and get info
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to remove
	}
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	// SECURITY: Only remove if it's actually a symlink
	// This prevents accidental deletion of regular files or directories
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("refusing to remove %s: not a symlink", path)
	}

	// Safe to remove - it's confirmed to be a symlink
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove symlink %s: %w", path, err)
	}

	return nil
}

// IsSymlinkValid checks if a symlink exists and points to the correct target
func IsSymlinkValid(linkPath, expectedTarget string) bool {
	linkPath = filepath.Clean(linkPath)
	expectedTarget = filepath.Clean(expectedTarget)

	// Check if symlink exists
	info, err := os.Lstat(linkPath)
	if err != nil {
		return false
	}

	// Check if it's actually a symlink
	if info.Mode()&os.ModeSymlink == 0 {
		return false
	}

	// Check if it points to the correct target
	actualTarget, err := os.Readlink(linkPath)
	if err != nil {
		return false
	}

	// Check if target path matches
	if filepath.Clean(actualTarget) != expectedTarget {
		return false
	}

	// Check if the target actually exists (for broken symlinks)
	_, err = os.Stat(linkPath) // This will fail for broken symlinks
	return err == nil
}

// CheckAllSymlinks verifies all symlinks are properly created
func CheckAllSymlinks() (bool, []string) {
	var missing []string
	allValid := true

	targets := GetAllTargets()
	for _, target := range targets {
		for _, link := range target.Links {
			if !IsSymlinkValid(link.Target, link.Source) {
				missing = append(missing, fmt.Sprintf("%s (%s)", link.Target, target.Name))
				allValid = false
			}
		}
	}

	return allValid, missing
}

// CreateTargetSymlinks creates symlinks for a specific target
func CreateTargetSymlinks(targetName string) error {
	targets := GetAllTargets()

	for _, target := range targets {
		if target.Name == targetName {
			// Create required directories first
			for _, dir := range GetRequiredDirectories() {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", dir, err)
				}
			}

			// Create symlinks for this target
			for _, link := range target.Links {
				if err := createSymlink(link.Source, link.Target); err != nil {
					return fmt.Errorf("failed to create symlink: %w", err)
				}
			}
			return nil
		}
	}

	return fmt.Errorf("target %s not found", targetName)
}

// RemoveTargetSymlinks removes symlinks for a specific target
func RemoveTargetSymlinks(targetName string) error {
	targets := GetAllTargets()

	for _, target := range targets {
		if target.Name == targetName {
			for _, link := range target.Links {
				if err := removeSymlink(link.Target); err != nil {
					return fmt.Errorf("failed to remove symlink: %w", err)
				}
			}
			return nil
		}
	}

	return fmt.Errorf("target %s not found", targetName)
}
