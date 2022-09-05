package review

import "time"

type CreateReviewDTO struct {
	UserId     *int   `json:"userId,omitempty"`
	UserName   string `json:"userName" binding:"required"`
	Mark       int    `json:"mark" binding:"required"`
	ReviewText string `json:"reviewText" binding:"required"`
}

type SelectReviewsDTO struct {
	Id         int       `json:"id" db:"id"`
	UserId     *int      `json:"userId" db:"user_id"`
	UserName   string    `json:"userName" db:"user_name"`
	Mark       int       `json:"mark" db:"mark"`
	ReviewText string    `json:"reviewText" db:"review_text"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}
