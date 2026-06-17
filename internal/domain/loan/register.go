package loan

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
	loanRepository := NewRepository(db)
	loanService := NewService(loanRepository)
	loanHandler := NewHandler(loanService)

	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	adminOnly := middlewares.RequireRole(string(constants.RoleAdmin), string(constants.RoleSuperAdmin))

	api := e.Group("/api/v1/loans")
	api.Use(middlewares.AuthMiddleware(jwtService))

	// user
	api.POST("", loanHandler.RequestLoan)
	api.GET("/me/active", loanHandler.GetMyActiveLoans)
	api.GET("/me/history", loanHandler.GetMyLoanHistory)

	// admin
	api.GET("", loanHandler.GetAllLoans, adminOnly)
	api.PATCH("/:id/approve", loanHandler.ApproveLoan, adminOnly)
	api.PATCH("/:id/reject", loanHandler.RejectLoan, adminOnly)
	api.PATCH("/:id/return", loanHandler.ReturnLoan, adminOnly)
}
