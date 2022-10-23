package product

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
)

type CreateDTO struct {
	Article         string          `json:"article" db:"article" example:"000123"`
	CategoryTitle   string          `json:"categoryTitle" db:"category_title" example:"Одноразові стакани"`
	ProductTitle    string          `json:"productTitle" db:"product_title" example:"Стакан одноразовий Крафт 110мл"`
	Img             *models.FileDTO `json:"img"`
	ImgUrl          *string         `json:"imgUrl" db:"img_url" example:"https://google-images.com/some-img123"`
	AmountInStock   float64         `json:"amountInStock" db:"amount_in_stock" example:"120"`
	Price           float64         `json:"price" db:"price" example:"3.75"`
	Characteristics *types.JSONText `json:"characteristics" db:"characteristics"`
	Packaging       *types.JSONText `json:"packaging" db:"packaging"`
	Description     *string         `json:"description" db:"description"`
}
