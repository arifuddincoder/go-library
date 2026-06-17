package loan

import (
	"errors"
	"go-library/internal/domain/loan/dto"
	"go-library/internal/httpresponse"
	"go-library/internal/query"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) userID(c *echo.Context) (uint, bool) {
	id, ok := c.Get("user_id").(uint)
	return id, ok
}

func (h *handler) RequestLoan(c *echo.Context) error {
	userID, ok := h.userID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot identify user",
			Details: "missing user id in context",
		})
	}

	var req dto.RequestLoan
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

	response, err := h.service.RequestLoan(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code: http.StatusNotFound, Message: "Cannot create loan", Details: err.Error(),
			})
		case errors.Is(err, ErrNoCopiesLeft), errors.Is(err, ErrAlreadyRequested):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "Cannot create loan", Details: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code: http.StatusInternalServerError, Message: "Failed to create loan request", Details: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *handler) ApproveLoan(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid loan id",
			Details: err.Error(),
		})
	}

	response, err := h.service.ApproveLoan(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrLoanNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code: http.StatusNotFound, Message: "Loan not found", Details: err.Error(),
			})
		case errors.Is(err, ErrNotPending):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "Cannot approve loan", Details: err.Error(),
			})
		case errors.Is(err, ErrBookNotFound), errors.Is(err, ErrNoCopiesLeft):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "Cannot approve loan", Details: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code: http.StatusInternalServerError, Message: "Failed to approve loan", Details: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, response)
}

func (h *handler) RejectLoan(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code: http.StatusBadRequest, Message: "Invalid loan id", Details: err.Error(),
		})
	}

	response, err := h.service.RejectLoan(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrLoanNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code: http.StatusNotFound, Message: "Loan not found", Details: err.Error(),
			})
		case errors.Is(err, ErrNotPending):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "Cannot reject loan", Details: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code: http.StatusInternalServerError, Message: "Failed to reject loan", Details: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, response)
}

func (h *handler) ReturnLoan(c *echo.Context) error {
	id, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code: http.StatusBadRequest, Message: "Invalid loan id", Details: err.Error(),
		})
	}

	response, err := h.service.ReturnLoan(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrLoanNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code: http.StatusNotFound, Message: "Loan not found", Details: err.Error(),
			})
		case errors.Is(err, ErrNotBorrowed):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "Cannot return book", Details: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code: http.StatusInternalServerError, Message: "Failed to return book", Details: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, response)
}

// present loan — j user nisilo se nijer borrowed loan dekhbe
func (h *handler) GetMyActiveLoans(c *echo.Context) error {
	userID, ok := h.userID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code: http.StatusUnauthorized, Message: "Cannot identify user", Details: "missing user id in context",
		})
	}
	p := query.Parse(c)
	result, err := h.service.GetMyActiveLoans(userID, p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code: http.StatusInternalServerError, Message: "Failed to fetch active loans", Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}

// old loan — returned hoye geche emon gula
func (h *handler) GetMyLoanHistory(c *echo.Context) error {
	userID, ok := h.userID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code: http.StatusUnauthorized, Message: "Cannot identify user", Details: "missing user id in context",
		})
	}
	p := query.Parse(c)
	result, err := h.service.GetMyLoanHistory(userID, p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code: http.StatusInternalServerError, Message: "Failed to fetch loan history", Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}

// admin — sob loan, chaile ?status=pending diye filter
func (h *handler) GetAllLoans(c *echo.Context) error {
	p := query.Parse(c)

	var statuses []string
	if st := c.QueryParam("status"); st != "" {
		statuses = []string{st}
	}

	result, err := h.service.GetAllLoans(statuses, p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code: http.StatusInternalServerError, Message: "Failed to fetch loans", Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, result)
}
