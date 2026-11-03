package httpx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/apperror"
)

// ErrorResponse is the JSON body of every error response: {"error": "..."}.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError converts an error into an HTTP response.
//
// Errors created with the apperror package are mapped to a 4xx status and
// their message is shown to the client. Any other error is unexpected: it is
// logged and the client only sees "internal server error".
func WriteError(c *gin.Context, err error) {
	message := "internal server error"
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		message = appErr.Message
	}

	status := statusFor(err)
	if status == http.StatusInternalServerError {
		slog.ErrorContext(c.Request.Context(), "request failed", "path", c.FullPath(), "error", err)
	}
	c.JSON(status, ErrorResponse{Error: message})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperror.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, apperror.ErrInvalid):
		return http.StatusUnprocessableEntity
	case errors.Is(err, apperror.ErrMalformed):
		return http.StatusBadRequest
	case errors.Is(err, apperror.ErrTooLarge):
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

// WriteBadRequest writes a 400 response with the given message.
func WriteBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{Error: message})
}
