package product

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type Product struct {
	Id              int             `json:"id" db:"id" example:"5"`
	Article         string          `json:"article" binding:"required" db:"article" example:"000123"`
	CategoryTitle   string          `json:"categoryTitle" binding:"required" db:"category_title" example:"Одноразові стакани"`
	ProductTitle    string          `json:"productTitle" binding:"required" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl          *string         `json:"img_url" db:"img_url" example:"https://google-images.com/some-img123"`
	AmountInStock   float64         `json:"amountInStock" db:"amount_in_stock" example:"120"`
	Price           float64         `json:"price" binding:"required" db:"price" example:"3.75"`
	Characteristics *types.JSONText `json:"characteristics" db:"characteristics"`
	Packaging       *types.JSONText `json:"packaging" db:"packaging"`
	Description     string          `json:"description" db:"description"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
}
