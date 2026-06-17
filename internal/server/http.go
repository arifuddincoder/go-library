package server

import (
	"fmt"
	"go-library/internal/auth"
	"go-library/internal/config"
	"go-library/internal/domain/book"
	"go-library/internal/domain/category"
	"go-library/internal/domain/dashboard"
	"go-library/internal/domain/loan"
	"go-library/internal/domain/user"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	// ek bar e jwtService banai, sob domain ke pass kori
	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello from go library")
	})

	user.RegisterRoutes(e, db, cfg, jwtService)
	category.RegisterRoutes(e, db, jwtService)
	book.RegisterRoutes(e, db, jwtService)
	loan.RegisterRoutes(e, db, jwtService)
	dashboard.RegisterRoutes(e, db, jwtService)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
