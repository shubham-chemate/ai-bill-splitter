package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
)

// 7 MB
const MAX_IMAGE_SIZE = 7 * (1 << 20)

func handleBill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	// max 10 MB allowed in memory, others are in disk
	err := r.ParseMultipartForm(10 * (1 << 20))
	if err != nil {
		slog.Info("error while parsing multipart form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("bill-image")
	if err != nil {
		slog.Info("error while getting image of bill")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Size > MAX_IMAGE_SIZE {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("image exceeds allowed size limit", "size", fileHeader.Size, "limit", MAX_IMAGE_SIZE)
		return
	}

	billImgFileData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		slog.Info("mimetype must be present, it should be of image", "mimeType received is", mimeType)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	rules := r.FormValue("split-rules")
	if rules == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("rules cannot be empty")
		return
	}

	slog.Info("billImage received", "size", len(billImgFileData))
	slog.Info("split rules received", "split-rules", rules)

	response := map[string]interface{}{
		"status":  "success",
		"message": "bill split calculated",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	http.HandleFunc("/split", handleBill)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
