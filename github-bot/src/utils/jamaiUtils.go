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

# Developer Roadmap
## v0.3 (upcoming)
(Breaking) Add versioning into table metadata - [FEAT] [Backend] GenTable metadata should include a "version" str column #137
(Breaking) Add object type into gen_config - [FEAT] [Backend] GenTable generation config should include an "object" key to support more config object types #136
(Breaking) Change ChatRequest.stop default to None
Deprecate deploy in duplicate_table
Note that the above must be implemented in order to ease the implementation of the following:

Existing Template import into Discover Template (Backend) - [FEAT][Backend] Method to publish table (template) to "Discover Template"  #238
Discover Template import into project (UI) - [FEAT][UI] Discover Templates - Import Templates  #237
Table import and export - [FEAT] [UI] [Backend] Allow the import or export of projects and tables #124
Project import and export - [FEAT] [UI] [Backend] Allow the import or export of projects and tables #124
Extras

Billing refactor
New pricing model - [FEAT] [Backend] New pricing and billing model #235
Use ELLM API keys if user does not supply their own - [FEAT] [Backend] New pricing and billing model #235
Fix storage usage computation - [BUG] [Backend] Billing: Storage usage sometimes shows negative usage #234
Project rename
TS SDK: Update get_conversation_thread to accept 2 additional optional parameters for filtering by row_id and specifying whether to include the row in the returned thread.
Expose parent_id in table list method in Python and TS SDKs - [BUG] [Python SDK] Expose parent_id in table list method #252
Use one Infinity instance for multiple models - Heads up: Infinity:0.0.42 now supports multiple models JamAIBase#5
Automatic context length management: Reduce the number of chunks dynamically depending on context length

## v0.4
(Breaking) Remove deploy in duplicate_table
Template marketplace - [FEAT] Link Sharing and Table Duplication #211
Table sharing / iframe views - [FEAT] Link Sharing and Table Duplication #211
Role-based access control - [FEAT] [UI, Backend] Allow role-based access control #121

## Others
Janitor process - [FEAT] [Backend] Implement a separate janitor / maintenance process #233
Task queues for long-running operations such as row add
Future
(Breaking) Image input - [FEAT] [UI, Backend] Allow Image and File Inputs #120
(Breaking) File input - [FEAT] [UI, Backend] Allow Image and File Inputs #120
Flexible Knowledge Table and Chat Table schemas
KT: User defined embedding columns
CT: More than one chat thread
CT: Allow input interpolation in chat column
Unify all tables into one type

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

Analyze the issue described by User Input and respond in the same format as the examples above. Your responses and suggestions should be helpful, and fitting of an assistant software engineer. You must provide debugging help and substantial suggestions. Your suggestions should not have puns and should be serious.

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
		const checkSecretsPrompt = `# Instructions

Based on the diff provided, check if there are any sensitive keys, secrets, passwords, or information accidentally added. Note that the leak can be in any type of form. Note that a commit SHA is NOT a secret leak. Also check if the users are correctly utilizing environment variables instead of directly adding secretsinto the code.If there is, provide the commit SHA where it was leaked, and the suspected section and file name.

# Response Template

Your response must be in the template of:

{
  "leak": true or false,
  "commit": "the SHA of the leaked commit, if any.",
  "response": "The file name, your response on the suspected leak."
}

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

{
  "leak": false,
  "commit": "",
  "response": "Jambo! I am Jambu, your github assistant. There are no leaks in detected. CHANGELOG.md mentioned secrets obfuscation but I did not find any keys in CHANGELOG.md!"
}

## Example 2
### Pull Request Secrets Body
Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608
Diff:
File: .env
@@ -9,7 +9,6 @@ DOCIO_URL=http://docio:6979/api/docio
UNSTRUCTUREDIO_URL=http://unstructuredio:6989

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

{
  "leak": false,
  "commit": "",
  "response": "Jamboree! No suspected leaks found."
}

## Example 3
### Pull Request Secrets Body
Commit: 90ea6326f26465cf4ea71d7393fcc0bbb7053608
Diff:
File: client.py
+def my_func:
+ print("Some random function")
+ABYAAQ=ew11465098_1111

### Response
{
  "leak": true,
  "commit": "",
  "response": "Jambo! I am Jambu, your github assistant. I suspect a secret key leaked in client.py. If this is not a false positive, please squash your commits!"
}

# Your Task

Analyze the git diff described in the User Input and respond in the same format as the examples above.

Ensure your response is clear and concise, providing meaningful and accurate information about any suspected leaks. Adhere to the JSON-friendly format for parsing, inlcuding both key value pairs, and follow the template:

{
  "leak": true or false,
  "commit": "the SHA of the leaked commit, if any.",
  "response": "The file name, your response on the suspected leak."
}

# User Input
${PullReqSecretsBody}`

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

	}
	return nil
}
