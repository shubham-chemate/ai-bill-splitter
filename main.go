package main

import (
	"context"
	"log/slog"
	"mime"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
		os.Exit(1)
	}

	var apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info("api key loaded", "length", len(apiKey))

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}

	billFileName := "sample-bill.jpg"
	imgData, err := os.ReadFile(billFileName)
	if err != nil {
		slog.Error("failed to read image file", "error", err)
		os.Exit(1)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(billFileName))
	if mimeType == "" {
		slog.Error("could not determine MIME type from file extension", "filename", billFileName)
		os.Exit(1)
	}

	promptFileName := "bill-prompt.txt"
	textPrompt, err := os.ReadFile(promptFileName)
	if err != nil {
		slog.Error("failed to read prompt file", "error", err)
		os.Exit(1)
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{
						Text: string(textPrompt),
					},
					{
						InlineData: &genai.Blob{
							Data:     imgData,
							MIMEType: mimeType,
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		os.Exit(1)
	}

	text := ""
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			text += string(part.Text)
		}
	}

	slog.Info("received response from model", "text", text)
}
