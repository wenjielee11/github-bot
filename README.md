
# GitHub Bot

## Overview
The `github-bot` repository contains a GitHub bot powered by JamAIBase, written in Go. This bot is designed to automate various GitHub repository tasks, including issue management, pull request handling, and more. This documentation provides a detailed guide on how to set up, configure, and use the bot, including an explanation of its main features and components.

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
    ```sh
    git clone https://github.com/wenjielee1/github-bot.git
    cd github-bot
    ```

2. **Install Dependencies:**
    ```sh
    go mod tidy
    ```

3. **Build the Project:**
    ```sh
    go build
    ```

## Configuration
The bot requires several environment variables to function correctly. These can be set in a `.env` file or directly in your environment.

### Required Environment Variables:
- `GITHUB_APP_ID`: Your GitHub App ID.
- `GITHUB_APP_PRIVATE_KEY`: The private key for your GitHub App.
- `GITHUB_WEBHOOK_SECRET`: The secret used to validate incoming webhooks.
- `TRIAGE_BOT_JAMAI_KEY`: The API key for JamAIBase.
- `TRIAGE_BOT_JAMAI_PROJECT_ID`: The project ID for JamAIBase.

## Features
The bot includes several key features:

### 1. Issue Handling
- **Create and Close Issues:** Automatically creates and closes issues based on predefined conditions.
- **Comment on Issues:** Adds comments to issues when specific events occur.

### 2. Pull Request Handling
- **Review Pull Requests:** Automatically reviews pull requests for certain conditions, such as missing CHANGELOG updates or potential secret key leaks.
- **Suggest Labels:** Automatically suggests labels for new pull requests.

### 3. Event Handling
- **Webhook Events:** Handles various GitHub webhook events, such as `issues`, `pull_request`, and `push`.

## Usage
1. **Run the Bot:**
    ```sh
    ./github-bot
    ```

2. **Set Up Webhooks:**
   - Configure your GitHub repository to send webhooks to the URL where the bot is running.

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
