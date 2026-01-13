package main

import (
	"context"
	"strings"

	"google.golang.org/genai"
)

func extractText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return ""
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(string(part.Text))
	}
	return sb.String()
}

func queryModelForBillReceipt(client *genai.Client, billReceipt []byte, mimeType string, prompt []byte) (string, error) {
	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: string(prompt)},
					{
						InlineData: &genai.Blob{
							Data:     billReceipt,
							MIMEType: mimeType,
						},
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

func queryModelForSplitConvo(client *genai.Client, prompt string) (string, error) {
	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: string(prompt)},
				},
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}
