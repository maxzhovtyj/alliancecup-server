package server

type Category struct {
	Id            int    `json:"-" db:"id"`
	CategoryTitle string `json:"category_title" binding:"required" db:"category_title"`
	ImgUrl        string `json:"img_url" binding:"required" db:"img_url"`
}

type Product struct {
	Id             int     `json:"-" db:"id"`
	Article        string  `json:"article" db:"article"`
	CategoryTitle  string  `json:"category_title" db:"category_title"`
	ProductTitle   string  `json:"product_title" db:"product_title"`
	ImgUrl         string  `json:"img_url" db:"img_url"`
	TypeTitle      string  `json:"type_title" db:"type_title"`
	AmountInStock  int     `json:"amount_in_stock" db:"amount_in_stock"`
	Price          float64 `json:"price" db:"price"`
	UnitsInPackage int     `json:"units_in_package" db:"units_in_package"`
	PackagesInBox  int     `json:"packages_in_box" db:"packages_in_box"`
}

type ProductInfo struct {
	InfoTitle   string `json:"info_title"`
	Description string `json:"description"`
}
