package dashboard

import (
	"go-library/internal/httpresponse"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) GetStats(c *echo.Context) error {
	result, err := h.service.GetStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch dashboard stats",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *handler) GetUserStats(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot identify user",
			Details: "missing user id in context",
		})
	}

	result, err := h.service.GetUserStats(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch user stats",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}
