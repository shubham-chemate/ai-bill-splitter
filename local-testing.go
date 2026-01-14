package main

import "fmt"

func getSplitConvo() string {
	return `Akash and Amey buy Office Chair
			Dipti buy queen size bed
			Aditya, Suyog and Viraj buys recliner
			Bookshelf is shared among everyone`
}

func printBill(personsSplit []PersonSplit) {
	fmt.Println("Bill Splitted as below:")
	for _, personSplit := range personsSplit {
		fmt.Printf("(%s, Amt: %0.3f)\n", personSplit.PersonName, personSplit.TotalAmount)
		for _, items := range personSplit.SplitByItem {
			if items.Amount >= 0.01 {
				fmt.Printf("- item: %s, amount: %.3f\n", items.ItemName, items.Amount)
			}
		}
	}
}
