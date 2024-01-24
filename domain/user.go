package domain

import (
	"context"

	"github.com/zaza-hikayat/go-fiber/dto/models"
	"github.com/zaza-hikayat/go-fiber/dto/request"
	"github.com/zaza-hikayat/go-fiber/dto/response"
)

type UserRepository interface {
	FindByUserId(ctx context.Context, userId string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type UserUseCase interface {
	Authenticate(ctx context.Context, req request.SignReq) (user models.User, token string, err error)
	ValidateUserToken(ctx context.Context, token string) (models.User, error)
	RegisterUser(ctx context.Context, req request.RegisterReq) (response.UserLoginResp, error)
	SendOtp(ctx context.Context, req request.SendOtpReq) (token string, err error)
	VerifyOtp(ctx context.Context, req request.VerifyOtpReq) (err error)
}
