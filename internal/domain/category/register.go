package category

import (
	"go-library/internal/auth"
	"go-library/internal/constants"
	middlewares "go-library/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	categoryRepository := NewRepository(db)
	categoryService := NewService(categoryRepository)
	categoryHandler := NewHandler(categoryService)

	api := e.Group("/api/v1/categories")
	api.POST("", categoryHandler.CreateCategory, middlewares.AuthMiddleware(jwtService), middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)))
	api.GET("", categoryHandler.GetAllCategories, middlewares.AuthMiddleware(jwtService))
	api.DELETE("/:id", categoryHandler.DeleteCategory,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
}
