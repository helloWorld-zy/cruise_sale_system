package main

import (
	"log"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/config"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/cruisebooking/backend/internal/pkg/database"
	"github.com/cruisebooking/backend/internal/pkg/logger"
	"github.com/cruisebooking/backend/internal/repository"
	"github.com/cruisebooking/backend/internal/router"
	"github.com/cruisebooking/backend/internal/service"

	_ "github.com/cruisebooking/backend/docs" // Swagger generated docs
)

// @title CruiseBooking API
// @version 1.0
// @description Backend API for CruiseBooking
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load configuration (env vars override config.yaml)
	cfg := config.Load(".")

	// 2. Initialize logger
	appLogger := logger.New(cfg.Log.Level, cfg.Log.Filename)
	defer func() { _ = appLogger.Sync() }()

	// 3. Connect database
	db, err := database.Connect(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// 4. Initialize repositories
	staffRepo := repository.NewStaffRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	cruiseRepo := repository.NewCruiseRepository(db)
	cabinTypeRepo := repository.NewCabinTypeRepository(db)
	facilityCategoryRepo := repository.NewFacilityCategoryRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)

	// 5. Initialize services
	authSvc := service.NewAuthService(staffRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	companySvc := service.NewCompanyService(companyRepo, cruiseRepo)
	cruiseSvc := service.NewCruiseService(cruiseRepo, cabinTypeRepo, companyRepo)
	cabinTypeSvc := service.NewCabinTypeService(cabinTypeRepo)
	facilityCategorySvc := service.NewFacilityCategoryService(facilityCategoryRepo)
	facilitySvc := service.NewFacilityService(facilityRepo)

	// 6. Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	cruiseHandler := handler.NewCruiseHandler(cruiseSvc)
	cabinTypeHandler := handler.NewCabinTypeHandler(cabinTypeSvc)
	facilityCategoryHandler := handler.NewFacilityCategoryHandler(facilityCategorySvc)
	facilityHandler := handler.NewFacilityHandler(facilitySvc)
	uploadHandler := handler.NewUploadHandler()

	// 7. Initialize Casbin enforcer
	enforcer, err := casbinv2.NewEnforcer("rbac/model.conf", "rbac/policy.csv")
	if err != nil {
		log.Fatalf("casbin enforcer init failed: %v", err)
	}

	// 8. Setup router and start server
	r := router.Setup(router.Dependencies{
		Auth:             authHandler,
		Company:          companyHandler,
		Cruise:           cruiseHandler,
		CabinType:        cabinTypeHandler,
		FacilityCategory: facilityCategoryHandler,
		Facility:         facilityHandler,
		Upload:           uploadHandler,
		JWTSecret:        cfg.JWT.Secret,
		Enforcer:         enforcer,
	})

	log.Printf("server starting on %s (mode: %s)", cfg.Server.Port, cfg.Server.Mode)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
