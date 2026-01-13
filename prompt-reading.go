package main

import (
	"log/slog"
	"os"
)

func getBillReceptPrompt() []byte {
	promptFileName := "bill-prompt.txt"
	prompt, err := os.ReadFile(promptFileName)
	if err != nil {
		slog.Error("failed to read bill prompt file", "error", err)
		os.Exit(1)
	}
	return prompt
}

func readSplitConvoPrompt() []byte {
	splitConvoFileName := "rules-prompt.txt"
	splitConvoPromptBytes, err := os.ReadFile(splitConvoFileName)
	if err != nil {
		slog.Error("failed to read split convo file", "error", err)
		os.Exit(1)
	}
	return splitConvoPromptBytes
}

func getSplitConvoPrompt(billItems []BillItem, splitConvo string) string {
	splitConvoPrompt := string(readSplitConvoPrompt())
	splitConvoPrompt += getItemsListAsString(billItems)
	splitConvoPrompt += splitConvo

	return splitConvoPrompt
}
