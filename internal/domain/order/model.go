package order

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type Order struct {
	Id                int        `json:"id" db:"id"`
	UserId            *int       `json:"user_id" db:"user_id"`
	UserLastName      string     `json:"user_lastname" db:"user_lastname" binding:"required"`
	UserFirstName     string     `json:"user_firstname" db:"user_firstname" binding:"required"`
	UserMiddleName    string     `json:"user_middle_name" db:"user_middle_name" binding:"required"`
	UserPhoneNumber   string     `json:"user_phone_number" db:"user_phone_number" binding:"required"`
	UserEmail         string     `json:"user_email" db:"user_email" binding:"required"`
	Status            string     `json:"status" db:"status"`
	Comment           string     `json:"comment" db:"comment"`
	SumPrice          float64    `json:"sum_price" db:"sum_price" binding:"required"`
	DeliveryTypeTitle string     `json:"delivery_type_title" db:"delivery_type_title" binding:"required"`
	PaymentTypeTitle  string     `json:"payment_type_title" db:"payment_type_title" binding:"required"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ClosedAt          *time.Time `json:"closed_at" db:"closed_at"`
}

type Product struct {
	OrderId          int     `json:"-" db:"order_id"`
	ProductId        int     `json:"productId" db:"product_id"`
	Quantity         int     `json:"quantity" db:"quantity"`
	PriceForQuantity float64 `json:"priceForQuantity" db:"price_for_quantity"`
}

type OrdersDelivery struct {
	OrderId             int    `json:"-" db:"order_id"`
	DeliveryTitle       string `json:"deliveryTitle" db:"delivery_title"`
	DeliveryDescription string `json:"deliveryDescription" db:"delivery_description"`
}

type CreateDTO struct {
	Order    Order            `json:"order"`
	Products []Product        `json:"products"`
	Delivery []OrdersDelivery `json:"delivery"`
}

type ProductFullInfo struct {
	Id               int            `json:"id" db:"id"`
	OrderId          int            `json:"orderId" db:"order_id"`
	Article          string         `json:"article" db:"article"`
	ProductTitle     string         `json:"productTitle" db:"product_title"`
	ImgUrl           *string        `json:"imgUrl" db:"img_url"`
	AmountInStock    float64        `json:"amountInStock" db:"amount_in_stock"`
	Price            float64        `json:"price" db:"price"`
	Packaging        types.JSONText `json:"packaging" db:"packaging"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	Quantity         int            `json:"quantity" db:"quantity"`
	PriceForQuantity float64        `json:"priceForQuantity" db:"price_for_quantity"`
}

type SelectDTO struct {
	Info     Order             `json:"info"`
	Products []ProductFullInfo `json:"products"`
	Delivery []OrdersDelivery  `json:"delivery"`
}
