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

	_ "github.com/cruisebooking/backend/docs" // 导入 Swagger 自动生成的文档
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
	// 1. 加载配置（环境变量可覆盖 config.yaml 中的配置项）
	cfg := config.Load(".")

	// 2. 初始化日志记录器
	appLogger := logger.New(cfg.Log.Level, cfg.Log.Filename)
	defer func() { _ = appLogger.Sync() }()

	// 3. 连接数据库
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
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 4. 初始化数据仓储层
	staffRepo := repository.NewStaffRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	cruiseRepo := repository.NewCruiseRepository(db)
	cabinTypeRepo := repository.NewCabinTypeRepository(db)
	facilityCategoryRepo := repository.NewFacilityCategoryRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)

	// 5. 初始化业务服务层
	authSvc := service.NewAuthService(staffRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	companySvc := service.NewCompanyService(companyRepo, cruiseRepo)
	cruiseSvc := service.NewCruiseService(cruiseRepo, cabinTypeRepo, companyRepo)
	cabinTypeSvc := service.NewCabinTypeService(cabinTypeRepo)
	facilityCategorySvc := service.NewFacilityCategoryService(facilityCategoryRepo)
	facilitySvc := service.NewFacilityService(facilityRepo)

	// 6. 初始化 HTTP 处理器层
	authHandler := handler.NewAuthHandler(authSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	cruiseHandler := handler.NewCruiseHandler(cruiseSvc)
	cabinTypeHandler := handler.NewCabinTypeHandler(cabinTypeSvc)
	facilityCategoryHandler := handler.NewFacilityCategoryHandler(facilityCategorySvc)
	facilityHandler := handler.NewFacilityHandler(facilitySvc)
	uploadHandler := handler.NewUploadHandler()

	// 7. 初始化 Casbin RBAC 权限执行器
	enforcer, err := casbinv2.NewEnforcer("rbac/model.conf", "rbac/policy.csv")
	if err != nil {
		log.Fatalf("Casbin 执行器初始化失败: %v", err)
	}

	// 8. 配置路由并启动 HTTP 服务器
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

	log.Printf("服务启动于 %s（模式: %s）", cfg.Server.Port, cfg.Server.Mode)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
