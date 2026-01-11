package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func queryModelForBill(client *genai.Client, billImage []byte, mimeType string, prompt []byte) (string, error) {
	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: string(prompt)},
					{
						InlineData: &genai.Blob{
							Data:     billImage,
							MIMEType: mimeType,
						},
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

func queryModelForRules(client *genai.Client, prompt string) (string, error) {
	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: string(prompt)},
				},
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

func extractText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return ""
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(string(part.Text))
	}
	return sb.String()
}

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

func validateBill(billItems BillItems) error {
	for _, billItem := range billItems {
		if billItem.TotalPrice == -1 {
			return fmt.Errorf("total item price is not present, billItem: %+v", billItem)
		}

		itemTotal := billItem.TotalPrice
		calculatedTotal := billItem.Tax
		if billItem.PricePerUnit != -1 && billItem.Quantity != -1 {
			calculatedTotal += billItem.PricePerUnit * float64(billItem.Quantity)
		}
		if itemTotal != calculatedTotal {
			return fmt.Errorf("item Total not matching calculated total, billItem: %+v", billItem)
		}
	}
	return nil
}

func getFriendsSplit(billItems BillItems, itemsSplit ItemsSplit) ([]PersonSplits, error) {
	if len(billItems) == 0 {
		return nil, fmt.Errorf("empty billItems")
	}
	if len(itemsSplit) == 0 {
		return nil, fmt.Errorf("empty items split")
	}
	return []PersonSplits{}, nil
}

func validateItemSplit(billItems BillItems, itemsSplit ItemsSplit) error {
	itemList := []string{}
	for _, billItem := range billItems {
		itemList = append(itemList, billItem.ItemName)
	}
	itemSplitItems := []string{}
	for _, splitItem := range itemsSplit {
		itemSplitItems = append(itemSplitItems, splitItem.ItemName)
	}

	if len(itemList) != len(itemSplitItems) {
		return fmt.Errorf("different number of items in bill and item split")
	}

	sort.Strings(itemList)
	sort.Strings(itemSplitItems)

	for i := range len(itemList) {
		if itemList[i] != itemSplitItems[i] {
			return fmt.Errorf("different item names in bill and split, billItems: %v, split items: %v", itemList, itemSplitItems)
		}
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found", "error", err)
		os.Exit(1)
	}

	var apiKey = os.Getenv("GEMINI_API_KEY")
	slog.Info("api key loaded", "length", len(apiKey))
	if len(apiKey) == 0 || apiKey == "" {
		slog.Warn("invalid api key", "api-key", apiKey)
		os.Exit(1)
	}

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

	promptFileName := "bill-prompt.txt"
	prompt, err := os.ReadFile(promptFileName)
	if err != nil {
		slog.Error("failed to read bill prompt file", "error", err)
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

	billRawJson, err := queryModelForBill(client, billImage, mimeType, prompt)
	if err != nil {
		slog.Error("failed to generate content", "error", err)
		os.Exit(1)
	}

	billRawJson = cleanRawJson(billRawJson)

	var itemsBill BillItems
	err = json.Unmarshal([]byte(billRawJson), &itemsBill)
	if err != nil {
		slog.Info("error while parsing billRawJson", "error", err, "billRawJson", billRawJson)
		os.Exit(1)
	}

	err = validateBill(itemsBill)
	if err != nil {
		slog.Info("bill validation failed", "error", err)
		os.Exit(1)
	}

	splitConvoFileName := "rules-prompt.txt"
	splitConvoPromptBytes, err := os.ReadFile(splitConvoFileName)
	if err != nil {
		slog.Error("failed to read split convo file", "error", err)
		os.Exit(1)
	}

	splitConvoPrompt := string(splitConvoPromptBytes)
	splitConvoPrompt += getItemsListAsString(itemsBill)
	splitConvoPrompt += `Akash and Amey buy Office Chair
						Dipti buy queen size bed
						Aditya, Suyog and Viraj buys recliner
						Bookshelf is shared among everyone`

	splitConvoRawJson, err := queryModelForRules(client, splitConvoPrompt)
	if err != nil {
		slog.Error("failed to get split convo json", "error", err)
		os.Exit(1)
	}

	splitConvoRawJson = cleanRawJson(splitConvoRawJson)

	var itemsSplit ItemsSplit
	err = json.Unmarshal([]byte(splitConvoRawJson), &itemsSplit)
	if err != nil {
		slog.Info("error while parsing splitConvRawJson", "error", err, "splitConvoRawJson", splitConvoRawJson)
		os.Exit(1)
	}

	err = validateItemSplit(itemsBill, itemsSplit)
	if err != nil {
		slog.Error("items validatation failed", "error", err)
		os.Exit(1)
	}

	for _, itemSplit := range itemsSplit {
		slog.Info("itemSplit of item", "itemSplit", itemSplit)
	}

	personSplits, err := getFriendsSplit(itemsBill, itemsSplit)
	if err != nil {
		slog.Info("error getting friends splits", "error", err)
		os.Exit(1)
	}

	for _, personSplit := range personSplits {
		slog.Info("person split received", "personSplit", personSplit)
	}

}
