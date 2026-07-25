package users_transport_http

import (
	"net/http"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
	core_http_utils "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/utils"
)

type GetUserResponse UserDtoResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId by path value")
	}

	user, err := h.usersService.GetUser(ctx, userId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDtoFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
