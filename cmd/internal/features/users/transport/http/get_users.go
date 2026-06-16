package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
	core_http_utils "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDtoResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	limit, offset, error := getLimitOffsetQueryParams(r)

	if error != nil {
		responseHandler.ErrorResponse(error, "failed to get 'limit' or 'offset' from query params")

		return
	}
	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDtoFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")

	if err != nil {
		return nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")

	if err != nil {
		return nil, nil, fmt.Errorf("get offset query param: %w", err)
	}
	return limit, offset, nil
}
