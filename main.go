package main

import (
	"context"
	"log/slog"
	"os"

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

	imgData, err := os.ReadFile("sample-bill.jpg")
	if err != nil {
		slog.Error("failed to read image file", "error", err)
		os.Exit(1)
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{
						Text: "what is this bill about? can you list down items as well. if this is not a bill then you can say so!",
					},
					{
						InlineData: &genai.Blob{
							Data:     imgData,
							MIMEType: "image/jpeg",
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
