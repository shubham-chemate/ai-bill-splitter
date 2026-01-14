package main

type BillItem struct {
	ItemName     string   `json:"itemName"`
	PricePerUnit float64  `json:"pricePerUnit"`
	Quantity     int      `json:"quantity"`
	Tax          float64  `json:"tax"`
	TotalPrice   float64  `json:"totalPrice"`
	Warnings     []string `json:"warnings"`
}

type PersonShare struct {
	Name  string  `json:"name"`
	Share float64 `json:"share"`
}

type ItemSplit struct {
	Name         string        `json:"itemName"`
	PersonSplits []PersonShare `json:"splits"`
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
