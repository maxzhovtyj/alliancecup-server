package server

type ShopItemCup struct {
	Id       int     `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Size     int     `json:"size"`
	Category string  `json:"category"`
}

type ShopList struct {
}

type Category struct {
}
