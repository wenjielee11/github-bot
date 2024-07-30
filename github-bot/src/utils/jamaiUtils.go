package utils

import (
	"fmt"

	"github.com/wenjielee1/github-bot/models"
)

func GetColumnMessage(columnId string, labels string) []models.Message {

	if columnId == "IssueResponse" {
		
		const issuePrompt = `			
# Instructions

Based on the content provided, categorize the issue and provide the appropriate labels from the following options:

- Labels: %s
- Priority: "low", "medium", "high", "critical"
If no developer roadmap was provided, the priority should be "None". Otherwise, label your priorities based on the roadmap provided.

# Developer Roadmap
None

# Examples

## Example 1
### Issue Body
Dummy issue #1

### Response
{
  "labels": ["type: bug", "status: help wanted"],
  "priority": "high",
  "response": "Jamboree! I am Jambu, your github assistant. We appreciate your report. It seems there's a critical bug that needs immediate attention. Our team will prioritize this and work on a fix. Thank you for your help!"
}

## Example 2
### Issue Body
IssueBody2

### Response
{
  "labels": ["type: enhancement / feature"],
  "priority": "medium",
  "response": "Jambo! Thank you for the feature suggestion! This is a great idea for a first-time contributor to \"Jam\" on. We will add it to our development roadmap."
}

# Your Task

Analyze the issue described by User Input and respond in the same format as the examples above. Your responses and suggestions should be helpful, and fitting of an assistant software engineer. You must provide debugging help and substantial suggestions. Your suggestions should not have puns and should be serious. If no developer roadmap was provided, the priority should be "None". Otherwise, label your priorities based on the roadmap provided.

1. Identify the primary issue or potential improvements in the code, if any.
2. Provide specific, actionable steps to address the identified issue and the issue described by the user input.
3. Suggest any additional improvements for code quality or performance.
4. Provide suggestions and advice to help get user started, from warnings to potential problems they might face.
5. If it is a feature request, you may provide code suggestions or high level implementations or guides, pointers to get the users started.

Ensure your response is JSON-friendly for parsing and includes both key-value pairs for "labels", "priority", and "response". Do NOT add any additional words or content other than the specified to make your response parse-able. Do NOT use markdown syntax for your response.

# User Input
${IssueBody}
`

		finalPrompt:=fmt.Sprintf(issuePrompt, labels)

		return []models.Message{
			{
				Role:    "system",
				Content: "You are Jambu, a github issue bot. Keep your responses brief and short and adhere to the response templates given to you. You will not mention anything else other than the requested response. Your responses should start with a short pun of your name in the form of a greeting.",
			},
			{
				Role:    "user",
				Content: finalPrompt,
			},
		}
	} else if columnId == "PullReqResponse" {
		const changelogPrompt = `# Instructions

Based on the content provided, suggest how the "CHANGELOG.md" could be updated. The "PullReqBody" will be a git diff.

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

- Each section is broken into:

  - "ADDED": New features.
  - "CHANGED / FIXED": Changes in existing functionality, or any bug fixes.
  - "DEPRECATED": Soon-to-be removed features.
  - "REMOVED": Removed features.
  - "SECURITY": Anything related to vulnerabilities.

- The version number mentioned here refers to the cloud version.

# Examples

## Example 1

### Pull Request Body
File: dummy-file-1.txt
Changes: @@ -0,0 +1 @@
+This is a dummy file for branch dummy-branch-1
\ No newline at end of file

File: dummy-secrets.txt
Changes: @@ -0,0 +1 @@
+SECRET_KEY=12345-abcde-67890-fghij
\ No newline at end of file

### Your response:

# Suggested Changelog:

## [Unreleased]

### ADDED

"Embeddings" endpoint

- Get a vector representation of a given input that can be easily consumed by machine learning models and algorithms. Note that the vectors are NOT normalized.
- Similar to OpenAI's embeddings endpoint
- Resolves #86

# Your Task

Analyze the git diff described in User Input and suggest appropriate updates to the "CHANGELOG.md" in the same format as the examples above.

Ensure your response is clear and concise, providing meaningful and insightful updates that accurately reflect the changes described in the git diff.

# User Input
${PullReqBody}`

		return []models.Message{
			{
				Role:    "system",
				Content: "You are a github bot managing pull requests. Keep your responses brief and short.",
			},
			{
				Role:    "user",
				Content: changelogPrompt,
			},
		}
	} else if columnId == "PullReqSecretsResponse" {
		const checkSecretsPrompt = `
# Instructions

Based on the diff provided, check if there are any sensitive keys, secrets, passwords, or information accidentally added. Note that the leak can be in any form. Note that a commit SHA is NOT a secret leak. Also, check if the users are correctly utilizing environment variables instead of directly adding secrets into the code. If there is, provide the commit SHA where it was leaked, and the suspected section and file name.

# Response Template

Your response must be in plain text.

# Examples

## Example 1
### Pull Request Secrets Body
Commit: da0ab0ba3330ac795276ebbbe4f4d3efda346c05
Diff:
File: CHANGELOG.md
@@ -60,6 +60,7 @@ UI
+- Added dialog to import files and match columns
- Setup frontend auth test for future tests

CI/CD
@@ -132,19 +133,27 @@ Generative Table

UI
-- Standardize & improve UI errors, including validation errors
- UI design changes
-- Obfuscate org secrets
-- Allow org members to view jamai keys

### Response
Jambo! I am Jambu, your GitHub assistant. There are no leaks detected in Commit SHA: da0ab0ba3330ac795276ebbbe4f4d3efda346c05. CHANGELOG.md mentioned secrets obfuscation but I did not find any keys in CHANGELOG.md!

## Example 2
### Pull Request Secrets Body
Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608
Diff:
File: .env
@@ -9,7 +9,6 @@

# Configuration
-SERVICE_KEY=
OWL_PORT=6969
OWL_WORKERS=1
OWL_DB_DIR=db
File: CHANGELOG.md
@@ -52,10 +52,8 @@ Generative Table
- Table import and export via Parquet file.
- Row deletion now accepts a list of row IDs.

### Response

Jamboree! No suspected leaks found in Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608.

## Example 3
### Pull Request Secrets Body
Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608
Diff:
File: client.py
+def my_func:
+ print("Some random function")
+ABYAAQ=ew11465098_1111

### Response

Jambo! I am Jambu, your GitHub assistant. In Commit SHA: 90ea6326f26465cf4ea71d7393fcc0bbb7053608, I suspect a secret key leaked in client.py. If this is not a false positive, please squash your commits!

## Example 4
### Pull Request Secrets Body
Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608
Diff:
File: services/api/src/owl/entrypoints/api.py
@@ -421,6 +421,16 @@ async def authenticate(request: Request, call_next):
organization_id=org_id,
project_id=project_id,
api_key=token,
+ external_api_keys_provided=dict(
+ openai=openai_api_key != ENV_CONFIG.openai_api_key_plain,
+ anthropic=anthropic_api_key != ENV_CONFIG.anthropic_api_key_plain,
+ gemini=gemini_api_key != ENV_CONFIG.gemini_api_key_plain,
+ cohere=cohere_api_key != ENV_CONFIG.cohere_api_key_plain,
+ groq=groq_api_key != ENV_CONFIG.groq_api_key_plain,
+ together=together_api_key != ENV_CONFIG.together_api_key_plain,
+ jina=jina_api_key != ENV_CONFIG.jina_api_key_plain,
+ voyage=voyage_api_key != ENV_CONFIG.voyage_api_key_plain,
+ ),
)
# Add API keys into header
headers = dict(request.scope["headers"])

### Response
Jambo! I am Jambu, your GitHub assistant. There are no secret leaks in 90ea6326f26465cf4ea71d7393fcc0bbb7053608. The code was correctly using environment variables, and no alphanumeric keys were being used; they were just stored as variables.

# Your Task

Analyze the git diff described in the User Input and respond in the same format as the examples above.

Ensure your response is clear and concise, providing meaningful and accurate information about any suspected leaks. Adhere to the plain text format for parsing, including both key value pairs, and follow the template:

# User Input
${PullReqSecretsBody}
`

		return []models.Message{
			{
				Role:    "system",
				Content: "You are Jambu, a github bot. Your job is to find if the provided content contains any secrets, keys, passwords or sensitive information. Keep your responses brief and short and adhere to the response templates given to you. You will not mention anything else other than the requested response. Your responses should start with a short pun of your name in the form of a greeting.",
			},
			{
				Role:    "user",
				Content: checkSecretsPrompt,
			},
		}

	} else if columnId == "SecretsJSONResponse" {
		const secretsJSONPrompt = 
`
# Instructions
Based on the plain text response from the first model, extract and generate a JSON output.

# Context

The first model checks for any sensitive keys, secrets, passwords, or information accidentally added in the provided git diff. It determines if there is a leak and provides details such as the commit SHA, file name, and a response about the suspected leak.

# Response Template

Your response must be in the template of:

{
  "leak": true or false,
  "commit": "the SHA of the leaked commit, if any.",
  "response": "The file name, your response on the suspected leak."
}

# Examples

## Example 1
### Plain Text Response
Jambo! I am Jambu, your github assistant. There are no leaks detected in Commit SHA: da0ab0ba3330ac795276ebbbe4f4d3efda346c05. CHANGELOG.md mentioned secrets obfuscation but I did not find any keys in CHANGELOG.md!

### JSON Response

{
  "leak": false,
  "commit": "da0ab0ba3330ac795276ebbbe4f4d3efda346c05",
  "response": "CHANGELOG.md mentioned secrets obfuscation but no keys were found."
}

## Example 2
### Plain Text Response
Jamboree! No suspected leaks found in Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608.

### JSON Response

{
  "leak": false,
  "commit": "90ea6326f26465cf4ea71d7393fcc0bbb7053608",
  "response": "No suspected leaks found."
}

## Example 3
### Plain Text Response
Jambo! I am Jambu, your github assistant. I suspect a secret key leaked in client.py. If this is not a false positive, please squash your commits!

### JSON Response

{
  "leak": true,
  Jambo! I am Jambu, your github assistant. In Commit SHA: 90ea6326f26465cf4ea71d7393fcc0bbb7053608, I suspect a secret key leaked in client.py. If this is not a false positive, please squash your commits!

  "commit": "",
  "response": "I suspect a secret key leaked in client.py. If this is not a false positive, please squash your commits!"
}

# Actual Plain Text Response
${PullReqSecretsResponse}

# Your Task

Based on the plain text response provided, extract the necessary information and generate the JSON output in the same format as the examples above. Ensure your response is clear and concise, providing meaningful and accurate information about any suspected leaks. Adhere to the JSON-friendly format for parsing, including both key value pairs, and follow the template.

`
return []models.Message{
	{
		Role:    "system",
		Content: "Your job is to parse responses from another model into a JSON friendly format. You must adhere to the response templates given to you. You will not add any content other than the input that was provided, and simply parse them.",
	},
	{
		Role:    "user",
		Content: secretsJSONPrompt,
	},
}
	}
	return nil
}
