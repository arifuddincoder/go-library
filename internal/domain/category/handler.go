package category

import (
	"errors"
	"go-library/internal/domain/category/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateCategory(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateCategory(req)
	if err != nil {
		if errors.Is(err, ErrCategoryAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Failed to create category",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create category",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetAllCategories(c *echo.Context) error {
	p := query.Parse(c)

	result, err := h.service.GetAllCategories(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch categories",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *handler) DeleteCategory(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid category id",
			Details: err.Error(),
		})
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Category not found",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete category",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
