package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Dependencies groups all handler dependencies for injection.
type Dependencies struct {
	Auth             *handler.AuthHandler
	Company          *handler.CompanyHandler
	Cruise           *handler.CruiseHandler
	CabinType        *handler.CabinTypeHandler
	FacilityCategory *handler.FacilityCategoryHandler
	Facility         *handler.FacilityHandler
	Upload           *handler.UploadHandler
	JWTSecret        string
	Enforcer         *casbin.Enforcer
}

// Setup creates and configures the Gin engine with all routes and middleware.
// CR-04 FIX: admin routes are protected by JWT+RBAC middleware; all handlers are properly injected.
func Setup(deps Dependencies) *gin.Engine {
	r := gin.New()

	// Global middleware: Recovery + Logger
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Swagger UI (unauthenticated)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// --- Public routes (no auth required) ---
	auth := api.Group("/admin/auth")
	{
		auth.POST("/login", deps.Auth.Login)
	}

	// --- Protected admin routes (JWT + RBAC required) ---
	admin := api.Group("/admin")
	admin.Use(middleware.JWT(middleware.JWTConfig{Secret: deps.JWTSecret}))
	if deps.Enforcer != nil {
		admin.Use(middleware.RBAC(deps.Enforcer))
	}

	// Auth profile (authenticated, any role)
	admin.GET("/auth/profile", deps.Auth.GetProfile)

	// Company management
	companies := admin.Group("/companies")
	{
		companies.GET("", deps.Company.List)
		companies.POST("", deps.Company.Create)
		companies.PUT("/:id", deps.Company.Update)
		companies.DELETE("/:id", deps.Company.Delete)
	}

	// Cruise management
	cruises := admin.Group("/cruises")
	{
		cruises.GET("", deps.Cruise.List)
		cruises.POST("", deps.Cruise.Create)
		cruises.PUT("/:id", deps.Cruise.Update)
		cruises.DELETE("/:id", deps.Cruise.Delete)
	}

	// Cabin type management
	cabinTypes := admin.Group("/cabin-types")
	{
		cabinTypes.GET("", deps.CabinType.List)
		cabinTypes.POST("", deps.CabinType.Create)
		cabinTypes.PUT("/:id", deps.CabinType.Update)
		cabinTypes.DELETE("/:id", deps.CabinType.Delete)
	}

	// Facility category management
	facilityCategories := admin.Group("/facility-categories")
	{
		facilityCategories.GET("", deps.FacilityCategory.List)
		facilityCategories.POST("", deps.FacilityCategory.Create)
		facilityCategories.DELETE("/:id", deps.FacilityCategory.Delete)
	}

	// Facility management
	facilities := admin.Group("/facilities")
	{
		facilities.GET("", deps.Facility.ListByCruise)
		facilities.POST("", deps.Facility.Create)
		facilities.DELETE("/:id", deps.Facility.Delete)
	}

	// Upload
	upload := admin.Group("/upload")
	{
		upload.POST("/image", deps.Upload.UploadImage)
	}

	return r
}
