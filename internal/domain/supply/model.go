package supply

import "time"

type InfoDTO struct {
	Id         int        `json:"id" db:"id"`
	Supplier   string     `json:"supplier" db:"supplier" binding:"required"`
	SupplyTime *time.Time `json:"supplyTime" db:"supply_time"`
	Comment    string     `json:"comment" db:"comment"`
	Sum        float64    `json:"sum" db:"sum"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
}

type PaymentDTO struct {
	PaymentAccount string     `json:"paymentType" db:"payment_account" binding:"required"`
	PaymentTime    *time.Time `json:"paymentTime" db:"payment_time"`
	PaymentSum     float64    `json:"paymentSum" db:"payment_sum" binding:"required"`
}

type ProductDTO struct {
	SupplyId      int     `json:"-" db:"supply_id"`
	ProductId     int     `json:"productId" binding:"required" db:"product_id"`
	Packaging     string  `json:"packaging" db:"packaging"`
	Amount        float64 `json:"amount" binding:"required" db:"amount"`
	PriceForUnit  float64 `json:"priceForUnit" binding:"required" db:"price_for_unit"`
	SumWithoutTax float64 `json:"sumWithoutTax" binding:"required" db:"sum_without_tax"`
	Tax           float64 `json:"tax" binding:"required" db:"tax"`
	TotalSum      float64 `json:"totalSum" binding:"required" db:"total_sum"`
}

type Supply struct {
	Info     InfoDTO      `json:"info"`
	Payment  []PaymentDTO `json:"payment"`
	Products []ProductDTO `json:"products"`
}
