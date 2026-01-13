package main

import (
	"fmt"
	"sort"
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

func validateItemsSplit(billItems []BillItem, itemsSplit []ItemSplit) error {
	itemList := []string{}
	for _, billItem := range billItems {
		itemList = append(itemList, billItem.ItemName)
	}

	itemSplitItems := []string{}
	for _, splitItem := range itemsSplit {
		itemSplitItems = append(itemSplitItems, splitItem.ItemName)

		// split share amoung all friends should be nearly 1
		sum := 0.0
		for _, split := range splitItem.Splits {
			sum += split.PersonShare
		}
		if sum > 1.010 || sum < 0.990 {
			return fmt.Errorf("item split is invalid, sum is not in [0.990, 1.010], item: %v, itemsplit: %v", splitItem.ItemName, splitItem.Splits)
		}
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
