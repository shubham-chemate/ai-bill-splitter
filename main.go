package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

func getItemsListAsString(billItems BillItems) string {
	itemsString := "the items are - "
	for _, billItem := range billItems {
		itemsString += billItem.ItemName
		itemsString += ", "
	}
	itemsString += "\n"
	return itemsString
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
		os.Exit(1)
	}

	var apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info("gemini api key loaded", "length", len(apiKey))
	if len(apiKey) == 0 || apiKey == "" {
		slog.Warn("invalid api key", "api-key", apiKey)
		os.Exit(1)
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}

	billImage, mimeType := getBillImage()
	prompt := getBillReceptPrompt()

	billItemsRawResp, err := queryModelForBillReceipt(client, billImage, mimeType, prompt)
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		os.Exit(1)
	}

	billItemsRawResp = cleanRawJson(billItemsRawResp)

	var billItems BillItems
	err = json.Unmarshal([]byte(billItemsRawResp), &billItems)
	if err != nil {
		slog.Info("error while parsing billRawJson", "error", err, "billRawJson", billItemsRawResp)
		os.Exit(1)
	}

	err = validateBillItems(billItems)
	if err != nil {
		slog.Info("bill validation failed", "error", err)
		os.Exit(1)
	}

	splitConvoPrompt := string(getSplitConvoPrompt())
	splitConvoPrompt += getItemsListAsString(billItems)
	splitConvoPrompt += getSplitConvo()

	splitConvoRawResp, err := queryModelForSplitConvo(client, splitConvoPrompt)
	if err != nil {
		slog.Error("failed to get split convo json", "error", err)
		os.Exit(1)
	}

	splitConvoRawResp = cleanRawJson(splitConvoRawResp)

	var itemsSplit ItemsSplit
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

	personSplits, err := getPersonsSplit(billItems, itemsSplit)
	if err != nil {
		slog.Info("error getting friends splits", "error", err)
		os.Exit(1)
	}

	fmt.Println("Bill Splitted as below:")
	for _, personSplit := range personSplits {
		// slog.Info("person split received", "personSplit", personSplit)

		fmt.Printf("(%s, Amt: %0.3f)\n", personSplit.PersonName, personSplit.TotalAmount)
		for _, items := range personSplit.SplitByItem {
			if items.Amount >= 0.01 {
				fmt.Printf("- item: %s, amount: %.3f\n", items.ItemName, items.Amount)
			}
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	http.HandleFunc("/split", handleBill)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
