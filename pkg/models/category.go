package models

type Category struct {
	Id                  int     `json:"id" db:"id"`
	CategoryTitle       string  `json:"category_title" binding:"required" db:"category_title"`
	ImgUrl              *string `json:"img_url" db:"img_url"`
	CategoryDescription *string `json:"category_description" db:"category_description"`
}

type CategoryFiltration struct {
	Id                    int     `json:"id" db:"id"`
	CategoryId            *int    `json:"categoryId" db:"category_id"`
	ImgUrl                *string `json:"imgUrl" db:"img_url"`
	InfoDescription       string  `json:"infoDescription" binding:"required" db:"info_description"`
	FiltrationTitle       string  `json:"filtrationTitle" binding:"required" db:"filtration_title"`
	FiltrationDescription *string `json:"filtrationDescription" db:"filtration_description"`
	FiltrationListId      *int    `json:"filtrationListId" db:"filtration_list_id"`
}
