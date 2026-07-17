package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// GenerateCommand represents a go generate command extracted from source files
type GenerateCommand struct {
	File    string // Source file path
	Line    int    // Line number in the file
	Command string // The actual command to execute
	Dir     string // Working directory for the command
}

// extractGoGenerateCommands scans Go files and extracts //go:generate commands
func extractGoGenerateCommands(rootDir string, patterns []string) ([]GenerateCommand, error) {
	var commands []GenerateCommand

	// Regex to match //go:generate comments
	generateRegex := regexp.MustCompile(`^//go:generate\s+(.+)$`)

	for _, pattern := range patterns {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip non-Go files and directories starting with . or _
			if info.IsDir() ||
				!strings.HasSuffix(path, ".go") ||
				strings.HasPrefix(info.Name(), ".") ||
				strings.HasPrefix(info.Name(), "_") {
				return nil
			}

			// Check if file matches the pattern
			matched, err := filepath.Match(pattern, info.Name())
			if err != nil {
				return err
			}
			if !matched && pattern != "*" {
				return nil
			}

			// Read file and extract go:generate commands
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()

				if matches := generateRegex.FindStringSubmatch(line); matches != nil {
					cmd := strings.TrimSpace(matches[1])
					// Only process mockgen commands
					if strings.Contains(cmd, "mockgen") {
						commands = append(commands, GenerateCommand{
							File:    path,
							Line:    lineNum,
							Command: cmd,
							Dir:     filepath.Dir(path),
						})
					}
				}
			}

			return scanner.Err()
		})

		if err != nil {
			return nil, err
		}
	}

	return commands, nil
}

// executeCommand runs a single generate command
func executeCommand(cmd GenerateCommand, verbose bool) error {
	if verbose {
		fmt.Printf("  📍 File: %s:%d\n", cmd.File, cmd.Line)
		fmt.Printf("  📂 Dir:  %s\n", cmd.Dir)
		fmt.Printf("  🔧 Cmd:  %s\n", cmd.Command)
	}

	// Parse the command into executable and arguments
	parts := strings.Fields(cmd.Command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	execCmd := exec.Command(parts[0], parts[1:]...)
	execCmd.Dir = cmd.Dir
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

func main() {
	fmt.Println("🔍 Scanning for go:generate mockgen commands...")
	fmt.Println("=" + strings.Repeat("=", 59))

	// Define search patterns - look for specific files that typically contain mockgen
	patterns := []string{
		"service.go",
		"store.go",
		"cmdbclient.go",
		"bk_cmdb_mock.go",
		"endpoint.go",
		"endpoints.go",
		"transport.go",
		"common.go",
		"troubleshooting_task_scheduler.go",
		"scheduler_engine.go",
	}

	// Extract all mockgen commands
	commands, err := extractGoGenerateCommands(".", patterns)
	if err != nil {
		fmt.Printf("❌ Error scanning files: %v\n", err)
		os.Exit(1)
	}

	if len(commands) == 0 {
		fmt.Println("⚠️  No mockgen commands found!")
		os.Exit(0)
	}

	fmt.Printf("✅ Found %d mockgen commands\n\n", len(commands))

	// Execute each command
	successCount := 0
	failCount := 0
	skipCount := 0

	for i, cmd := range commands {
		fmt.Printf("[%d/%d] Generating mock for %s\n", i+1, len(commands), filepath.Base(cmd.File))

		verbose := true // Set to false to reduce output
		if err := executeCommand(cmd, verbose); err != nil {
			fmt.Printf("  ❌ Failed: %v\n", err)
			failCount++
		} else {
			fmt.Printf("  ✅ Success\n")
			successCount++
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 Generation Summary:")
	fmt.Printf("  Total:    %d\n", len(commands))
	fmt.Printf("  Success:  %d\n", successCount)
	fmt.Printf("  Failed:   %d\n", failCount)
	if skipCount > 0 {
		fmt.Printf("  Skipped:  %d\n", skipCount)
	}

	if failCount > 0 {
		fmt.Println("\n❌ Some mock generations failed!")
		os.Exit(1)
	}

	fmt.Println("\n✅ All mocks generated successfully!")
}
