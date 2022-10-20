package models

import (
	"time"
)

type Session struct {
	Id           int       `json:"-" db:"id"`
	UserId       int       `json:"userId" db:"user_id"`
	RoleId       int       `json:"roleId" db:"role_id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token"`
	IsBlocked    bool      `json:"isBlocked" db:"is_blocked"`
	ClientIp     string    `json:"clientIp" db:"client_ip"`
	UserAgent    string    `json:"userAgent" db:"user_agent"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
