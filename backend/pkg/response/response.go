package response

import (
	"errors"
	"net/http"

	"epbms/internal/domain"
	"github.com/gin-gonic/gin"
)

// Envelope is the standard API response wrapper.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// OK sends a 200 JSON response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

// OKWithMeta sends a 200 JSON response with pagination metadata.
func OKWithMeta(c *gin.Context, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data, Meta: &meta})
}

// Created sends a 201 JSON response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

// Error maps a domain error to the appropriate HTTP status code and sends the response.
func Error(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, domain.ErrNotFound):
		code = http.StatusNotFound
		msg = err.Error()
	case errors.Is(err, domain.ErrConflict):
		code = http.StatusConflict
		msg = err.Error()
	case errors.Is(err, domain.ErrBookingConflict):
		code = http.StatusConflict
		msg = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		code = http.StatusUnauthorized
		msg = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		code = http.StatusForbidden
		msg = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		code = http.StatusBadRequest
		msg = err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials):
		code = http.StatusUnauthorized
		msg = err.Error()
	case errors.Is(err, domain.ErrInvalidStatus):
		code = http.StatusBadRequest
		msg = err.Error()
	}

	c.JSON(code, Envelope{Success: false, Error: msg})
}

// BadRequest sends a 400 JSON response with a custom message.
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Envelope{Success: false, Error: msg})
}
