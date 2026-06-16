package book

import (
	"errors"
	"go-library/internal/domain/book/dto"
	"go-library/internal/httpresponse"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateBook(c *echo.Context) error {
	var req dto.RegisterRequest

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

	response, err := h.service.CreateBook(req)
	if err != nil {
		if errors.Is(err, ErrBookAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Failed to create book",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create book",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetBookByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.PathParam("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid book id",
			Details: err.Error(),
		})
	}

	response, err := h.service.GetBookByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Book not found",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get book",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetAllBooks(c *echo.Context) error {
	responses, err := h.service.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get books",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *handler) DeleteBook(c *echo.Context) error {
	id, err := strconv.ParseUint(c.PathParam("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid book id",
			Details: err.Error(),
		})
	}

	if err := h.service.DeleteBook(uint(id)); err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Book not found",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete book",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
