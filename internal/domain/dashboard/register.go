package dashboard

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
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	api := e.Group("/api/v1/dashboard")
	api.GET("/stats", h.GetStats,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
	api.GET("/me/stats", h.GetUserStats, middlewares.AuthMiddleware(jwtService))
}
