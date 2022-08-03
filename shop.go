package server

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"time"
)

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Category struct {
	Id            int    `json:"id" db:"id"`
	CategoryTitle string `json:"category_title" binding:"required" db:"category_title"`
	ImgUrl        string `json:"img_url" binding:"required" db:"img_url"`
}

type Product struct {
	Id             int       `json:"id" db:"id" example:"5"`
	Article        string    `json:"article" db:"article" example:"000123"`
	CategoryTitle  string    `json:"category_title" db:"category_title" example:"Одноразові стакани"`
	ProductTitle   string    `json:"product_title" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl         string    `json:"img_url" db:"img_url" example:"https://google-images.com/some-img123"`
	TypeTitle      string    `json:"type_title" db:"type_title" example:"Стакан"`
	AmountInStock  float64   `json:"amount_in_stock" db:"amount_in_stock" example:"120"`
	Price          float64   `json:"price" db:"price" example:"3.75"`
	UnitsInPackage int       `json:"units_in_package" db:"units_in_package" example:"30"`
	PackagesInBox  int       `json:"packages_in_box" db:"packages_in_box" example:"50"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
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

type SearchParams struct {
	CategoryTitle  string `json:"category_title"`
	Size           int    `json:"size"`
	Price          string `json:"price"`
	Characteristic string `json:"characteristic"`
}

type Order struct {
	Id                uuid.UUID    `json:"-" db:"id"`
	UserId            int          `json:"-" db:"user_id"`
	UserLastName      string       `json:"user_lastname" db:"user_lastname"`
	UserFirstName     string       `json:"user_firstname" db:"user_firstname"`
	UserMiddleName    string       `json:"user_middle_name" db:"user_middle_name"`
	UserPhoneNumber   string       `json:"user_phone_number" db:"user_phone_number"`
	UserEmail         string       `json:"user_email" db:"user_email"`
	OrderStatus       string       `json:"order_status" db:"order_status"`
	OrderComment      string       `json:"order_comment" db:"order_comment"`
	OrderSumPrice     float64      `json:"order_sum_price" db:"order_sum_price"`
	DeliveryTypeTitle string       `json:"delivery_type_title" db:"delivery_type_title"`
	PaymentTypeTitle  string       `json:"payment_type_title" db:"payment_type_title"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	ClosedAt          sql.NullTime `json:"closed_at" db:"closed_at"`
}

type OrderProducts struct {
	OrderId          uuid.UUID `json:"-"`
	ProductId        int       `json:"product_id"`
	Quantity         int       `json:"quantity"`
	PriceForQuantity float64   `json:"price_for_quantity"`
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
