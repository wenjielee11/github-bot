package utils

import (
	"github.com/wenjielee1/github-bot/models"
)

func GetColumnMessage(columnId string) []models.Message {
	if columnId == "IssueResponse" {
		return []models.Message{
			{
				Role:    "system",
				Content: "You are a github issue bot. Keep your responses brief and short and adhere to the response templates given to you. You will not mention anything else other than the requested response.",
			},
			{
				Role:    "user",
				Content: "Based on ${IssueBody}, provide the labels with BUG, FEATURE, DOCUMENTATION, ENHANCEMENT, GOOD FIRST ISSUE, HELP WANTED, INVALID, QUESTION, the priority LOW, MEDIUM, HIGH, CRITICAL. Finally, respond to the issue by suggesting an improvement, provide some comments, or provide words of appreciation.\n{\nlabels: [your_labels],\npriority: your_priority,\nresponse: your_response\n}\nYour entire response must adhere to a JSON friendly format for parsing, that includes both key value pairs.",
			},
		}
	} else if columnId == "PullReqResponse" {
		return []models.Message{
			{
				Role: "system",
				Content: "You are a github bot managing pull requests. Keep your responses brief and short.",
			},
			{
				Role: "user",
				Content: "Based on ${PullReqBody}, suggest how the CHANGELOG.md could be updated. This is a test case, so you will be given dummy content. Regardless of the dummy content, I want you to roleplay as if you are reviewing an actual PR. Generate coherent content and responses.",
			},
		}
	}
	return nil
}
