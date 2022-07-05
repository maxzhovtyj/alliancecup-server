package server

type ShopList struct {
}

type Category struct {
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
