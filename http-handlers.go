package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func parseMultipartRequest(r *http.Request) ([]byte, string, string, error) {
	// max 7 MB allowed in memory, others are in disk
	err := r.ParseMultipartForm(MaxImageSizeBytes)
	if err != nil {
		return nil, "", "", fmt.Errorf("error while parsing the form")
	}

	billReceiptImageFile, fileHeader, err := r.FormFile("bill-image")
	if err != nil {
		return nil, "", "", fmt.Errorf("error while parsing bill image file")
	}
	defer billReceiptImageFile.Close()

	if fileHeader.Size > MaxImageSizeBytes {
		return nil, "", "", fmt.Errorf("bill image exceeds allowed file size")
	}

	billReceipt, err := io.ReadAll(billReceiptImageFile)
	if err != nil {
		return nil, "", "", fmt.Errorf("error while reading bill image file")
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Error": "Invalid HTTP Method",
		})
		return
	}

	// we won't be reading files more than 7 MB from client
	r.Body = http.MaxBytesReader(w, r.Body, MaxImageSizeBytes)

	imageData, mimeType, splitConvo, err := parseMultipartRequest(r)
	if err != nil {
		slog.Error("error while parsing multipart form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("bill-image received", "number of bytes", len(imageData))
	slog.Info("split-rules received", "split-rules", splitConvo)

	// *************** PROCESS BILL IMAGE & CONVERSATION TO CALCULATE BILL ***************
	personsSplit, err := processBill(imageData, splitConvo, mimeType)
	if err != nil {
		slog.Info("error getting persons split", "error", err)
		os.Exit(1)
	}

	fmt.Println("Bill Splitted as below:")
	for _, personSplit := range personsSplit {
		fmt.Printf("(%s, Amt: %0.3f)\n", personSplit.PersonName, personSplit.TotalAmount)
		for _, items := range personSplit.SplitByItem {
			if items.Amount >= 0.01 {
				fmt.Printf("- item: %s, amount: %.3f\n", items.ItemName, items.Amount)
			}
		}
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "bill split calculated",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
