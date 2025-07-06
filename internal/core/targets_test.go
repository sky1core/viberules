package core

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetAllTargets(t *testing.T) {
	targets := GetAllTargets()

	// Should have 4 targets
	if len(targets) != 4 {
		t.Errorf("GetAllTargets() = %d targets, want 4", len(targets))
	}

	// Each target should have correct name
	expectedNames := []string{"claude", "amazonq", "gemini", "codex"}
	var actualNames []string
	for _, target := range targets {
		actualNames = append(actualNames, target.Name)
	}

	for _, expected := range expectedNames {
		found := false
		for _, actual := range actualNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected target %s not found in %v", expected, actualNames)
		}
	}
}

func TestTargetStructure(t *testing.T) {
	targets := GetAllTargets()

	tests := []struct {
		name              string
		targetName        string
		expectedLinkCount int
		expectedSources   []string
		expectedTargets   []string
	}{
		{
			name:              "claude target",
			targetName:        "claude",
			expectedLinkCount: 1,
			expectedSources:   []string{filepath.Join(".viberules", "rules.md")},
			expectedTargets:   []string{"CLAUDE.md"},
		},
		{
			name:              "amazonq target",
			targetName:        "amazonq",
			expectedLinkCount: 1,
			expectedSources:   []string{filepath.Join("..", "..", ".viberules", "rules.md")},
			expectedTargets:   []string{filepath.Join(".amazonq", "rules", "AMAZONQ.md")},
		},
		{
			name:              "gemini target",
			targetName:        "gemini",
			expectedLinkCount: 1,
			expectedSources:   []string{filepath.Join(".viberules", "rules.md")},
			expectedTargets:   []string{"GEMINI.md"},
		},
		{
			name:              "codex target",
			targetName:        "codex",
			expectedLinkCount: 1,
			expectedSources:   []string{filepath.Join(".viberules", "rules.md")},
			expectedTargets:   []string{"AGENTS.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target *Target
			for _, tgt := range targets {
				if tgt.Name == tt.targetName {
					target = &tgt
					break
				}
			}

			if target == nil {
				t.Fatalf("Target %s not found", tt.targetName)
			}

			if len(target.Links) != tt.expectedLinkCount {
				t.Errorf("Target %s has %d links, want %d", tt.targetName, len(target.Links), tt.expectedLinkCount)
			}

			// Verify source and target paths
			var actualSources, actualTargets []string
			for _, link := range target.Links {
				actualSources = append(actualSources, link.Source)
				actualTargets = append(actualTargets, link.Target)
			}

			if !reflect.DeepEqual(actualSources, tt.expectedSources) {
				t.Errorf("Target %s sources = %v, want %v", tt.targetName, actualSources, tt.expectedSources)
			}

			if !reflect.DeepEqual(actualTargets, tt.expectedTargets) {
				t.Errorf("Target %s targets = %v, want %v", tt.targetName, actualTargets, tt.expectedTargets)
			}
		})
	}
}

func TestGetRequiredDirectories(t *testing.T) {
	dirs := GetRequiredDirectories()

	expectedDirs := []string{
		filepath.Join(".amazonq", "rules"),
	}

	if !reflect.DeepEqual(dirs, expectedDirs) {
		t.Errorf("GetRequiredDirectories() = %v, want %v", dirs, expectedDirs)
	}
}
