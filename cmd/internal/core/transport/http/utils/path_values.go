package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	param := r.PathValue(key)
	if param == "" {
		return 0, fmt.Errorf("no key='%s' in path values: %w ", key, core_errors.ErrInvalidArgument)
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("param %s of key=%s not a valid integer: %v: %w", param, key, err, core_errors.ErrInvalidArgument)
	}

	return val, nil
}
