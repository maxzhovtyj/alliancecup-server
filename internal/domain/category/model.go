package category

import (
	"github.com/google/uuid"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
)

type Category struct {
	Id            int        `json:"id" db:"id"`
	CategoryTitle string     `json:"categoryTitle" db:"category_title" binding:"required"`
	ImgUrl        *string    `json:"imgUrl" db:"img_url"`
	ImgUUID       *uuid.UUID `json:"imgUUID" db:"img_uuid"`
	Description   *string    `json:"description" db:"description"`
}

type CreateDTO struct {
	CategoryTitle       string
	ImgUrl              *string
	Img                 *models.FileDTO
	CategoryDescription *string
}

type UpdateImageDTO struct {
	Id  int
	Img *models.FileDTO
}

type CreateFiltrationDTO struct {
	Id                    int
	CategoryId            *int
	ImgUrl                *string
	Img                   *models.FileDTO
	SearchKey             string
	SearchCharacteristic  string
	FiltrationTitle       string
	FiltrationDescription *string
	FiltrationListId      *int
}

type Filtration struct {
	Id                    int        `json:"id" db:"id"`
	CategoryId            *int       `json:"categoryId" db:"category_id"`
	ImgUrl                *string    `json:"imgUrl" db:"img_url"`
	ImgUUID               *uuid.UUID `json:"imgUUID" db:"img_uuid"`
	SearchKey             string     `json:"searchKey" db:"search_key" binding:"required"`
	SearchCharacteristic  string     `json:"searchCharacteristic" db:"search_characteristic" binding:"required"`
	FiltrationTitle       string     `json:"filtrationTitle" db:"filtration_title" binding:"required"`
	FiltrationDescription *string    `json:"filtrationDescription" db:"filtration_description"`
	FiltrationListId      *int       `json:"filtrationListId" db:"filtration_list_id"`
}
