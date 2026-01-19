package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

func getBillItems(billReceipt []byte, mimeType string) ([]BillItem, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		return nil, err
	}

	prompt := getBillReceptPrompt()

	billItemsRawResp, err := extractItemsFromImage(client, billReceipt, mimeType, prompt)
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		return nil, err
	}

	billItemsRawResp = cleanRawJson(billItemsRawResp)

	var billItems []BillItem
	err = json.Unmarshal([]byte(billItemsRawResp), &billItems)
	if err != nil {
		slog.Info("error while parsing billRawJson", "error", err, "billRawJson", billItemsRawResp)
		return nil, err
	}

	slog.Info("received bill items from llm", "bill items", billItems)

	err = validateBillItems(billItems)
	if err != nil {
		slog.Info("billItems validation failed", "error", err)
		return nil, err
	}

	return billItems, nil
}

func getItemsListAsString(billItems []BillItem) string {
	itemsString := "item list:\n"
	for i, billItem := range billItems {
		itemsString += fmt.Sprintf("%d. %s\n", i+1, billItem.ItemName)
	}
	return itemsString
}

func getItemsSplit(billItems []BillItem, splitRules string) ([]ItemSplit, error) {

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		return nil, err
	}

	splitRulesPrompt := getSplitRulesPrompt(billItems, splitRules)

	splitRulesRawResp, err := generateBillSplitRules(client, splitRulesPrompt)
	if err != nil {
		slog.Error("failed to get split rules json", "error", err)
		return nil, err
	}

	splitRulesRawResp = cleanRawJson(splitRulesRawResp)

	var itemsSplit []ItemSplit
	err = json.Unmarshal([]byte(splitRulesRawResp), &itemsSplit)
	if err != nil {
		slog.Info("error while parsing splitConvRawJson", "error", err, "splitRulesRawJson", splitRulesRawResp)
		return nil, err
	}

	slog.Info("received item split from llm", "items split", itemsSplit)

	err = validateItemsSplit(billItems, itemsSplit)
	if err != nil {
		slog.Error("items split validatation failed", "error", err)
		return nil, err
	}

	return itemsSplit, nil
}

// processing bill and conversation to generate persons split
func processBill(billRecept []byte, splitRules string, mimeType string) ([]PersonSplit, error) {

	billItems, err := getBillItems(billRecept, mimeType)
	if err != nil {
		return nil, err
	}

	itemsSplit, err := getItemsSplit(billItems, splitRules)
	if err != nil {
		return nil, err
	}

	return calculatePersonsSplit(billItems, itemsSplit)
}
