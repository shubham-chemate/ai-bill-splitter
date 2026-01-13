package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var apiKey string

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
		os.Exit(1)
	}

	apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info("gemini api key loaded", "length", len(apiKey))
	if len(apiKey) == 0 || apiKey == "" {
		slog.Warn("invalid api key", "api-key", apiKey)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})
	http.HandleFunc("/split", handleBillSplitRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
