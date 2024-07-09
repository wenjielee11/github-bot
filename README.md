
# JambuBot

## Overview
This repository contains a GitHub bot powered by JamAIBase, written in Go. This bot is designed to automate various GitHub repository tasks, including issue management, pull request handling, and more. This documentation provides a detailed guide on how to set up, configure, and use the bot, including an explanation of its main features and components.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Features](#features)
5. [Usage](#usage)
6. [Development](#development)
7. [Contributing](#contributing)

## Prerequisites
- Go (version 1.16 or later)
- A GitHub account
- GitHub App with necessary permissions
- JamAIBase account

## Installation
1. **Clone the Repository:**
You should clone this repository into your existing repository as a submodule, or simply copy it over. 
    ```sh
    git clone https://github.com/wenjielee1/github-bot.git
    mv github-bot path/to/your/repo
    ```
    or if you are adding it as a submodule, from the root folder,
   ```
   git submodule add https://github.com/wenjielee1/github-bot path/to/submodule/github-bot
   git submodule update --init --recursive
   ```

3. **Install Dependencies:**
    ```sh
    go mod tidy
    ```

4. **Build the Project:**
    ```sh
    go build
    ```

## Configuration
The bot requires several environment variables to function correctly. These can be set in a `.env` file or directly in your environment.

### Required Environment Variables:
- `TRIAGE_BOT_PRIVATE_KEY`: The private key for your GitHub App.
- `TRIAGE_BOT_APP_ID`: Your GitHub App ID.
- `TRIAGE_BOT_JAMAI_KEY`: The API key for JamAIBase.
- `TRIAGE_BOT_JAMAI_PROJECT_ID`: The project ID for JamAIBase.

### Installing JambuBot App
1. **Install the JambuBot App:**
   - Go to the GitHub Marketplace and install the [JambuBot app](https://github.com/marketplace/jambubot)  on your GitHub account.
   - Follow the instructions to install the app on your repositories.

2. **Retrieve App ID and Installation ID:**
   - After installing the app, note down the App ID and Installation ID from the app's settings page.

3. **Generate a Private Key:**
   - In the app's settings, generate a new private key. This will download a `.pem` file. Keep this file secure.

4. **Set Up Environment Variables:**
   - Set the required environment variables with the values obtained from the app installation and key generation.

5. **Setup workflow**
    - Copy the `github-bot/.github` folder, or if you already have existing workflows, `.github/workflows/github_bot.yml` into your project root.
    - You wil have to edit `github_bot.yml` to correctly navigate through your repository structure. You may do this under job `Build and run Go Script` 
        
## Features
The bot includes several key features, each powered by JamAIBase:

### 1. Issue Handling
- **Comment on Issues:** Adds comments to issues when specific events occur. JamAIBase's built-in RAG and chunk reranking ensure comments are relevant and context-aware.
- **Issue Labels:** Automatically suggests a label and adds them for you, zero code needed. 

### 2. Pull Request Handling
- **Review Pull Requests:** Automatically reviews pull requests for certain conditions, such as missing CHANGELOG updates or potential secret key leaks. JamAIBase powers the analysis by providing high accuracy checks and generating insightful feedback.
- **Suggest Labels:** Automatically suggests labels for new pull requests. Leveraging JamAIBase's advanced AI capabilities, the bot can suggest the most appropriate labels based on pull request content.

### 3. Event Handling
- **Workflow Events:** Handles various GitHub action events, such as `issues`, `pull_request`, and `push`. JamAIBase processes the events, ensuring that the bot's responses are timely and accurate.

## Usage
1. **Run the Bot:**
    ```sh
    ./github-bot
    ```

2. **Set Up Env Configs:**
   - Configure your GitHub repository secrets to match the defined env variables.

## Testing
To test the bot, you need to set the environment variables directly inside the test scripts under `github-bot/test` and generate a `.pem` file for the private key.

1. **Set Environment Variables:**
   - Edit the test scripts to include the values for `GITHUB_APP_ID`, `GITHUB_INSTALLATION_ID`, `GITHUB_APP_PRIVATE_KEY`, and other required environment variables.

2. **Generate Private Key:**
   - Generate a `.pem` file after installing the JambuBot app to get the private key. Use this file in your test scripts.

3. **Run Test Scripts:**
   - Execute the test scripts to verify the bot's functionality. For example:
     ```sh
     go run create.go
     go run delete.go
     ```

## Development
To contribute to the development of this bot, follow these steps:

1. **Fork the Repository:**
    - Create a fork of the repository on GitHub.

2. **Create a New Branch:**
    ```sh
    git checkout -b feature-branch
    ```

3. **Make Your Changes:**
    - Implement your changes and commit them with meaningful commit messages.

4. **Submit a Pull Request:**
    - Push your changes to your fork and submit a pull request to the main repository.

## Contributing
We welcome contributions to improve the bot. Please ensure that your contributions adhere to the project's coding standards and conventions. Before submitting a pull request, make sure to:

- Run tests to ensure your changes do not break existing functionality.
- Update documentation if your changes introduce new features or modify existing ones.

---

## File Descriptions

### `main.go`
The main entry point of the application. It initializes the server and sets up the routes for handling GitHub webhook events.

### `eventHandler.go`
Contains the event handling logic for different GitHub webhook events. It routes the events to the appropriate handlers based on the event type.

### `issueHandler.go`
Defines functions for handling issue-related events, such as creating, commenting on, and closing issues.

### `prHandler.go`
Contains logic for handling pull request events, including reviewing pull requests and suggesting labels.

### `githubModels.go`
Defines data structures for GitHub-related entities, such as issues and pull requests.

### `jamaiModels.go`
Defines data structures for JamAIBase-related entities.

### `authService.go`
Provides authentication services, including retrieving installation tokens for the GitHub App and JamAIBase headers.

### `issueService.go`
Implements the business logic for issue handling, such as creating and commenting on issues.

### `jamaiService.go`
Handles interactions with JamAIBase, such as sending requests and processing responses.

### `prService.go`
Implements the business logic for pull request handling, such as reviewing commits for potential issues and suggesting labels.

This documentation provides a comprehensive guide to understanding, setting up, and contributing to the `github-bot`. For detailed implementation, refer to the source code in the respective files.
