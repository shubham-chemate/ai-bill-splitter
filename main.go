package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
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

	chat, err := client.Chats.Create(ctx, "gemini-2.5-flash", nil, nil)
	if err != nil {
		slog.Error("failed to create chat", "error", err)
		os.Exit(1)
	}

	result, err := chat.SendMessage(ctx, genai.Part{Text: "how's the weather in pune today?"})
	if err != nil {
		slog.Error("error while chatting", "error", err)
	}

	debugPrint(result)

	result, err = chat.SendMessage(ctx, genai.Part{Text: "it's feeling cold here? for much more days we should expect low temperature?"})
	if err != nil {
		slog.Error("error while chatting", "error", err)
	}

	debugPrint(result)
}

func debugPrint[T any](r *T) {

	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response))
}
