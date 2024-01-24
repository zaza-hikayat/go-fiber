package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SignReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 20)),
	)
}

type RegisterReq struct {
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

func (r RegisterReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Fullname, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 20)),
	)
}
