package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"ollama-vt100/agent"
	"ollama-vt100/config"
	"ollama-vt100/ollama"
	"ollama-vt100/ui"
)

// SessionConfig holds the configuration for a single chat session.
type SessionConfig struct {
	AgentSystemPrompt string
	Model             string
	Temperature       float32
	ContextWindow     int
}

func main() {
	// 1. Load startup configuration from file FIRST
	startupConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal Error: Could not load configuration. %v", err)
	}

	fmt.Println("Ollama VT-100 Client")
	fmt.Println("--------------------")
	fmt.Printf("Connecting to Ollama at: %s\n\n", startupConfig.OllamaURL)

	var sessionCfg SessionConfig

	// 2. Agent Selection
	agentPrompt, err := agent.SelectAgent()
	if err != nil {
		log.Fatalf("Failed to select agent: %v", err)
	}
	sessionCfg.AgentSystemPrompt = agentPrompt

	// 3. Model Selection (pass the loaded URL to the client)
	ollamaClient := ollama.NewClient(startupConfig.OllamaURL)
	model, err := selectModel(ollamaClient)
	if err != nil {
		log.Fatalf("Failed to select model: %v", err)
	}
	sessionCfg.Model = model

	// 4. Parameter Configuration
	temp, context, err := configureParameters()
	if err != nil {
		log.Fatalf("Failed to configure parameters: %v", err)
	}
	sessionCfg.Temperature = temp
	sessionCfg.ContextWindow = context

	log.Println("\nConfiguration complete:")
	log.Printf("  - Model: %s", sessionCfg.Model)
	log.Printf("  - Temperature: %.1f", sessionCfg.Temperature)
	log.Printf("  - Context Window: %d", sessionCfg.ContextWindow)

	// Updated startup message
	fmt.Println("\nStarting chat. Type '/quit' or '/exit' to end the session.")
	fmt.Println("-------------------------------------------------")

	chatLoop(sessionCfg, ollamaClient)
}

// selectModel fetches models and prompts the user for a selection.
func selectModel(client *ollama.Client) (string, error) {
	models, err := client.ListModels()
	if err != nil {
		return "", err
	}
	if len(models) == 0 {
		return "", fmt.Errorf("no local models found in Ollama")
	}

	fmt.Println("\nPlease select a model:")
	for i, modelName := range models {
		fmt.Printf("  %d: %s\n", i+1, modelName)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter selection: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		selection, err := strconv.Atoi(input)
		if err == nil && selection > 0 && selection <= len(models) {
			selectedModel := models[selection-1]
			fmt.Printf("Model '%s' selected.\n", selectedModel)
			return selectedModel, nil
		}
		fmt.Println("Invalid selection. Please try again.")
	}
}

// configureParameters prompts the user to set temperature and context window.
func configureParameters() (float32, int, error) {
	reader := bufio.NewReader(os.Stdin)

	// Temperature
	var temp float32
	defaultTemp := 0.7
	for {
		fmt.Printf("\nEnter temperature (default: %.1f): ", defaultTemp)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			temp = float32(defaultTemp)
			break
		}
		parsedTemp, err := strconv.ParseFloat(input, 32)
		if err == nil && parsedTemp >= 0 {
			temp = float32(parsedTemp)
			break
		}
		fmt.Println("Invalid input. Please enter a positive number.")
	}

	// Context Window
	var context int
	defaultContext := 128000
	for {
		fmt.Printf("Enter context window size (default: %d): ", defaultContext)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			context = defaultContext
			break
		}
		parsedContext, err := strconv.Atoi(input)
		if err == nil && parsedContext > 0 {
			context = parsedContext
			break
		}
		fmt.Println("Invalid input. Please enter a positive integer.")
	}

	return temp, context, nil
}

// chatLoop manages the main interactive chat.
func chatLoop(config SessionConfig, client *ollama.Client) {
	history := []ollama.Message{
		{Role: "system", Content: config.AgentSystemPrompt},
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			continue
		}
		// Added /quit as an exit option
		if userInput == "/exit" || userInput == "/quit" {
			fmt.Println("Exiting.")
			break
		}

		history = append(history, ollama.Message{Role: "user", Content: userInput})

		req := ollama.ChatRequest{
			Model:    config.Model,
			Messages: history,
		}
		req.Options.Temperature = config.Temperature
		req.Options.NumCtx = config.ContextWindow

		var responseBuilder strings.Builder
		fmt.Print("\nAI: ")

		// Initialize the Pager with screen dimensions and the input reader
		// Using 23 lines for a standard 24-line terminal leaves room for the next prompt.
		pager := ui.NewPager(76, 21, reader)

		err := client.Chat(req, func(chunk string) {
			sanitizedChunk := ui.ToASCII(chunk)
			pager.Write(sanitizedChunk)
			responseBuilder.WriteString(chunk)
		})

		pager.Flush()

		if err != nil {
			log.Printf("Error during chat: %v", err)
			history = history[:len(history)-1]
			continue
		}

		history = append(history, ollama.Message{Role: "assistant", Content: responseBuilder.String()})
	}
}
