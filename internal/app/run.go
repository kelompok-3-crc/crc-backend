package app

import (
	"fmt"
	"log"
	"ml-prediction/config"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/routes"
	"ml-prediction/pkg/logger"
	"ml-prediction/pkg/validation"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Run() {

	app := fiber.New()
	cfg := config.NewConfig()
	var validate *validator.Validate

	db := config.SetupDatabase(cfg)
	db.AutoMigrate(&model.User{}, &model.Customer{})
	validate = validator.New()
	if err := validation.RegisterCustomValidation(validate); err != nil {
		log.Fatalf("error register custom validation")
	}

	logger, err := logger.Initialize(*cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	api := app.Group("/api/v1")
	routes.Register(api, db, *cfg, logger, validate)

	go func() {
		fmt.Println("Listen and Serve at port 8080")
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()
	log.Print("Server Started")

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt)
	<-stopped

	fmt.Println("shutting down gracefully...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("error in Server Shutdown: %s", err)
	}
	fmt.Println("server stopped")
}
