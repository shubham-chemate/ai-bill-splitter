package main

type BillItem struct {
	ItemName     string   `json:"itemName"`
	PricePerUnit float64  `json:"pricePerUnit"`
	Quantity     int      `json:"quantity"`
	Tax          float64  `json:"tax"`
	TotalPrice   float64  `json:"totalPrice"`
	Warnings     []string `json:"warnings"`
}

type BillItems []BillItem

type Split struct {
	PersonName  string  `json:"personName"`
	PersonShare float64 `json:"personShare"`
}

type ItemSplits struct {
	ItemName string  `json:"itemName"`
	Splits   []Split `json:"splits"`
}

type ItemsSplit []ItemSplits
