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

// 7 MB
const MAX_IMAGE_SIZE = 7 * (1 << 20)

func handleBillSplitRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	// we won't be reading files more than 7 MB from client
	r.Body = http.MaxBytesReader(w, r.Body, MAX_IMAGE_SIZE)

	// max 7 MB allowed in memory, others are in disk
	err := r.ParseMultipartForm(MAX_IMAGE_SIZE)
	if err != nil {
		slog.Info("error while parsing multipart form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	billReceiptImageFile, fileHeader, err := r.FormFile("bill-image")
	if err != nil {
		slog.Info("error while getting image of bill")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer billReceiptImageFile.Close()

	if fileHeader.Size > MAX_IMAGE_SIZE {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("image exceeds allowed size limit", "size", fileHeader.Size, "limit", MAX_IMAGE_SIZE)
		return
	}

	billReceipt, err := io.ReadAll(billReceiptImageFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		slog.Info("mimetype must be present, it should be of image", "mimeType received is", mimeType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	splitConvo := r.FormValue("split-rules")
	if splitConvo == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("rules cannot be empty")
		return
	}

	slog.Info("billImage received", "number of bytes", len(billReceipt))
	slog.Info("split rules received", "split-rules", splitConvo)

	// *************** PROCESS BILL IMAGE & CONVERSATION TO CALCULATE BILL ***************

	personsSplit, err := processBill(billReceipt, splitConvo, mimeType)
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
