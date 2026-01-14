package main

import (
	"fmt"
)

// calculatedTotal = pricePerUnit * quantity + tax
// if quantity is not present or price per unit not present we will ignore validation
func validateBillItems(billItems []BillItem) error {
	for _, billItem := range billItems {
		if billItem.TotalPrice == -1 {
			return fmt.Errorf("total item price is not present, billItem: %+v", billItem)
		}

		itemTotal := billItem.TotalPrice
		calculatedTotal := billItem.Tax
		if billItem.PricePerUnit != -1 && billItem.Quantity != -1 {
			calculatedTotal += billItem.PricePerUnit * float64(billItem.Quantity)
		} else {
			continue
		}
		if itemTotal != calculatedTotal {
			return fmt.Errorf("Item Total not matching calculated total, billItem: %+v", billItem)
		}
	}
	return nil
}

// 1. total item share among all friends must be 0.99 to 1.01
// 2. item from items split must be in bill item (llm hellusination)
func validateItemsSplit(billItems []BillItem, itemsSplit []ItemSplit) error {
	billItemsSet := make(map[string]struct{})
	for _, billItem := range billItems {
		billItemsSet[billItem.ItemName] = struct{}{}
	}

	splitItemsSet := make(map[string]struct{})
	for _, splitItem := range itemsSplit {
		splitItemsSet[splitItem.Name] = struct{}{}

		// split share amoung all friends should be nearly 1
		sum := 0.0
		for _, split := range splitItem.PersonSplits {
			sum += split.Share
		}
		if sum > 1.010 || sum < 0.990 {
			return fmt.Errorf("item split is invalid, sum is not in [0.990, 1.010], item: %v, itemsplit: %v", splitItem.Name, splitItem.PersonSplits)
		}
	}

	if len(billItems) != len(splitItemsSet) {
		return fmt.Errorf("different number of items in bill and item split")
	}

	for billItem := range billItemsSet {
		if _, found := splitItemsSet[billItem]; !found {
			return fmt.Errorf("different item names in bill and split, bill-item: %v, split-item: %v", billItem, splitItemsSet[billItem])
		}
	}

	return nil
}
