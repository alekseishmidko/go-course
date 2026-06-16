package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)


func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param:= r.URL.Query().Get(key)
	if param == "" {
		return nil , nil
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("param %s of key=%s not a valid integer: %v: %w", param, key, err, core_errors.ErrInvalidArgument)
	}

	return &val, nil
}