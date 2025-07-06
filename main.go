package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sky1core/viberules/internal/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const version = "0.2.0"

// WARNING: DO NOT CHANGE THESE CONSTANTS!
// These strings are used to identify viberules sections in existing .gitignore files.
// Changing them will break gitignore updates for users who already have viberules installed.
const (
	gitignoreSectionPrefix = "# viberules"
	gitignoreLocalMode     = "# viberules (local mode"
	gitignoreLocalFiles    = "# viberules local files"
	gitignoreConfigFile    = "# viberules config file"
	gitignoreOutputFiles   = "# viberules output files"
)

var (
	silent bool
	force  bool
)

var rootCmd = &cobra.Command{
	Use:   "viberules",
	Short: "AI assistant rules management tool using symlinks",
	Long: `viberules is a CLI tool for managing AI coding assistant rules 
(Claude Code, Amazon Q Developer, Gemini Code Assist, etc.) using symlinks.

Key features:
- Manage only 2 files: viberules.md, viberules.local.md
- Real-time sync via symlinks
- Individual target management (add/remove)`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("Windows is not supported. Please use macOS or Linux")
		}
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize viberules project",
	Long: `Create viberules files and symlinks in the current directory.

Created files:
- rules.md (single rules file for all AI tools)
- Symlinks for each AI tool (CLAUDE.md, GEMINI.md, AGENTS.md, etc.)
- Mode-aware .gitignore configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initProject()
	},
}

var addCmd = &cobra.Command{
	Use:   "add [target]",
	Short: "Add target",
	Long: `Enable the specified AI assistant target.
Available targets: claude, amazonq, gemini`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addTarget(args[0])
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [target]",
	Short: "Remove target",
	Long: `Disable the specified AI assistant target.
Available targets: claude, amazonq, gemini`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeTarget(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List enabled targets",
	Long:  "Show currently enabled AI assistant targets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listTargets()
	},
}

var modeCmd = &cobra.Command{
	Use:   "mode [public|local]",
	Short: "Get or set project mode",
	Long: `Get or set the project mode.
	
Modes:
- public: .viberules directory is tracked by git (shared rules)
- local: .viberules directory is ignored by git (personal rules)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show current mode
			mode := getProjectMode()
			fmt.Printf("Current mode: %s\n", mode)
			return nil
		}
		
		if len(args) != 1 {
			return fmt.Errorf("usage: viberules mode [public|local]")
		}
		
		return setModeCommand(args[0])
	},
}

func initProject() error {
	if !silent {
		fmt.Println("🚀 Initializing viberules project...")
	}

	// Check if .viberules directory already exists
	if stat, err := os.Stat(".viberules"); err == nil && stat.IsDir() {
		if !force {
			return fmt.Errorf(".viberules directory already exists. Use --force to reinitialize")
		}
		if !silent {
			fmt.Println("⚠️  Reinitializing existing project...")
			fmt.Println("   - Existing .viberules/rules.md will be preserved")
			fmt.Println("   - Missing files will be created")
			fmt.Println("   - Symlinks will be recreated")
		}
	}

	// Create .viberules directory
	if err := os.MkdirAll(".viberules", 0755); err != nil {
		return fmt.Errorf("failed to create .viberules directory: %w", err)
	}

	// Create single rules.md file only if it doesn't exist
	rulesFile := ".viberules/rules.md"
	if !fileExists(rulesFile) {
		rulesContent := `# AI Assistant Rules

> ⚠️ IMPORTANT: Edit THIS FILE (rules.md) to update rules for ALL AI assistants
> Changes here automatically apply to Claude, Amazon Q, Gemini, Codex, etc.

## Project Overview
Describe your project, tech stack, and coding standards here.

## Coding Standards
- Use TypeScript with strict mode
- Follow ESLint configuration
- Write unit tests for all functions
- Use descriptive variable names

## Architecture Guidelines
- Follow clean architecture principles
- Separate business logic from UI
- Use dependency injection

## Git Workflow
- Use conventional commits
- Create feature branches
- Require code review for main branch

---
*This file is automatically linked to all AI assistants via viberules*
`

		if err := os.WriteFile(rulesFile, []byte(rulesContent), 0644); err != nil {
			return fmt.Errorf("failed to create .viberules/rules.md: %w", err)
		}
		if !silent && force {
			fmt.Println("📝 Created .viberules/rules.md")
		}
	} else if !silent && force {
		fmt.Println("📋 Preserved existing .viberules/rules.md")
	}

	// Add to .gitignore
	if err := addToGitignore(); err != nil {
		if !silent {
			fmt.Printf("⚠️  Failed to update .gitignore: %v\n", err)
		}
	} else if !silent {
		fmt.Println("📝 Added *.local.md to .gitignore")
	}

	// Create symlinks
	if err := core.CreateAllSymlinks(); err != nil {
		return fmt.Errorf("failed to create symlinks: %w", err)
	}

	// Initialize default config (local mode, all targets)
	defaultConfig := &Config{
		Mode:    "local",
		Targets: []string{"claude", "amazonq", "gemini", "codex"},
	}
	if err := saveConfig(defaultConfig); err != nil {
		if !silent {
			fmt.Printf("⚠️  Failed to create config file: %v\n", err)
		}
	}

	if !silent {
		fmt.Println("✅ viberules project initialized successfully!")
		fmt.Println("📁 Created files:")
		fmt.Println("   - .viberules/rules.md (rules shared by all AI tools)")
		fmt.Println("   - Symlinks for each AI tool")
		fmt.Println("")
		fmt.Println("Next steps:")
		fmt.Println("1. Edit .viberules/rules.md to write your project rules")
		fmt.Println("2. Use 'viberules remove [target]' to remove unnecessary targets")
	}

	return nil
}

func addTarget(target string) error {
	if !isValidTarget(target) {
		return fmt.Errorf("invalid target: %s (available: claude, amazonq, gemini, codex)", target)
	}

	if !fileExists(".viberules/rules.md") {
		return fmt.Errorf(".viberules/rules.md not found. Run 'viberules init' first")
	}

	// Load current targets
	enabledTargets, err := loadEnabledTargets()
	if err != nil {
		return fmt.Errorf("failed to load target settings: %w", err)
	}

	// Check if already enabled
	for _, enabled := range enabledTargets {
		if enabled == target {
			fmt.Printf("Target '%s' is already enabled\n", target)
			return nil
		}
	}

	// Add target
	enabledTargets = append(enabledTargets, target)

	// Save configuration
	if err := saveEnabledTargets(enabledTargets); err != nil {
		return fmt.Errorf("failed to save target settings: %w", err)
	}

	// Create symlinks for this target
	if err := core.CreateTargetSymlinks(target); err != nil {
		return fmt.Errorf("failed to create symlinks for target '%s': %w", target, err)
	}

	fmt.Printf("✅ Target '%s' added successfully\n", target)
	return nil
}

func removeTarget(target string) error {
	if !isValidTarget(target) {
		return fmt.Errorf("invalid target: %s (available: claude, amazonq, gemini, codex)", target)
	}

	// Load current targets
	enabledTargets, err := loadEnabledTargets()
	if err != nil {
		return fmt.Errorf("failed to load target settings: %w", err)
	}

	// Remove target from list
	newTargets := make([]string, 0)
	found := false
	for _, enabled := range enabledTargets {
		if enabled != target {
			newTargets = append(newTargets, enabled)
		} else {
			found = true
		}
	}

	if !found {
		fmt.Printf("Target '%s' is not enabled\n", target)
		return nil
	}

	// Save configuration
	if err := saveEnabledTargets(newTargets); err != nil {
		return fmt.Errorf("failed to save target settings: %w", err)
	}

	// Remove symlinks for this target
	if err := core.RemoveTargetSymlinks(target); err != nil {
		return fmt.Errorf("failed to remove symlinks for target '%s': %w", target, err)
	}

	fmt.Printf("✅ Target '%s' removed successfully\n", target)
	return nil
}

func listTargets() error {
	enabledTargets, err := loadEnabledTargets()
	if err != nil {
		return fmt.Errorf("failed to load target settings: %w", err)
	}

	fmt.Println("Enabled targets:")
	if len(enabledTargets) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, target := range enabledTargets {
			fmt.Printf("  - %s\n", target)
		}
	}

	fmt.Println("\nAvailable targets:")
	for _, target := range []string{"claude", "amazonq", "gemini", "codex"} {
		fmt.Printf("  - %s\n", target)
	}

	return nil
}

func setModeCommand(mode string) error {
	if !fileExists(".viberules/rules.md") {
		return fmt.Errorf(".viberules/rules.md not found. Run 'viberules init' first")
	}
	
	if err := setProjectMode(mode); err != nil {
		return err
	}
	
	// Update gitignore based on new mode
	if err := addToGitignore(); err != nil {
		fmt.Printf("⚠️  Failed to update .gitignore: %v\n", err)
	}
	
	fmt.Printf("✅ Project mode set to '%s'\n", mode)
	if mode == "public" {
		fmt.Println("📁 .viberules/rules.md will be tracked by git")
		fmt.Println("🔒 .viberules/.config.yaml will be ignored by git")
	} else {
		fmt.Println("🔒 .viberules directory will be ignored by git")
	}
	
	return nil
}

func isValidTarget(target string) bool {
	for _, valid := range []string{"claude", "amazonq", "gemini", "codex"} {
		if target == valid {
			return true
		}
	}
	return false
}

type Config struct {
	Mode    string   `yaml:"mode"`
	Targets []string `yaml:"targets"`
}

func loadConfig() (*Config, error) {
	configPath := ".viberules/.config.yaml"
	if !fileExists(configPath) {
		// Return default config if no config file exists
		return &Config{
			Mode:    "local", // Default mode changed to local
			Targets: []string{"claude", "amazonq", "gemini", "codex"},
		}, nil
	}

	// Security: Limit config file size to prevent YAML bomb attacks
	info, err := os.Stat(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}
	const maxConfigSize = 1 * 1024 * 1024 // 1MB
	if info.Size() > maxConfigSize {
		return nil, fmt.Errorf("config file too large: %d bytes (max %d)", info.Size(), maxConfigSize)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate mode
	if config.Mode != "local" && config.Mode != "public" {
		config.Mode = "local" // Default value
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configPath := ".viberules/.config.yaml"
	
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func loadEnabledTargets() ([]string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Targets, nil
}

func saveEnabledTargets(targets []string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}
	config.Targets = targets
	return saveConfig(config)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func addToGitignore() error {
	gitignorePath := ".gitignore"
	mode := getProjectMode()

	// Create gitignore content based on mode
	var viberulesSection string
	if mode == "local" {
		// Local mode: ignore entire .viberules directory
		viberulesSection = fmt.Sprintf(`
%s - entire directory ignored)
.viberules/

%s (always ignored)
.viberules/.config.yaml

%s (personal files only)
*.local.md

%s (symlinked)
.amazonq/
CLAUDE.md
GEMINI.md
AGENTS.md
`, gitignoreLocalMode, gitignoreConfigFile, gitignoreLocalFiles, gitignoreOutputFiles)
	} else {
		// Public mode: track .viberules/rules.md but ignore config
		viberulesSection = fmt.Sprintf(`
%s (always ignored)
.viberules/.config.yaml

%s (personal files only)
*.local.md

%s (symlinked)
.amazonq/
CLAUDE.md
GEMINI.md
AGENTS.md
`, gitignoreConfigFile, gitignoreLocalFiles, gitignoreOutputFiles)
	}

	// Read existing .gitignore
	var content []byte
	var err error
	if fileExists(gitignorePath) {
		content, err = os.ReadFile(gitignorePath)
		if err != nil {
			return fmt.Errorf("failed to read .gitignore: %w", err)
		}
	}

	contentStr := string(content)

	// Remove existing viberules section if present
	if contains(contentStr, gitignoreLocalFiles) || contains(contentStr, gitignoreLocalMode) {
		// Simple approach: split by lines and rebuild without viberules section
		lines := strings.Split(contentStr, "\n")
		var newLines []string
		skipSection := false
		
		for _, line := range lines {
			if strings.HasPrefix(line, gitignoreSectionPrefix) {
				skipSection = true
				continue
			}
			
			// End of section when we hit another comment section (not viberules)
			if skipSection {
				if strings.HasPrefix(line, "#") && !strings.Contains(line, "viberules") {
					skipSection = false
					newLines = append(newLines, line)
				}
				// Skip all lines until we hit another section or end of file
				continue
			}
			
			newLines = append(newLines, line)
		}
		
		contentStr = strings.Join(newLines, "\n")
		// Remove trailing empty lines
		contentStr = strings.TrimRight(contentStr, "\n")
	}

	// Add viberules section
	if len(contentStr) > 0 && contentStr[len(contentStr)-1] != '\n' {
		contentStr += "\n"
	}
	contentStr += viberulesSection

	// Write back
	if err := os.WriteFile(gitignorePath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	return nil
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// getProjectMode returns the current project mode (public or local)
func getProjectMode() string {
	config, err := loadConfig()
	if err != nil {
		return "local" // fallback to default
	}
	return config.Mode
}

// setProjectMode sets the project mode (public or local)
func setProjectMode(mode string) error {
	if mode != "public" && mode != "local" {
		return fmt.Errorf("invalid mode: %s (must be 'public' or 'local')", mode)
	}
	
	config, err := loadConfig()
	if err != nil {
		return err
	}
	
	config.Mode = mode
	return saveConfig(config)
}

func init() {
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force reinitialize existing project")
	
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(modeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
