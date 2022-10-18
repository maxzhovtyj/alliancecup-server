package user

type User struct {
	Id          int    `json:"-" db:"id"`
	RoleId      int    `json:"role_id" db:"role_id"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type InfoDTO struct {
	Email       string `json:"email" db:"email"`
	Name        string `json:"name" db:"name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
}
