package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/zaza-hikayat/go-fiber/domain"
	constant "github.com/zaza-hikayat/go-fiber/internal/constants"
	infra_error "github.com/zaza-hikayat/go-fiber/internal/infrastructure/errors"
)

func Authenticate(useCase domain.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.ReplaceAll(c.Get("Authorization"), "Bearer", "")
		token = strings.TrimSpace(token)
		if token == "" {
			cerr := infra_error.NewCommonError(infra_error.INVALID_USER, nil)
			return c.Status(http.StatusForbidden).JSON(cerr.ToHttpError())
		}

		user, err := useCase.ValidateUserToken(context.TODO(), token)
		if err != nil {
			cerr := infra_error.NewCommonError(infra_error.INVALID_TOKEN, nil)
			return c.Status(http.StatusForbidden).JSON(cerr.ToHttpError())
		}
		userJson, _ := json.Marshal(user)
		// Set a custom header on all responses:
		c.Locals(constant.X_USER, string(userJson))

		// Go to next middleware:
		return c.Next()
	}
}
