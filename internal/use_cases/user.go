package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/zaza-hikayat/go-fiber/configs"
	"github.com/zaza-hikayat/go-fiber/domain"
	constant "github.com/zaza-hikayat/go-fiber/internal/constants"
	infra_error "github.com/zaza-hikayat/go-fiber/internal/infrastructure/errors"
	"github.com/zaza-hikayat/go-fiber/internal/utility"

	"github.com/zaza-hikayat/go-fiber/dto/models"
	"github.com/zaza-hikayat/go-fiber/dto/request"
	"github.com/zaza-hikayat/go-fiber/dto/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUseCase struct {
	userRepo  domain.UserRepository
	cacheRepo domain.RedisCache
	config    *configs.Config
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	cacheRepo domain.RedisCache,
	config *configs.Config,
) domain.UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		config:    config,
	}
}

func (u *userUseCase) Authenticate(ctx context.Context, req request.SignReq) (models.User, string, error) {
	var tokenStr string
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, tokenStr, infra_error.NewCommonError(infra_error.INVALID_DATA, err)
		}
		return user, tokenStr, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {

		return user, tokenStr, err
	}
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["id"] = user.ID
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := []byte(u.config.Application.SecretKey)
	tokenStr, _ = token.SignedString(signingKey)
	userJson, _ := json.Marshal(user)
	u.cacheRepo.Set(ctx, "user:"+tokenStr, string(userJson), time.Duration(u.config.Redis.DefaultExpire)*time.Minute)
	return user, tokenStr, nil
}

func (u *userUseCase) ValidateUserToken(ctx context.Context, token string) (user models.User, err error) {
	cacheKey := "user:" + token
	result, _ := u.cacheRepo.Get(ctx, cacheKey).Result()

	if result != "" {
		_ = json.Unmarshal([]byte(result), &user)
	} else {
		signingKey := []byte(u.config.Application.SecretKey)
		claims := make(jwt.MapClaims)

		parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}
			return signingKey, nil
		})

		if err != nil || parsedToken == nil {
			cerr := infra_error.NewCommonError(infra_error.INVALID_TOKEN, err)
			return user, cerr
		}

		if !parsedToken.Valid {
			cerr := infra_error.NewCommonError(infra_error.INVALID_TOKEN, err)
			return user, cerr
		}

		if _, ok := claims["id"]; ok {
			idNum, _ := strconv.Atoi(claims["id"].(string))
			if idNum != 0 {
				user, err = u.userRepo.FindByUserId(ctx, claims["id"].(string))
			}
		}

		return user, err
	}

	return user, nil
}

func (u *userUseCase) RegisterUser(ctx context.Context, req request.RegisterReq) (resp response.UserLoginResp, err error) {
	var user models.User
	// check if user is already registered
	user, _ = u.userRepo.FindByEmail(ctx, req.Email)
	if user != (models.User{}) {
		cerr := infra_error.NewCommonError(infra_error.INVALID_DATA, err)
		cerr.SetSystemMessage(fmt.Sprintf("user %v is already registered", req.Email))
		return resp, infra_error.NewCommonError(infra_error.INVALID_DATA, cerr)
	}
	paswordHashedByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return resp, err
	}
	now := time.Now()
	user = models.User{
		Fullname:     req.Fullname,
		Email:        req.Email,
		PasswordHash: string(paswordHashedByte),
		PhoneNumber:  req.PhoneNumber,
		VerifiedAt:   &now,
	}

	user, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return resp, err
	}
	mapstructure.Decode(user, &resp)
	return resp, nil
}

func (u *userUseCase) SendOtp(ctx context.Context, req request.SendOtpReq) (string, error) {
	otpUser := utility.GenerateRandomNumber()
	expiresTime := 5 * time.Minute
	r := u.cacheRepo.Set(ctx, fmt.Sprintf("otp_user_%s", req.Email), otpUser, expiresTime)
	if r.Err() != nil {
		return constant.EMPTY_STR, r.Err()
	}

	// TODO send notification email or sms
	return otpUser, nil
}
func (u *userUseCase) VerifyOtp(ctx context.Context, req request.VerifyOtpReq) (err error) {
	r, err := u.cacheRepo.Get(ctx, fmt.Sprintf("otp_user_%s", req.Email)).Result()
	if err != nil || r == constant.EMPTY_STR {
		cerr := infra_error.NewCommonError(infra_error.INVALID_PAYLOAD, err)
		cerr.SetSystemMessage("Invalid OTP token")
		return cerr
	}

	return nil
}
