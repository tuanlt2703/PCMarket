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
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"product-service/internal/platform/apperror"
	"product-service/internal/platform/optional"
)

const maxBodyBytes = 1 << 20

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("httpx: unexpected validator engine")
	}
	// Report JSON field names in validation messages.
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		if name == "" {
			return f.Name
		}
		return name
	})
	v.RegisterCustomTypeFunc(func(f reflect.Value) any {
		return f.Interface().(interface{ ValidationValue() any }).ValidationValue()
	},
		optional.Field[string]{}, optional.Field[bool]{},
		optional.Field[int32]{}, optional.Field[int64]{}, optional.Field[float64]{},
		optional.Field[[]int64]{},
	)
}

// ReadBody reads the request body up to the size limit.
func ReadBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return nil, apperror.TooLarge(fmt.Sprintf("request body must not exceed %d bytes", maxBodyBytes))
		}
		return nil, apperror.Malformed("could not read request body")
	}
	return body, nil
}

// DecodeJSON strictly decodes exactly one JSON object into dst and validates it.
// Unparseable input returns ErrMalformed (400); failed binding tags or a
// Validate() error method on dst return ErrInvalid (422).
func DecodeJSON(body []byte, dst any) error {
	if trimmed := bytes.TrimSpace(body); len(trimmed) == 0 || trimmed[0] != '{' {
		return apperror.Malformed("request body must be a JSON object")
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return apperror.Malformed(decodeMessage(err))
	}
	if err := dec.Decode(new(any)); !errors.Is(err, io.EOF) {
		return apperror.Malformed("request body must contain exactly one JSON object")
	}
	if err := binding.Validator.ValidateStruct(dst); err != nil {
		return apperror.Invalid(validationMessage(err))
	}
	if v, ok := dst.(interface{ Validate() error }); ok {
		return v.Validate()
	}
	return nil
}

func decodeMessage(err error) string {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		return "request body is not valid JSON"
	case errors.As(err, &typeErr):
		if typeErr.Field != "" {
			return fmt.Sprintf("%s must be of type %s", typeErr.Field, jsonType(typeErr.Type))
		}
		return "request body has a value of the wrong type"
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		return "unknown field " + strings.TrimPrefix(err.Error(), "json: unknown field ")
	}
	return "invalid request body"
}

func jsonType(t reflect.Type) string {
	switch t.Kind() {
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
	}
	return "object"
}

func validationMessage(err error) string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return "invalid request body"
	}
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, fieldMessage(e))
	}
	return strings.Join(msgs, "; ")
}

func fieldMessage(e validator.FieldError) string {
	field := e.Namespace()
	if i := strings.IndexByte(field, '.'); i >= 0 {
		field = field[i+1:] // drop the root struct name
	}
	isString := e.Kind() == reflect.String
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "max", "lte":
		if isString {
			return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
		}
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "min", "gte":
		if isString {
			return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		}
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "datetime":
		return fmt.Sprintf("%s must be a date in the format YYYY-MM-DD", field)
	}
	return fmt.Sprintf("%s is invalid (%s)", field, e.Tag())
}
