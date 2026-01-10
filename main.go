package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, reading from environment variable")
	}

	var apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info(fmt.Sprintf("%d", len(apiKey)))

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("hey"),
		nil,
	)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info(result.Text())
}
