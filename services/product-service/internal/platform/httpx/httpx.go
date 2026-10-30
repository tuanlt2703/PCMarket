package httpx

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/apperror"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError maps application errors to HTTP status codes. Unknown errors are
// logged and hidden from the client.
func WriteError(c *gin.Context, err error) {
	var appErr *apperror.Error
	msg := "internal server error"
	if errors.As(err, &appErr) {
		msg = appErr.Message
	}

	switch {
	case errors.Is(err, apperror.ErrNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{msg})
	case errors.Is(err, apperror.ErrConflict):
		c.JSON(http.StatusConflict, ErrorResponse{msg})
	case errors.Is(err, apperror.ErrInvalid):
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{msg})
	case errors.Is(err, apperror.ErrMalformed):
		c.JSON(http.StatusBadRequest, ErrorResponse{msg})
	case errors.Is(err, apperror.ErrTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{msg})
	default:
		slog.ErrorContext(c.Request.Context(), "request failed", "path", c.FullPath(), "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{msg})
	}
}

// BindJSON reads, strictly decodes and validates the body, writing the error
// response on failure: 400 for unparseable JSON, 422 for failed validation.
func BindJSON(c *gin.Context, dst any) bool {
	body, err := ReadBody(c)
	if err == nil {
		err = DecodeJSON(body, dst)
	}
	if err != nil {
		WriteError(c, err)
		return false
	}
	return true
}

// ParamID parses a positive int64 path parameter, writing a 400 response on failure.
func ParamID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{"invalid " + name})
		return 0, false
	}
	return id, true
}

// BadRequest writes a 400 response with the given message.
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{msg})
}
