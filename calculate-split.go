package main

import "fmt"

func getPersonsSplit(billItems BillItems, itemsSplit ItemsSplit) ([]PersonSplits, error) {
	if len(billItems) == 0 {
		return nil, fmt.Errorf("empty billItems")
	}
	if len(itemsSplit) == 0 {
		return nil, fmt.Errorf("empty items split")
	}

	itemsPrices := make(map[string]float64)
	for _, billItem := range billItems {
		itemsPrices[billItem.ItemName] = billItem.TotalPrice
	}

	personsSplits := make(map[string][]SplitByItem)
	for _, itemSplit := range itemsSplit {
		itemName := itemSplit.ItemName
		itemPrice := itemsPrices[itemName]

		for _, split := range itemSplit.Splits {
			personName := split.PersonName
			personShare := split.PersonShare

			_, exist := personsSplits[personName]
			if !exist {
				personsSplits[personName] = make([]SplitByItem, 0)
			}

			splitForPerson := SplitByItem{
				ItemName: itemName,
				Amount:   personShare * itemPrice,
			}

			personsSplits[personName] = append(personsSplits[personName], splitForPerson)
		}
	}

	personsSplitsArray := make([]PersonSplits, 0)
	for personName, personSplits := range personsSplits {
		totalAmount := 0.0
		for _, split := range personSplits {
			totalAmount += split.Amount
		}
		record := PersonSplits{
			PersonName:  personName,
			SplitByItem: personSplits,
			TotalAmount: totalAmount,
		}
		personsSplitsArray = append(personsSplitsArray, record)
	}

	return personsSplitsArray, nil
}
