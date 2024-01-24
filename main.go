package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/zaza-hikayat/go-fiber/configs"
	"github.com/zaza-hikayat/go-fiber/dto/models"
	"github.com/zaza-hikayat/go-fiber/internal/infrastructure/database"

	infra_logger "github.com/zaza-hikayat/go-fiber/internal/infrastructure/logger"
	"github.com/zaza-hikayat/go-fiber/internal/infrastructure/persistance"
	"github.com/zaza-hikayat/go-fiber/internal/infrastructure/redis"
	"github.com/zaza-hikayat/go-fiber/internal/interface/rest"

	usecases "github.com/zaza-hikayat/go-fiber/internal/use_cases"
)

func main() {
	// init configuration
	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := infra_logger.InitLogger(conf)
	defer logger.Sync()

	// set connection database
	dbConn, err := database.NewDBConnection(*conf)
	if err != nil {
		panic(err)
	}

	// initiliaze repository
	userRepository := persistance.NewUserRepository(dbConn)
	cacheRepository := redis.NewRedisCahce(conf)

	// initiliaze useCase
	allUseCase := usecases.AllUseCases{
		UserUseCase: usecases.NewUserUseCase(userRepository, cacheRepository, conf),
	}

	app := fiber.New(fiber.Config{
		Immutable: true,
		// make encoding json faster than default library
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	rest.NewRouter(app, allUseCase)
	dbConn.AutoMigrate(&models.User{})

	// run http server
	go func() {
		app.Listen(fmt.Sprintf(":%d", conf.Server.Port))
	}()

	// graceful shutdown
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-exitSignal
	shutdownDelay := 5 * time.Second
	logger.Info("Application shutting down")
	db, _ := dbConn.DB()
	db.Close()
	app.ShutdownWithTimeout(shutdownDelay)
	time.Sleep(shutdownDelay)
	logger.Info("Exited")

}
