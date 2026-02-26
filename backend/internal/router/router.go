package router

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Dependencies 聚合了所有处理器依赖，用于路由初始化时的依赖注入。
type Dependencies struct {
	Auth             *handler.AuthHandler             // 认证处理器
	Company          *handler.CompanyHandler          // 邮轮公司处理器
	Cruise           *handler.CruiseHandler           // 邮轮处理器
	CabinType        *handler.CabinTypeHandler        // 舱房类型处理器
	FacilityCategory *handler.FacilityCategoryHandler // 设施分类处理器
	Facility         *handler.FacilityHandler         // 设施处理器
	Route            *handler.RouteHandler            // 航线处理器
	Voyage           *handler.VoyageHandler           // 航次处理器
	Cabin            *handler.CabinHandler            // 舱房处理器
	Booking          *handler.BookingHandler          // 订单处理器
	User             *handler.UserHandler             // C端用户处理器
	Upload           *handler.UploadHandler           // 文件上传处理器
	Payment          *handler.PaymentHandler          // 支付回调处理器
	Refund           *handler.RefundHandler           // 退款处理器
	Analytics        *handler.AnalyticsHandler        // 统计分析处理器
	JWTSecret        string                           // JWT 签名密钥
	Enforcer         *casbin.Enforcer                 // Casbin RBAC 执行器
}

// Setup 创建并配置 Gin 引擎，注册所有路由和中间件。
// CR-04 修复：管理后台路由受 JWT + RBAC 中间件保护；所有处理器均通过依赖注入传入。
func Setup(deps Dependencies) *gin.Engine {
	r := gin.New()

	// 全局中间件：崩溃恢复 + 请求日志
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger API 文档界面（无需认证）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// --- 公开路由（无需认证） ---
	auth := api.Group("/admin/auth")
	{
		auth.POST("/login", deps.Auth.Login) // 管理员登录
	}

	// --- 受保护的管理后台路由（需要 JWT + RBAC 认证） ---
	admin := api.Group("/admin")
	admin.Use(middleware.JWT(middleware.JWTConfig{Secret: deps.JWTSecret}))
	if deps.Enforcer != nil {
		admin.Use(middleware.RBAC(deps.Enforcer))
	}

	// 获取当前员工信息（已认证，任意角色）
	admin.GET("/auth/profile", deps.Auth.GetProfile)

	// 邮轮公司管理
	companies := admin.Group("/companies")
	{
		companies.GET("", deps.Company.List)          // 查询公司列表
		companies.POST("", deps.Company.Create)       // 创建公司
		companies.PUT("/:id", deps.Company.Update)    // 更新公司
		companies.DELETE("/:id", deps.Company.Delete) // 删除公司
	}

	// 邮轮管理
	cruises := admin.Group("/cruises")
	{
		cruises.GET("", deps.Cruise.List)          // 查询邮轮列表
		cruises.POST("", deps.Cruise.Create)       // 创建邮轮
		cruises.PUT("/:id", deps.Cruise.Update)    // 更新邮轮
		cruises.DELETE("/:id", deps.Cruise.Delete) // 删除邮轮
	}

	// 舱房类型管理
	cabinTypes := admin.Group("/cabin-types")
	{
		cabinTypes.GET("", deps.CabinType.List)          // 查询舱房类型列表
		cabinTypes.POST("", deps.CabinType.Create)       // 创建舱房类型
		cabinTypes.PUT("/:id", deps.CabinType.Update)    // 更新舱房类型
		cabinTypes.DELETE("/:id", deps.CabinType.Delete) // 删除舱房类型
	}

	// 设施分类管理
	facilityCategories := admin.Group("/facility-categories")
	{
		facilityCategories.GET("", deps.FacilityCategory.List)          // 查询分类列表
		facilityCategories.POST("", deps.FacilityCategory.Create)       // 创建分类
		facilityCategories.DELETE("/:id", deps.FacilityCategory.Delete) // 删除分类
	}

	// 设施管理
	facilities := admin.Group("/facilities")
	{
		facilities.GET("", deps.Facility.ListByCruise)  // 按邮轮查询设施
		facilities.POST("", deps.Facility.Create)       // 创建设施
		facilities.DELETE("/:id", deps.Facility.Delete) // 删除设施
	}

	// 文件上传
	upload := admin.Group("/upload")
	{
		upload.POST("/image", deps.Upload.UploadImage) // 上传图片
	}

	// 航线、航次、舱房管理（Sprint 2）—— 完整 CRUD
	routes := admin.Group("/routes")
	{
		routes.GET("", deps.Route.List)          // 查询航线列表
		routes.POST("", deps.Route.Create)       // 创建航线
		routes.PUT("/:id", deps.Route.Update)    // 更新航线
		routes.DELETE("/:id", deps.Route.Delete) // 删除航线
	}

	voyages := admin.Group("/voyages")
	{
		voyages.GET("", deps.Voyage.List)          // 查询航次列表
		voyages.POST("", deps.Voyage.Create)       // 创建航次
		voyages.PUT("/:id", deps.Voyage.Update)    // 更新航次
		voyages.DELETE("/:id", deps.Voyage.Delete) // 删除航次
	}

	cabins := admin.Group("/cabins")
	{
		cabins.GET("", deps.Cabin.List)                                  // 查询舱房列表
		cabins.POST("", deps.Cabin.Create)                               // 创建舱房 SKU
		cabins.PUT("/:id", deps.Cabin.Update)                            // 更新舱房 SKU
		cabins.DELETE("/:id", deps.Cabin.Delete)                         // 删除舱房 SKU
		cabins.GET("/:id/inventory", deps.Cabin.GetInventory)            // 查询库存
		cabins.POST("/:id/inventory/adjust", deps.Cabin.AdjustInventory) // 调整库存
		cabins.GET("/:id/prices", deps.Cabin.ListPrices)                 // 查询价格日历
		cabins.POST("/:id/prices", deps.Cabin.UpsertPrice)               // 设置价格
	}

	// ------------------------------------------
	// 小程序/Web C端 API（无需 admin 权限，部分需要 user auth）
	// C 端 JWT 使用独立 ContextKey（ContextKeyUserID）区分管理员身份
	// ------------------------------------------
	cUserJWT := middleware.JWT(middleware.JWTConfig{Secret: deps.JWTSecret, ContextKey: middleware.ContextKeyUserID})

	users := api.Group("/users")
	{
		users.POST("/login", deps.User.Login)
		users.POST("/sms-code", deps.User.SendCode)
		users.Use(cUserJWT)
		users.GET("/profile", deps.User.Profile)
	}

	bookings := api.Group("/bookings")
	{
		bookings.Use(cUserJWT)
		bookings.POST("", deps.Booking.Create)
	}

	// --- 支付回调（公开路由，由支付平台调用） ---
	api.POST("/pay/callback", deps.Payment.Callback)

	// --- 退款（需要用户认证） ---
	refunds := api.Group("/refunds")
	{
		refunds.Use(cUserJWT)
		refunds.POST("", deps.Refund.Create)
	}

	// --- 管理后台统计分析 ---
	admin.GET("/analytics/summary", deps.Analytics.Summary)

	return r
}
