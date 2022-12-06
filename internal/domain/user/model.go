package user

type User struct {
	Id          int    `json:"id" db:"id"`
	RoleCode    string `json:"roleCode" db:"role_code"`
	Email       string `json:"email" db:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Lastname    string `json:"lastname" db:"lastname" binding:"required"`
	Firstname   string `json:"firstname" db:"firstname" binding:"required"`
	MiddleName  string `json:"middleName" db:"middle_name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number" binding:"required"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}

type InfoDTO struct {
	Email       string `json:"email" db:"email"`
	Lastname    string `json:"lastname" db:"lastname"`
	Firstname   string `json:"firstname" db:"firstname"`
	MiddleName  string `json:"middleName" db:"middle_name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
}
