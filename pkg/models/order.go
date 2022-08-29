package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Id                uuid.UUID  `json:"id" db:"id"`
	UserId            int        `json:"-" db:"user_id"`
	UserLastName      string     `json:"user_lastname" binding:"required" db:"user_lastname"`
	UserFirstName     string     `json:"user_firstname" binding:"required" db:"user_firstname"`
	UserMiddleName    string     `json:"user_middle_name" binding:"required" db:"user_middle_name"`
	UserPhoneNumber   string     `json:"user_phone_number" binding:"required" db:"user_phone_number"`
	UserEmail         string     `json:"user_email" binding:"required" db:"user_email"`
	OrderStatus       string     `json:"order_status" db:"order_status"`
	OrderComment      string     `json:"order_comment" db:"order_comment"`
	OrderSumPrice     float64    `json:"order_sum_price" binding:"required" db:"order_sum_price"`
	DeliveryTypeTitle string     `json:"delivery_type_title" binding:"required" db:"delivery_type_title"`
	PaymentTypeTitle  string     `json:"payment_type_title" binding:"required" db:"payment_type_title"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ClosedAt          *time.Time `json:"closed_at" db:"closed_at"`
}

type OrderProducts struct {
	OrderId          uuid.UUID `json:"-" db:"order_id"`
	ProductId        int       `json:"product_id" db:"product_id"`
	Quantity         int       `json:"quantity" db:"quantity"`
	PriceForQuantity float64   `json:"price_for_quantity" db:"price_for_quantity"`
}

type OrdersDelivery struct {
	OrderId             uuid.UUID `json:"-" db:"order_id"`
	DeliveryTitle       string    `json:"delivery_title" db:"delivery_title"`
	DeliveryDescription string    `json:"delivery_description" db:"delivery_description"`
}

type OrderFullInfo struct {
	Info     Order            `json:"info"`
	Products []OrderProducts  `json:"products"`
	Delivery []OrdersDelivery `json:"delivery"`
}

type OrderProductFullInfo struct {
	Id               int       `json:"id" db:"id"`
	OrderId          uuid.UUID `json:"order_id" db:"order_id"`
	Article          string    `json:"article" db:"article"`
	ProductTitle     string    `json:"product_title" db:"product_title"`
	ImgUrl           string    `json:"img_url" db:"img_url"`
	AmountInStock    float64   `json:"amount_in_stock" db:"amount_in_stock"`
	Price            float64   `json:"price" db:"price"`
	UnitsInPackage   int       `json:"units_in_package" db:"units_in_package"`
	PackagesInBox    int       `json:"packages_in_box" db:"packages_in_box"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Quantity         int       `json:"quantity" db:"quantity"`
	PriceForQuantity float64   `json:"price_for_quantity" db:"price_for_quantity"`
}

type OrderInfo struct {
	Info     Order                  `json:"info"`
	Products []OrderProductFullInfo `json:"products"`
	Delivery []OrdersDelivery       `json:"delivery"`
}
