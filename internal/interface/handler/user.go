package handler

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/zaza-hikayat/go-fiber/domain"
	"github.com/zaza-hikayat/go-fiber/dto/models"
	"github.com/zaza-hikayat/go-fiber/dto/request"
	constant "github.com/zaza-hikayat/go-fiber/internal/constants"
	infra_error "github.com/zaza-hikayat/go-fiber/internal/infrastructure/errors"
	"github.com/zaza-hikayat/go-fiber/internal/utility"
)

type UserHandler interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
	SendOtp(c *fiber.Ctx) error
	VerifyOtp(c *fiber.Ctx) error
}

type userHandler struct {
	userUseCase domain.UserUseCase
	respClient  utility.IResponseClient
}

func NewUserHandler(
	userUseCase domain.UserUseCase,
	respClient utility.IResponseClient,
) UserHandler {
	return &userHandler{
		userUseCase: userUseCase,
		respClient:  respClient,
	}
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	req := new(request.SignReq)
	// parse request
	if err := c.BodyParser(req); err != nil {
		return h.respClient.HttpError(c, err)
	}

	if err := req.Validate(); err != nil {
		cerr := infra_error.NewCommonError(infra_error.INVALID_PAYLOAD, err)
		return h.respClient.HttpError(c, cerr)
	}

	user, tokenStr, err := h.userUseCase.Authenticate(c.Context(), *req)
	if err != nil {
		return h.respClient.HttpError(c, err)
	}

	h.respClient.JSON(c, fiber.Map{"user": user, "token": tokenStr}, nil)
	return nil
}

// Register implements UserHandler.
func (h *userHandler) Register(c *fiber.Ctx) error {
	req := new(request.RegisterReq)
	if err := c.BodyParser(req); err != nil {
		return h.respClient.HttpError(c, err)
	}

	if err := req.Validate(); err != nil {
		cerr := infra_error.NewCommonError(infra_error.INVALID_PAYLOAD, err)
		return h.respClient.HttpError(c, cerr)
	}
	user, err := h.userUseCase.RegisterUser(c.Context(), *req)
	if err != nil {
		return h.respClient.HttpError(c, err)
	}

	return h.respClient.JSON(c, user, nil)
}

func (h *userHandler) GetMe(c *fiber.Ctx) error {
	var userModel models.User
	userJson := c.Locals(constant.X_USER)
	_ = json.Unmarshal([]byte(userJson.(string)), &userModel)
	return h.respClient.JSON(c, userModel, nil)
}

func (h *userHandler) SendOtp(c *fiber.Ctx) error {
	req := new(request.SendOtpReq)
	if err := c.BodyParser(req); err != nil {
		cerr := infra_error.NewCommonError(infra_error.INVALID_PAYLOAD, err)
		cerr.SetSystemMessage(err.Error())
		return h.respClient.HttpError(c, cerr)
	}
	otp, err := h.userUseCase.SendOtp(c.Context(), *req)
	if err != nil {
		return h.respClient.HttpError(c, err)
	}
	return h.respClient.JSON(c, fiber.Map{"otp": otp}, nil)
}

func (h *userHandler) VerifyOtp(c *fiber.Ctx) error {
	panic("unimplemented")
}
