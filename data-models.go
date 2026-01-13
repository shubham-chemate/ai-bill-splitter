package main

type BillItem struct {
	ItemName     string   `json:"itemName"`
	PricePerUnit float64  `json:"pricePerUnit"`
	Quantity     int      `json:"quantity"`
	Tax          float64  `json:"tax"`
	TotalPrice   float64  `json:"totalPrice"`
	Warnings     []string `json:"warnings"`
}

type SplitByPerson struct {
	PersonName  string  `json:"personName"`
	PersonShare float64 `json:"personShare"`
}

type ItemSplit struct {
	ItemName string          `json:"itemName"`
	Splits   []SplitByPerson `json:"splits"`
}

type SplitByItem struct {
	ItemName string
	Amount   float64
}

type PersonSplit struct {
	PersonName  string
	SplitByItem []SplitByItem
	TotalAmount float64
}
