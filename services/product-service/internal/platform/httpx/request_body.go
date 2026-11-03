package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"product-service/internal/platform/apperror"
)

const maxBodyBytes = 1 << 20 // 1 MB

// Request is implemented by every JSON request body type. Validate checks
// the business rules of the decoded values and returns an apperror.Invalid
// error (HTTP 422) when they are broken.
type Request interface {
	Validate() error
}

// BindJSON reads, decodes and validates the request body into request.
// On failure it writes the error response and returns false, so a handler
// only needs:
//
//	if !httpx.BindJSON(c, &request) {
//		return
//	}
func BindJSON(c *gin.Context, request Request) bool {
	body, err := ReadBody(c)
	if err != nil {
		WriteError(c, err)
		return false
	}
	if err := DecodeJSON(body, request); err != nil {
		WriteError(c, err)
		return false
	}
	return true
}

// ReadBody reads the whole request body, refusing bodies larger than 1 MB.
func ReadBody(c *gin.Context) ([]byte, error) {
	limitedBody := http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)
	body, err := io.ReadAll(limitedBody)
	if err == nil {
		return body, nil
	}

	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return nil, apperror.TooLarge(fmt.Sprintf("request body must not exceed %d bytes", maxBodyBytes))
	}
	return nil, apperror.Malformed("could not read request body")
}

// DecodeJSON strictly decodes exactly one JSON object into request, then calls
// request.Validate().
//
//   - Input that cannot be parsed (bad JSON, unknown field, wrong type, not an
//     object, extra data after the object) returns apperror.Malformed (400).
//   - Parsed input that breaks a rule returns the error from Validate (422).
func DecodeJSON(body []byte, request Request) error {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return apperror.Malformed("request body must be a JSON object")
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(request); err != nil {
		return apperror.Malformed(describeDecodeError(err))
	}

	// A second Decode must hit end of input, otherwise there is trailing data.
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return apperror.Malformed("request body must contain exactly one JSON object")
	}

	return request.Validate()
}

// describeDecodeError turns an encoding/json error into a client-friendly message.
func describeDecodeError(err error) string {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
		return "request body is not valid JSON"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		if typeErr.Field == "" {
			return "request body has a value of the wrong type"
		}
		return fmt.Sprintf("%s must be of type %s", typeErr.Field, jsonTypeName(typeErr.Type))
	}

	// encoding/json has no error type for unknown fields, only this message format.
	const unknownFieldPrefix = "json: unknown field "
	if strings.HasPrefix(err.Error(), unknownFieldPrefix) {
		return "unknown field " + strings.TrimPrefix(err.Error(), unknownFieldPrefix)
	}

	return "invalid request body"
}

// jsonTypeName returns the JSON name of the Go type the decoder expected.
func jsonTypeName(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		return "array"
	default:
		return "object"
	}
}
