Ollama VT-100 Client
A minimalist, terminal-based client for interacting with the Ollama API, designed specifically for low-bandwidth environments and classic terminals like the VT-100/VT-220. It is fully ASCII-compliant, supports custom "agents," and provides a more-style pager for long responses.

Features
Agent System: Load custom system prompts from .agent files.

VT-100 Compliant: Strictly ASCII output with 80-column word wrapping.

Configurable: Set the Ollama server URL, model, temperature, and context window at runtime.

Built-in Pager: Automatically pages through long responses, waiting for user input to continue.

Portable: Compiles to a single binary with no external dependencies.

Installation and Setup
To run the application, you need three components placed together in the same directory (e.g., ~/termapps/ on Linux):

The Binary: The compiled executable file for your operating system.

config.json: The configuration file for the Ollama server URL.

agents/ directory: A folder containing your custom agent files.

Your directory structure should look like this:

/path/to/your/app/
├── ollama-vt100-linux  (or your OS-specific binary)
├── config.json
└── agents/
    ├── researcher.agent
    └── coder.agent

For ease of use, add the directory containing the binary to your system's PATH.

Configuration
config.json
This file tells the client where to find your Ollama server. The application will automatically create a default version on its first run if one is not found.

To connect to a remote server, simply edit the ollama_url value.

Example config.json:

{
  "ollama_url": "[http://192.168.1.10:11434](http://192.168.1.10:11434)"
}

Creating Agent Files
Agents are simple text files that contain a system prompt. The client will find any file ending in .agent inside the agents/ directory and let you choose one at startup.

Create a new text file inside the agents/ directory.

The name of the file (before the extension) will be the agent's name (e.g., coder.agent will be listed as coder).

Write your system prompt directly into the file. There is no special formatting required.

Example agents/coder.agent:

You are an expert Go programmer. Provide clean, idiomatic, and well-commented code. Do not provide explanations unless explicitly asked. Focus on correctness and efficiency.

Usage
Run the binary from your terminal:

./ollama-vt100-linux

Follow the on-screen prompts to select an agent, a model, and to configure the temperature and context window (press Enter to accept defaults).

Type your message at the > prompt and press Enter.

If the AI's response is longer than the screen, a -- More -- prompt will appear. Press Enter to view the next page.

To end the session, type /quit or /exit and press Enter.
