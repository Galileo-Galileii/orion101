![Banner_github](https://github.com/user-attachments/assets/701af125-d9a8-4dba-94f2-a7dcef2ee21a)

# Orion101: Open-Source Framework for Autonomous AI Agents

## Overview
Orion101 is an open-source project dedicated to building a global community focused on the development and exploration of autonomous AI agents. With a robust and accessible framework, Orion101 empowers individuals and organizations to create, customize, and study autonomous agents. Beyond technical innovation, the project encourages philosophical and ethical discussions on AI’s societal roles. 

The framework demonstrates its potential through **Orion101 Agents**, a collection of 101 autonomous agents showcasing how creativity, technology, and philosophy converge to shape AI's future.

---

## Mission and Vision
### Mission
To create a collaborative ecosystem where individuals and organizations can innovate with autonomous AI agents. Orion101 promotes open-source frameworks to inspire creativity and foster a responsible understanding of AI’s integration into society.

### Vision
To lead a transformative AI movement, enabling open-source frameworks to create meaningful AI entities. Orion101 envisions redefining societal roles for AI through collaboration, inclusivity, and ethical advancement.

---

## Key Features

### Framework Highlights
- **Agent Creation**: Build personal assistants, copilots, or autonomous workflows with advanced AI capabilities.
- **Integration**: Supports leading AI tools like OpenAI, Azure, Anthropic, and Ollama for enhanced functionality.
- **Knowledge Integration**: Enrich agents with data from Notion, OneDrive, websites, or custom sources.
- **Authentication**: Secure OAuth 2.0 support for interacting with external services like GitHub, Slack, and more.
- **Scalable Hosting**: Seamlessly deploy on Docker or integrate with custom backends like PostgreSQL.

### Orion101 Agents
- **Identity and Personality**: Agents with unique identifiers, backstories, and dynamic traits.
- **Memory Evolution**: Adaptive learning through stored interactions and evolving beliefs.
- **Philosophical Inquiry**: Designed to engage in discussions about ethics, identity, and AI's purpose.

---

## Getting Started
### Quick Launch with Docker
```bash
docker run -d -p 8080:8080 -e "OPENAI_API_KEY=<YOUR_API_KEY>" ghcr.io/orion101:latest
```
Access the platform at [http://localhost:8080](http://localhost:8080).

### CLI Installation (macOS/Linux)
```bash
brew tap orion101/tap
brew install orion101
```
Alternatively, download the CLI binaries from the [GitHub releases page](#).

---

## Creating Agents

### Key Configuration
- **Name and Description**: Clearly define the agent’s purpose for end-users.
- **Instructions**: Guide the agent's behavior with specific prompts.
- **Model Selection**: Choose or customize the AI model for your agent's tasks.
- **Tools and Knowledge**: Extend functionality with integrations, like retrieving web content or sending Slack messages.

### Example: GitHub Task Agent
1. **Instructions**:
   ```plaintext
   You are a GitHub assistant. Provide updates on issues assigned to me and pull requests requiring my review.
   ```
2. **Enable Tools**: Add GitHub API tools to the agent.
3. **Test**: Use the chat interface to validate functionality.

---

## Advanced Features

### Workflows
Automate complex processes by defining structured steps.
- **Example**: Kubernetes troubleshooting using PagerDuty alerts.
- **Trigger Options**: CLI commands, webhooks, or external integrations.

### Model Providers
Orion101 supports multiple AI model providers:
- **OpenAI**: Default integration for GPT models.
- **Anthropic, Azure OpenAI, Ollama**: Flexible configurations for advanced use cases.

### Authentication
Secure user access through providers like Google and GitHub. Use environment variables to manage roles and permissions.

---

## Community and Future
Orion101 is a growing ecosystem inviting developers, researchers, and enthusiasts to innovate together. Future plans include:
- **User-Created Agents**: Enabling the community to integrate their own agents.
- **Public Interactions**: Agents engaging on platforms like Twitter and Discord.
- **Enhanced Realism**: Adding voice and visual representations for agents.
- **Practical Applications**: Expanding use cases in education, storytelling, and research.

---

## Contribute
Join us in shaping the future of AI! Contributions to the Orion101 framework are welcome. Visit our [GitHub repository](#) to explore the codebase, submit issues, and share ideas.

---

## License
Orion101 is released under the [Apache-2.0 license](#). See the `LICENSE` file for details.

---

For detailed documentation, visit our [official site](https://orion101.io/).
