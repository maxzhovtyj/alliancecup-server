package shopping

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type CartProduct struct {
	CartId           int            `json:"-" db:"cart_id"`
	ProductId        int            `json:"productId" db:"product_id" binding:"required"`
	Article          string         `json:"article" example:"000123"`
	ProductTitle     string         `json:"productTitle" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl           *string        `json:"imgUrl" example:"https://google-images.com/some-img123"`
	ImgUUID          *string        `json:"imgUUID"`
	AmountInStock    float64        `json:"amountInStock" example:"120"`
	Price            float64        `json:"price" db:"price" example:"3.75"`
	Packaging        types.JSONText `json:"packaging"`
	CreatedAt        time.Time      `json:"createdAt"`
	Quantity         int            `json:"quantity" db:"quantity" binding:"required"`
	PriceForQuantity float64        `json:"priceForQuantity" db:"price_for_quantity" binding:"required"`
}

type CharacteristicParam struct {
	Name  string
	Value string
}

type SearchParams struct {
	CategoryId     int
	PriceRange     string
	CreatedAt      string
	Characteristic []CharacteristicParam
	Search         string
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
