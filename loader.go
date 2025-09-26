// agent/loader.go

package agent

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const agentDir = "./agents"

// Agent represents a loaded agent file.
type Agent struct {
	Name   string
	Path   string
	Prompt string
}

// SelectAgent scans the agent directory, prompts the user to select an agent,
// and returns the selected agent's system prompt.
func SelectAgent() (string, error) {
	agents, err := findAgents()
	if err != nil {
		return "", fmt.Errorf("error finding agents: %w", err)
	}

	if len(agents) == 0 {
		return "", fmt.Errorf("no .agent files found in the '%s' directory", agentDir)
	}

	fmt.Println("Please select an agent:")
	for i, agent := range agents {
		// Display index 1-based for user-friendliness
		fmt.Printf("  %d: %s\n", i+1, agent.Name)
	}

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter selection: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		selection, err := strconv.Atoi(input)

		// Validate selection
		if err == nil && selection > 0 && selection <= len(agents) {
			// Read the selected agent's prompt
			selectedAgent := agents[selection-1]
			promptBytes, err := os.ReadFile(selectedAgent.Path)
			if err != nil {
				return "", fmt.Errorf("failed to read agent file '%s': %w", selectedAgent.Path, err)
			}
			fmt.Printf("\nAgent '%s' selected.\n", selectedAgent.Name)
			return string(promptBytes), nil
		}

		fmt.Println("Invalid selection. Please try again.")
	}
}

// findAgents scans the agent directory for files ending in .agent.
func findAgents() ([]Agent, error) {
	var agents []Agent

	entries, err := os.ReadDir(agentDir)
	if err != nil {
		// If the directory doesn't exist, create it.
		if os.IsNotExist(err) {
			fmt.Printf("Creating agent directory at: %s\n", agentDir)
			if err := os.Mkdir(agentDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create agent directory: %w", err)
			}
			return agents, nil // Return empty slice
		}
		return nil, fmt.Errorf("could not read agent directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".agent") {
			agentName := strings.TrimSuffix(entry.Name(), ".agent")
			agents = append(agents, Agent{
				Name: agentName,
				Path: filepath.Join(agentDir, entry.Name()),
			})
		}
	}

	return agents, nil
}
