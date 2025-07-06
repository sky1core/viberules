package core

import "path/filepath"

// Target represents an AI assistant target with its symlink paths
type Target struct {
	Name  string
	Links []SymlinkDef
}

// SymlinkDef defines a symlink mapping
type SymlinkDef struct {
	Source string // relative path to single rules file
	Target string // destination path for the symlink
}

// GetAllTargets returns all supported AI assistant targets
func GetAllTargets() []Target {
	return []Target{
		{
			Name: "claude",
			Links: []SymlinkDef{
				{Source: filepath.Join(".viberules", "rules.md"), Target: "CLAUDE.md"},
			},
		},
		{
			Name: "amazonq",
			Links: []SymlinkDef{
				{Source: filepath.Join("..", "..", ".viberules", "rules.md"), Target: filepath.Join(".amazonq", "rules", "AMAZONQ.md")},
			},
		},
		{
			Name: "gemini",
			Links: []SymlinkDef{
				{Source: filepath.Join(".viberules", "rules.md"), Target: "GEMINI.md"},
			},
		},
		{
			Name: "codex",
			Links: []SymlinkDef{
				{Source: filepath.Join(".viberules", "rules.md"), Target: "AGENTS.md"},
			},
		},
	}
}

// GetRequiredDirectories returns directories that need to be created
func GetRequiredDirectories() []string {
	return []string{
		filepath.Join(".amazonq", "rules"),
	}
}
