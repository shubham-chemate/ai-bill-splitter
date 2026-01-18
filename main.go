package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	MaxImageSizeMb    = 7
	MaxImageSizeBytes = MaxImageSizeMb * 1024 * 1024

	DefaultPort = ":8080"
	GeminiModel = "gemini-2.5-flash-lite"

	SplitTolerance = 0.01
)

var apiKey string

func loadAPIKey() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
	}

	apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info("gemini api key loaded", "length", len(apiKey))
	if len(apiKey) == 0 || apiKey == "" {
		slog.Warn("invalid api key", "api-key", apiKey)
		os.Exit(1)
	}
}

func main() {

	loadAPIKey()

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})
	http.HandleFunc("/split", handleBillSplitRequest)

	slog.Info("server starting at", "port", DefaultPort)

	log.Fatal(http.ListenAndServe(DefaultPort, nil))
}
