package main

import (
	"log/slog"
	"mime"
	"os"
	"path/filepath"
)

func getBillImage() ([]byte, string) {
	billFileName := "furniture-bill.jpg"
	billImage, err := os.ReadFile(billFileName)
	if err != nil {
		slog.Error("failed to read image file", "error", err)
		os.Exit(1)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(billFileName))
	if mimeType == "" {
		slog.Error("could not determine MIME type from file extension", "filename", billFileName)
		os.Exit(1)
	}

	return billImage, mimeType
}

func getSplitConvo() string {
	return `Akash and Amey buy Office Chair
			Dipti buy queen size bed
			Aditya, Suyog and Viraj buys recliner
			Bookshelf is shared among everyone`
}
