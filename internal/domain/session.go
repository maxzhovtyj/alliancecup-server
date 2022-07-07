package domain

import "time"

type Session struct {
	Id           int       `json:"-" db:"id"`
	UserId       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IsBlocked    bool      `json:"is_blocked" db:"is_blocked"`
	ClientIp     string    `json:"client_ip" db:"client_ip"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
