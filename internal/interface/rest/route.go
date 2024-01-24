package rest

import (
	"github.com/gofiber/fiber/v2"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/zaza-hikayat/go-fiber/internal/interface/handler"
	"github.com/zaza-hikayat/go-fiber/internal/interface/middleware"
	usecases "github.com/zaza-hikayat/go-fiber/internal/use_cases"
	"github.com/zaza-hikayat/go-fiber/internal/utility"
)

func NewRouter(app *fiber.App, useCases usecases.AllUseCases) {
	respClient := utility.NewResponseClient()

	// initialize handlers
	authHandler := handler.NewUserHandler(useCases.UserUseCase, respClient)

	// middlewarres
	app.Use(requestid.New())
	app.Use(fiber_logger.New(fiber_logger.Config{
		Format: "${locals:requestid} ${status} - ${method} ${path}​ ${latency}\n",
	}))

	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	v1.Post("/login", authHandler.Login)
	v1.Post("/register", authHandler.Register)
	v1.Get("/me", middleware.Authenticate(useCases.UserUseCase), authHandler.GetMe)
}
