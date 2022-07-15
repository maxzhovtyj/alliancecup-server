package server

type Category struct {
	Id            int    `json:"-" db:"id"`
	CategoryTitle string `json:"category_title"`
}

type Product struct {
	Id             int     `json:"-" db:"id"`
	Article        string  `json:"article"`
	CategoryId     int     `json:"category_id"`
	ProductTitle   string  `json:"product_title"`
	ImgUrl         string  `json:"img_url"`
	TypeId         int     `json:"-" db:"type_id"`
	AmountInStock  int     `json:"amount_in_stock"`
	Price          float64 `json:"price"`
	UnitsInPackage int     `json:"units_in_package"`
	PackagesInBox  int     `json:"packages_in_box"`
}

type ProductInfo struct {
	ProductId   int    `json:"product_id"`
	InfoTitle   string `json:"info_title"`
	Description string `json:"description"`
}
