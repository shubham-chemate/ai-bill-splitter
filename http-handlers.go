package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
)

func parseMultipartRequest(r *http.Request) ([]byte, string, string, error) {
	// max 7 MB allowed in memory, others are in disk
	err := r.ParseMultipartForm(MaxImageSizeBytes)
	if err != nil {
		return nil, "", "", fmt.Errorf("error while parsing the form, %v", err)
	}

	billReceiptImageFile, fileHeader, err := r.FormFile("bill-image")
	if err != nil {
		return nil, "", "", fmt.Errorf("error while parsing bill image file, %v", err)
	}
	defer billReceiptImageFile.Close()

	if fileHeader.Size > MaxImageSizeBytes {
		return nil, "", "", fmt.Errorf("bill image exceeds allowed file size")
	}

	billReceipt, err := io.ReadAll(billReceiptImageFile)
	if err != nil {
		return nil, "", "", fmt.Errorf("error while reading bill image file, %v", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		slog.Info("mimetype must be present, it should be of image", "mimeType received is", mimeType)
		return nil, "", "", fmt.Errorf("only image is allowed, mime type unknown")
	}

	splitConvo := r.FormValue("split-rules")
	if splitConvo == "" {
		return nil, "", "", fmt.Errorf("split rules cannot be empty")
	}

	return billReceipt, mimeType, splitConvo, nil
}

func handleBillSplitRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Invalid HTTP method", nil)
		return
	}

	// we won't be reading files more than 7 MB from client
	r.Body = http.MaxBytesReader(w, r.Body, MaxImageSizeBytes)

	imageData, mimeType, splitConvo, err := parseMultipartRequest(r)
	if err != nil {
		slog.Error("error while parsing multipart form", "error", err)
		respondError(w, http.StatusBadRequest, "Failed to parse the request", "Invalid form data or file format")
		return
	}

	slog.Info("bill-image received", "number of bytes", len(imageData))
	slog.Info("split-rules received", "split-rules", splitConvo)

	personsSplit, err := processBill(imageData, splitConvo, mimeType)
	if err != nil {
		slog.Info("error getting persons split", "error", err)
		respondError(w, http.StatusInternalServerError, "error getting person splits", err.Error())
		return
	}

	printBill(personsSplit)

	respondSuccess(w, http.StatusOK, "bill split calculated", personsSplit)
}

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, statusCode int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(StandardResponse{
		Status:  "error",
		Message: message,
		Error:   details,
	})
}
