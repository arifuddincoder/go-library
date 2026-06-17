package book

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
	bookRepository := NewRepository(db)
	bookService := NewService(bookRepository)
	bookHandler := NewHandler(bookService)

	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	api := e.Group("/api/v1/books")
	api.GET("", bookHandler.GetAllBooks)
	api.GET("/:id", bookHandler.GetBookByID)

	api.POST("", bookHandler.CreateBook,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
	api.PATCH("/:id", bookHandler.UpdateBook,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
	api.DELETE("/:id", bookHandler.DeleteBook,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin)),
	)
}
