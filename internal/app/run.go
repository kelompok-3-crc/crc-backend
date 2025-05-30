package app

import (
	"context"
	"fmt"
	"log"
	"ml-prediction/config"
	"ml-prediction/internal/app/routes"
	"ml-prediction/pkg/logger"
	"ml-prediction/pkg/utils"
	"ml-prediction/pkg/validation"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run() {

	cors := cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, X-Total-Count",
		MaxAge:           3600,
	})
	app := fiber.New()
	app.Use(cors)
	cfg := config.NewConfig()
	var validate *validator.Validate

	db := config.SetupDatabase(cfg)
	validate = validator.New()
	if err := validation.RegisterCustomValidation(validate, db); err != nil {
		log.Fatalf("error register custom validation")
	}

	logger, err := logger.Initialize(*cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	err = utils.ImportInitialCustomerData(context.Background(), db)
	if err != nil {
		log.Fatalf("failed to run import data: %v", err)
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
