package user

import (
	"go-library/internal/auth"
	"go-library/internal/config"
	"go-library/internal/constants"
	middlewares "go-library/internal/middleware"
	"log"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.GET("/me", userHandler.GetMe, middlewares.AuthMiddleware(jwtService))
	api.POST("/admin", userHandler.CreateAdmin,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleSuperAdmin)),
	)
	api.DELETE("/users/:id", userHandler.DeleteUser,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
}
