package main

import (
	"embed"
	"log/slog"
	"os"
)

//go:embed prompts/*.txt
var promptFS embed.FS

func getBillReceptPrompt() []byte {
	promptFileName := "prompts/bill-prompt.txt"
	prompt, err := promptFS.ReadFile(promptFileName)
	if err != nil {
		slog.Error("failed to read bill prompt file", "error", err)
		os.Exit(1)
	}
	return prompt
}

func readSplitRulesPrompt() []byte {
	splitRulesFileName := "prompts/rules-prompt.txt"
	splitRulesPromptBytes, err := promptFS.ReadFile(splitRulesFileName)
	if err != nil {
		slog.Error("failed to read split rules file", "error", err)
		os.Exit(1)
	}
	return splitRulesPromptBytes
}

func getSplitRulesPrompt(billItems []BillItem, splitRules string) string {
	splitRulesPrompt := string(readSplitRulesPrompt())
	splitRulesPrompt += getItemsListAsString(billItems)
	splitRulesPrompt += splitRules

	return splitRulesPrompt
}
