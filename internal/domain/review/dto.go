package review

type CreateReviewDTO struct {
	UserId     *int   `json:"userId,omitempty"`
	UserName   string `json:"userName" binding:"required"`
	Mark       int    `json:"mark" binding:"required"`
	ReviewText string `json:"reviewText" binding:"required"`
}
