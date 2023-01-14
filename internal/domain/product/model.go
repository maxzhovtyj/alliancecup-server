package product

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"time"
)

type Product struct {
	Id              int             `json:"id" db:"id" example:"5"`
	Article         string          `json:"article" db:"article" binding:"required" example:"000123"`
	CategoryTitle   *string         `json:"categoryTitle" db:"category_title" binding:"required" example:"Одноразові стакани"`
	ProductTitle    string          `json:"productTitle" db:"product_title" binding:"required" example:"Стакан одноразовий Крафт 110мл"`
	ImgUrl          *string         `json:"imgUrl" db:"img_url" example:"https://google-images.com/some-img123"`
	ImgUUID         *uuid.UUID      `json:"imgUUID" db:"img_uuid" example:"https://google-images.com/some-img123"`
	AmountInStock   float64         `json:"amountInStock" db:"amount_in_stock" example:"120"`
	Price           float64         `json:"price" db:"price" binding:"required" example:"3.75"`
	Characteristics *types.JSONText `json:"characteristics" db:"characteristics"`
	Packaging       *types.JSONText `json:"packaging" db:"packaging"`
	Description     *string         `json:"description" db:"description"`
	IsActive        bool            `json:"isActive" db:"is_active"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
}
