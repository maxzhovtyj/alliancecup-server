package server

type Category struct {
	Id            int    `json:"id" db:"id"`
	CategoryTitle string `json:"category_title" binding:"required" db:"category_title"`
	ImgUrl        string `json:"img_url" binding:"required" db:"img_url"`
}

type Product struct {
	Id             int     `json:"id" db:"id"`
	Article        string  `json:"article" db:"article"`
	CategoryTitle  string  `json:"category_title" db:"category_title"`
	ProductTitle   string  `json:"product_title" db:"product_title"`
	ImgUrl         string  `json:"img_url" db:"img_url"`
	TypeTitle      string  `json:"type_title" db:"type_title"`
	AmountInStock  float64 `json:"amount_in_stock" db:"amount_in_stock"`
	Price          float64 `json:"price" db:"price"`
	UnitsInPackage int     `json:"units_in_package" db:"units_in_package"`
	PackagesInBox  int     `json:"packages_in_box" db:"packages_in_box"`
}

type ProductInfo struct {
	ProductId   int    `json:"product_id" db:"product_id"`
	InfoTitle   string `json:"info_title" db:"info_title"`
	Description string `json:"description" db:"description"`
}

type ProductInfoDescription struct {
	Info        Product       `json:"info"`
	Description []ProductInfo `json:"description"`
}

type CartProduct struct {
	CartId           int     `json:"-" db:"cart_id"`
	ProductId        int     `json:"product_id" binding:"required" db:"product_id"`
	Quantity         int     `json:"quantity" binding:"required" db:"quantity"`
	PriceForQuantity float64 `json:"price_for_quantity" binding:"required" db:"price_for_quantity"`
}
