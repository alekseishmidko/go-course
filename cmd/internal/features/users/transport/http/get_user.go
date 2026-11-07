package users_transport_http

import (
	"net/http"

	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
)

type GetUserResponse UserDtoResponse

// GetUser godoc
// @Summary Получить пользователя
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} core_http_response.ErrorResponse
// @Failure 404 {object} core_http_response.ErrorResponse
// @Failure 500 {object} core_http_response.ErrorResponse
// @Router /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId by path value")
		return
	}

	user, err := h.usersService.GetUser(ctx, userId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDtoFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
