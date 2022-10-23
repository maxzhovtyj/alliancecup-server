package category

import (
	"github.com/google/uuid"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
)

type Category struct {
	Id                  int        `json:"id" db:"id"`
	CategoryTitle       string     `json:"categoryTitle" db:"category_title" binding:"required"`
	ImgUrl              *string    `json:"imgUrl" db:"img_url"`
	ImgUUID             *uuid.UUID `json:"imgUUID" db:"img_uuid"`
	CategoryDescription *string    `json:"categoryDescription" db:"category_description"`
}
type CreateDTO struct {
	CategoryTitle       string
	ImgUrl              *string
	Img                 models.FileDTO
	CategoryDescription *string
}

type Filtration struct {
	Id                    int     `json:"id" db:"id"`
	CategoryId            *int    `json:"categoryId" db:"category_id"`
	ImgUrl                *string `json:"imgUrl" db:"img_url"`
	InfoDescription       string  `json:"infoDescription" binding:"required" db:"info_description"`
	FiltrationTitle       string  `json:"filtrationTitle" binding:"required" db:"filtration_title"`
	FiltrationDescription *string `json:"filtrationDescription" db:"filtration_description"`
	FiltrationListId      *int    `json:"filtrationListId" db:"filtration_list_id"`
}
