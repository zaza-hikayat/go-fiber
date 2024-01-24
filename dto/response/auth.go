package response

import "time"

type UserLoginResp struct {
	ID           int64      `json:"id"`
	CreatedAt    *time.Time `json:"createdAt"`
	Fullname     string     `json:"fullname"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	PhoneNumber  string     `json:"phone_number"`
	VerifiedAt   *time.Time `json:"verified_at"`
}
