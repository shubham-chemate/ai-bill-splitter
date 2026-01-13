package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"google.golang.org/genai"
)

func cleanRawJson(rawJson string) string {
	rawJson = strings.TrimSpace(rawJson)
	rawJson = strings.TrimPrefix(rawJson, "```json")
	rawJson = strings.TrimSuffix(rawJson, "```")
	rawJson = strings.TrimSpace(rawJson)
	rawJson = strings.TrimPrefix(rawJson, "\ufeff")

	return rawJson
}

func getBillItems(billReceipt []byte, mimeType string) []BillItem {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}

	prompt := getBillReceptPrompt()

	billItemsRawResp, err := extractItemsFromImage(client, billReceipt, mimeType, prompt)
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		os.Exit(1)
	}

	billItemsRawResp = cleanRawJson(billItemsRawResp)

	var billItems []BillItem
	err = json.Unmarshal([]byte(billItemsRawResp), &billItems)
	if err != nil {
		slog.Info("error while parsing billRawJson", "error", err, "billRawJson", billItemsRawResp)
		os.Exit(1)
	}

	err = validateBillItems(billItems)
	if err != nil {
		slog.Info("billItems validation failed", "error", err)
		os.Exit(1)
	}

	return billItems
}

func getItemsListAsString(billItems []BillItem) string {
	itemsString := "the items are - "
	for _, billItem := range billItems {
		itemsString += billItem.ItemName
		itemsString += ", "
	}
	itemsString += "\n"
	return itemsString
}

func getItemsSplit(billItems []BillItem, splitConvo string) []ItemSplit {

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}

	splitConvoPrompt := getSplitConvoPrompt(billItems, splitConvo)

	splitConvoRawResp, err := generateBillSplitRules(client, splitConvoPrompt)
	if err != nil {
		slog.Error("failed to get split convo json", "error", err)
		os.Exit(1)
	}

	splitConvoRawResp = cleanRawJson(splitConvoRawResp)

	var itemsSplit []ItemSplit
	err = json.Unmarshal([]byte(splitConvoRawResp), &itemsSplit)
	if err != nil {
		slog.Info("error while parsing splitConvRawJson", "error", err, "splitConvoRawJson", splitConvoRawResp)
		os.Exit(1)
	}

	err = validateItemsSplit(billItems, itemsSplit)
	if err != nil {
		slog.Error("items split validatation failed", "error", err)
		os.Exit(1)
	}

	return itemsSplit
}

// processing bill and conversation to generate persons split
func processBill(billRecept []byte, splitConvo string, mimeType string) ([]PersonSplit, error) {

	billItems := getBillItems(billRecept, mimeType)
	itemsSplit := getItemsSplit(billItems, splitConvo)

	return calculatePersonsSplit(billItems, itemsSplit)
}
