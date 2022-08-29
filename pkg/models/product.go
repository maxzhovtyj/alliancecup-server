package models

import "time"

type Product struct {
	Id             int       `json:"id" db:"id" example:"5"`
	Article        string    `json:"article" binding:"required" db:"article" example:"000123"`
	CategoryTitle  string    `json:"category_title" binding:"required" db:"category_title" example:"Одноразові стакани"`
	ProductTitle   string    `json:"product_title" binding:"required" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl         string    `json:"img_url" db:"img_url" example:"https://google-images.com/some-img123"`
	TypeTitle      string    `json:"type_title" binding:"required" db:"type_title" example:"Стакан"`
	AmountInStock  float64   `json:"amount_in_stock" db:"amount_in_stock" example:"120"`
	Price          float64   `json:"price" binding:"required" db:"price" example:"3.75"`
	UnitsInPackage int       `json:"units_in_package" binding:"required" db:"units_in_package" example:"30"`
	PackagesInBox  int       `json:"packages_in_box" binding:"required" db:"packages_in_box" example:"50"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ProductInfo struct {
	ProductId   int    `json:"product_id" db:"product_id" example:"1"`
	InfoTitle   string `json:"info_title" db:"info_title" example:"Колір"`
	Description string `json:"description" db:"description" example:"Білий"`
}

type ProductInfoDescription struct {
	Info        Product       `json:"info"`
	Description []ProductInfo `json:"description"`
}
