package routes

import (
	"ml-prediction/config"
	"ml-prediction/internal/app/handler"
	"ml-prediction/internal/app/repository"
	"ml-prediction/internal/app/usecase"
	"ml-prediction/internal/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Register(api fiber.Router, db *gorm.DB, cfg config.Configuration, log *zap.Logger, val *validator.Validate) {

	kantorCabangRepo := repository.NewKantorCabangRepository(db, log)
	kantorCabangService := usecase.NewKantorCabangUsecase(kantorCabangRepo)
	kcHandler := handler.NewKantorCabangHandler(kantorCabangService, cfg, val)

	userRepo := repository.NewUserRepo(db, log)
	authService := usecase.NewAuthUsecase(userRepo, kantorCabangRepo)
	authHandler := handler.NewAuthHandler(authService, cfg, val)

	customerRepo := repository.NewCustomerRepo(db, log)
	productRepo := repository.NewProductRepo(db, log)
	customerService := usecase.NewcustomerUsecase(customerRepo, userRepo, productRepo, db)
	customerHandler := handler.NewCustomerHandler(customerService, cfg, val)

	targetRepo := repository.NewTargetRepository(db, log)
	TargetUsecase := usecase.NewTargetUsecase(userRepo, targetRepo, db)
	targetHandler := handler.NewTargetHandler(TargetUsecase, cfg, val)

	marketingTargetUsecase := usecase.NewMarketingTargetUsecase(targetRepo, userRepo, db)
	marketingTargetHandler := handler.NewMarketingTargetHandler(marketingTargetUsecase, cfg, val)

	marketingCustomerRepo := repository.NewMarketingCustomerRepository(db, log)
	marketingCustomerUsecase := usecase.NewMarketingCustomerUsecase(marketingCustomerRepo, userRepo, db)
	marketingCustomerHandler := handler.NewMarketingCustomerHandler(marketingCustomerUsecase, cfg, val)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/create", middleware.JWTMiddleware("admin"), authHandler.CreateUser)

	predict := api.Group("/predictions")
	predict.Post("/", customerHandler.CreateCustomer)

	targetsRoute := api.Group("/profile", middleware.JWTMiddleware("marketing", "bm"))
	targetsRoute.Get("/summary", targetHandler.GetTargetSummary)

	marketing := api.Group("/marketing", middleware.JWTMiddleware("marketing"))
	marketing.Get("/customers", customerHandler.GetNewCustomers)
	marketing.Get("/customers/me", customerHandler.GetAssignedCustomers)
	marketing.Post("/customer/:cif", marketingCustomerHandler.UpdateCustomerStatus)
	marketing.Get("/customers/:cif", customerHandler.GetCustomerDetail)

	marketing.Get("/monitoring/target", marketingCustomerHandler.GetMonthlyMonitoringMarketing)

	kc := api.Group("/kantor-cabang", middleware.JWTMiddleware("admin"))
	kc.Post("/", kcHandler.Create)
	kc.Get("/", kcHandler.GetAll)

	bm := api.Group("/bm", middleware.JWTMiddleware("bm"))
	bm.Get("/monitoring/target", marketingCustomerHandler.GetMonthlyMonitoring)
	bm.Get("/monitoring/assignment", marketingTargetHandler.GetMarketingTargets)
	bm.Get("/monitoring/assignment/:nip", marketingTargetHandler.GetMarketingTargetsDetail)
	bm.Post("/monitoring/assignment/:nip", marketingTargetHandler.AssignMarketingTarget)
	bm.Get("/monitoring/product-performance", marketingCustomerHandler.GetProductPerformance)

}
