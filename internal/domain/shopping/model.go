package shopping

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type CartProduct struct {
	CartId           int             `json:"-" db:"cart_id"`
	Id               int             `json:"id" db:"product_id" binding:"required"`
	Article          string          `json:"article" db:"article" example:"000123"`
	ProductTitle     string          `json:"productTitle" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl           *string         `json:"imgUrl" db:"img_url" example:"https://google-images.com/some-img123"`
	ImgUUID          *string         `json:"imgUUID" db:"img_uuid"`
	AmountInStock    float64         `json:"amountInStock" db:"amount_in_stock" example:"120"`
	Price            float64         `json:"price" db:"price" example:"3.75"`
	Packaging        *types.JSONText `json:"packaging" db:"packaging"`
	Characteristics  *types.JSONText `json:"characteristics" db:"characteristics"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	Quantity         int             `json:"quantity" db:"quantity" binding:"required"`
	PriceForQuantity float64         `json:"priceForQuantity" db:"price_for_quantity" binding:"required"`
}

type CharacteristicParam struct {
	Name  string
	Value string
}

type SearchParams struct {
	CategoryId     int
	PriceRange     string
	CreatedAt      string
	IsActive       *bool
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
