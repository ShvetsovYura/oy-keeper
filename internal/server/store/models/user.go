package models

import "time"

type UserDB struct {
	Uuid         string
	Login        string
	Otp_secret   string
	Otp_auth     string
	Otp_verified bool
	Is_active    bool
	Created_at   time.Time
	Updated_at   time.Time
}
