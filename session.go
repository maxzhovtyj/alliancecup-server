package server

import (
	"time"
)

type Session struct {
	Id           int       `json:"-" db:"id"`
	UserId       int       `json:"user_id"`
	RoleId       int       `json:"role_id"`
	RefreshToken string    `json:"refresh_token"`
	IsBlocked    bool      `json:"is_blocked"`
	ClientIp     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
