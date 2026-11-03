package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// The functions in this file read path and query parameters.
// When a value is invalid they write a 400 response and return ok == false,
// so the handler only has to return:
//
//	id, ok := httpx.ParamID(c, "id")
//	if !ok {
//		return
//	}

// ParamID parses a positive int64 path parameter such as /models/:id.
func ParamID(c *gin.Context, name string) (id int64, ok bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		WriteBadRequest(c, "invalid "+name)
		return 0, false
	}
	return id, true
}

// QueryBool parses an optional boolean query parameter such as ?is_active=true.
// It returns nil when the parameter is absent or empty.
func QueryBool(c *gin.Context, name string) (value *bool, ok bool) {
	raw := c.Query(name)
	if raw == "" {
		return nil, true
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		WriteBadRequest(c, "invalid "+name)
		return nil, false
	}
	return &parsed, true
}

// QueryID parses an optional positive int64 query parameter such as
// ?manufacturer_id=3. It returns nil when the parameter is absent or empty.
func QueryID(c *gin.Context, name string) (value *int64, ok bool) {
	raw := c.Query(name)
	if raw == "" {
		return nil, true
	}
	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || parsed <= 0 {
		WriteBadRequest(c, "invalid "+name)
		return nil, false
	}
	return &parsed, true
}
