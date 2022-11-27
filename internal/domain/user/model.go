package user

type User struct {
	Id          int    `json:"-" db:"id"`
	RoleCode    string `json:"roleCode" db:"role_code"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Firstname   string `json:"firstname" binding:"required"`
	MiddleName  string `json:"middleName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type InfoDTO struct {
	Email       string `json:"email" db:"email"`
	Lastname    string `json:"lastname" db:"lastname"`
	Firstname   string `json:"firstname" db:"firstname"`
	MiddleName  string `json:"middleName" db:"middle_name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
}
