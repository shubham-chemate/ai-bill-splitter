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
		return nil, "", "", fmt.Errorf("failed while parsing the form, %w", err)
	}

	billReceiptImageFile, fileHeader, err := r.FormFile("bill-image")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed while parsing bill image file, %w", err)
	}
	defer billReceiptImageFile.Close()

	if fileHeader.Size > MaxImageSizeBytes {
		return nil, "", "", fmt.Errorf("bill image exceeds allowed file size")
	}

	billReceipt, err := io.ReadAll(billReceiptImageFile)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to read bill image file, %w", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		slog.Warn("mimetype must be present, it should be of image", "received mime type is", mimeType)
		return nil, "", "", fmt.Errorf("only image is allowed, mime type unknown")
	}

	splitRules := r.FormValue("split-rules")
	if splitRules == "" {
		return nil, "", "", fmt.Errorf("split rules cannot be empty")
	}

	return billReceipt, mimeType, splitRules, nil
}

func handleBillSplitRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Invalid HTTP method", nil)
		return
	}

	// we won't be reading files more than 7 MB from client
	r.Body = http.MaxBytesReader(w, r.Body, MaxImageSizeBytes)

	imageData, mimeType, splitRules, err := parseMultipartRequest(r)
	if err != nil {
		slog.Error("Failed to parse multipart form", "error", err)
		respondError(w, http.StatusBadRequest, "Unable to parse the request, invalid file format", nil)
		return
	}

	slog.Info("bill-image received", "number of bytes", len(imageData))
	slog.Info("split-rules received", "split-rules", splitRules)

	personsSplit, err := processBill(imageData, splitRules, mimeType)
	if err != nil {
		slog.Error("Failed to process bukk", "error", err)
		respondError(w, http.StatusInternalServerError, "Unable to process bill", nil)
		return
	}

	printBill(personsSplit)

	respondSuccess(w, http.StatusOK, "Bill split calculated", personsSplit)
}

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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

func respondError(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(StandardResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	})
}
