package models

import (
	_ "github.com/google/uuid"
	"time"
)

type CartProduct struct {
	CartId           int     `json:"-" db:"cart_id"`
	ProductId        int     `json:"product_id" binding:"required" db:"product_id"`
	Quantity         int     `json:"quantity" binding:"required" db:"quantity"`
	PriceForQuantity float64 `json:"price_for_quantity" binding:"required" db:"price_for_quantity"`
}

type CartProductFullInfo struct {
	CartId           int       `json:"-" db:"cart_id"`
	ProductId        int       `json:"product_id" binding:"required" db:"product_id"`
	Article          string    `json:"article" db:"article" example:"000123"`
	ProductTitle     string    `json:"product_title" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl           string    `json:"img_url" db:"img_url" example:"https://google-images.com/some-img123"`
	AmountInStock    float64   `json:"amount_in_stock" db:"amount_in_stock" example:"120"`
	Price            float64   `json:"price" db:"price" example:"3.75"`
	UnitsInPackage   int       `json:"units_in_package" db:"units_in_package" example:"30"`
	PackagesInBox    int       `json:"packages_in_box" db:"packages_in_box" example:"50"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Quantity         int       `json:"quantity" binding:"required" db:"quantity"`
	PriceForQuantity float64   `json:"price_for_quantity" binding:"required" db:"price_for_quantity"`
}

type SearchParams struct {
	CategoryId     int    `json:"category_id"`
	Size           string `json:"size"`
	Price          string `json:"price"`
	Characteristic string `json:"characteristic"`
}

type DeliveryType struct {
	Id                string `json:"id" db:"id"`
	DeliveryTypeTitle string `json:"delivery_type_title" db:"delivery_type_title"`
}

type PaymentType struct {
	Id               string `json:"id" db:"id"`
	PaymentTypeTitle string `json:"payment_type_title" db:"payment_type_title"`
}

type DeliveryPaymentTypes struct {
	DeliveryTypes []DeliveryType `json:"deliveryTypes"`
	PaymentTypes  []PaymentType  `json:"paymentTypes"`
}
