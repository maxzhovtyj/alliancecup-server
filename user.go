package server

type User struct {
	Id       int    `json:"-" db:"id"`
	RoleId   int    `json:"roleId" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
