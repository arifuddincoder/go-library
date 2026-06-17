package book

import (
	"go-library/internal/auth"
	"go-library/internal/constants"
	middlewares "go-library/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	bookRepository := NewRepository(db)
	bookService := NewService(bookRepository)
	bookHandler := NewHandler(bookService)

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
