package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/config"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/cruisebooking/backend/internal/pkg/database"
	"github.com/cruisebooking/backend/internal/pkg/logger"
	"github.com/cruisebooking/backend/internal/pkg/search"
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
var osExit = os.Exit

// main 为服务进程入口。
func main() {
	if err := RunApp("./"); err != nil {
		log.Printf("服务启动失败: %v", err)
		osExit(1)
	}
}

// RunApp 包含了应用程序的启动逻辑，提取出来以便单元测试覆盖。
func RunApp(configDir string) error {
	// 1. 加载配置（环境变量可覆盖 config.yaml 中的配置项）
	cfg := config.Load(configDir)

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
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 4. 初始化数据仓储层
	staffRepo := repository.NewStaffRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	cruiseRepo := repository.NewCruiseRepository(db)
	cabinTypeRepo := repository.NewCabinTypeRepository(db)
	facilityCategoryRepo := repository.NewFacilityCategoryRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)
	routeRepo := repository.NewRouteRepository(db)
	voyageRepo := repository.NewVoyageRepository(db)
	cabinRepo := repository.NewCabinRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 5. 初始化业务服务层
	authSvc := service.NewAuthService(staffRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	companySvc := service.NewCompanyService(companyRepo, cruiseRepo)
	cruiseSvc := service.NewCruiseService(cruiseRepo, cabinTypeRepo, companyRepo)
	cabinTypeSvc := service.NewCabinTypeService(cabinTypeRepo)
	facilityCategorySvc := service.NewFacilityCategoryService(facilityCategoryRepo)
	facilitySvc := service.NewFacilityService(facilityRepo)
	pricingSvc := service.NewPricingService(cabinRepo)
	cabinAdminSvc := service.NewCabinAdminService(cabinRepo)
	meiliIndexer := search.NewMeiliIndexer(cfg.Meilis.Host, cfg.Meilis.APIKey)
	searchRetryQueue := service.NewSearchRetryQueue(meiliIndexer, 3, 128)
	searchRetryQueue.Start()
	holdRepo := repository.NewCabinHoldRepository(db)
	holdSvc := service.NewCabinHoldService(holdRepo, 15*time.Minute)

	// 6. 初始化 HTTP 处理器层
	authHandler := handler.NewAuthHandler(authSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	cruiseHandler := handler.NewCruiseHandler(cruiseSvc)
	cabinTypeHandler := handler.NewCabinTypeHandler(cabinTypeSvc)
	facilityCategoryHandler := handler.NewFacilityCategoryHandler(facilityCategorySvc)
	facilityHandler := handler.NewFacilityHandler(facilitySvc)
	uploadHandler := handler.NewUploadHandler()
	routeHandler := handler.NewRouteHandler(routeRepo)    // L-02: 直接用 repo 满足 RouteService 接口
	voyageHandler := handler.NewVoyageHandler(voyageRepo) // L-02: 直接用 repo 满足 VoyageService 接口
	cabinHandler := handler.NewCabinHandlerWithIndexing(cabinAdminSvc, meiliIndexer, searchRetryQueue)

	bookingRepo := repository.NewBookingRepository(db)
	bookingSvc := service.NewBookingService(bookingRepo, pricingSvc, holdSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc, bookingRepo)
	userAuthSvc := service.NewUserAuthService(service.NewInMemoryCodeStore())
	userHandler := handler.NewUserHandlerWithRepo(userAuthSvc, userRepo, cfg.JWT.Secret) // M-03

	// Sprint 04: 支付 / 退款 / 通知 / 统计分析 依赖注入
	paymentRepo := repository.NewPaymentRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	payVerifiers := map[string]service.PaymentVerifier{
		"wechat": service.NewHMACVerifier(cfg.JWT.Secret),
		"alipay": service.NewHMACVerifier(cfg.JWT.Secret),
	}
	payCallbackSvc := service.NewPaymentCallbackService(paymentRepo, bookingRepo, payVerifiers)
	refundSvc := service.NewRefundService(paymentRepo, refundRepo)
	notifySvc := service.NewNotifyService(notifRepo)
	_ = notifySvc // 通知服务供预订/支付流程调用
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)

	paymentHandler := handler.NewPaymentHandler(payCallbackSvc)
	refundHandler := handler.NewRefundHandler(refundSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)

	// 7. 初始化 Casbin RBAC 权限执行器
	mPath := filepath.Join(configDir, "rbac/model.conf")
	pPath := filepath.Join(configDir, "rbac/policy.csv")
	enforcer, err := casbinv2.NewEnforcer(mPath, pPath)
	if err != nil {
		return fmt.Errorf("Casbin 执行器初始化失败: %w", err)
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
		Route:            routeHandler,
		Voyage:           voyageHandler,
		Cabin:            cabinHandler,
		Booking:          bookingHandler,
		User:             userHandler,
		Payment:          paymentHandler,
		Refund:           refundHandler,
		Analytics:        analyticsHandler,
		JWTSecret:        cfg.JWT.Secret,
		Enforcer:         enforcer,
	})

	log.Printf("服务启动于 %s（模式: %s）", cfg.Server.Port, cfg.Server.Mode)
	return r.Run(cfg.Server.Port)
}
