package user

import (
	"go-library/internal/auth"
	"go-library/internal/config"
	"go-library/internal/constants"
	middlewares "go-library/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, jwtService auth.JWTService) {
	userRepository := NewRepository(db)
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/refresh", userHandler.Refresh)
	api.GET("/me", userHandler.GetMe, middlewares.AuthMiddleware(jwtService))
	api.POST("/admin", userHandler.CreateAdmin,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleSuperAdmin)),
	)
	api.DELETE("/users/:id", userHandler.DeleteUser,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
	api.GET("/users", userHandler.GetAllUsers,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
}
