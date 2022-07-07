package server

import "time"

type Category struct {
	Id    int    `json:"-" db:"id"`
	Title string `json:"title"`
}

type Product struct {
	Id             int     `json:"-" db:"id"`
	CategoryId     int     `json:"category_id"`
	Title          string  `json:"title"`
	Price          float64 `json:"price"`
	Size           int     `json:"size"`
	Characteristic string  `json:"characteristic"`
	Description    string  `json:"description"`
	Amount         int     `json:"amount"`
	InStock        bool    `json:"in_stock"`
}

type Order struct {
	Id        int
	ProductId int
	UserId    int
	OrderDate time.Time
	Amount    int
	EndStatus bool
}
