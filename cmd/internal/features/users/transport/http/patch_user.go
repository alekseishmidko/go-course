package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
	core_http_types "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" example:"Алексей Иванов" minLength:"3" maxLength:"100" swaggertype:"string"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" example:"+79997654321" minLength:"10" maxLength:"15" extensions:"x-nullable" swaggertype:"string"`
}

type PatchUserResponse UserDtoResponse

// PatchUser godoc
// @Summary Обновить пользователя
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Accept json
// @Param request body PatchUserRequest true "Поля для обновления; непереданные поля сохраняются"
// @Description Пустой объект допустим. phone_number: null очищает телефон; full_name: null запрещён.
// @Success 200 {object} PatchUserResponse
// @Failure 400 {object} core_http_response.ErrorResponse
// @Failure 404 {object} core_http_response.ErrorResponse
// @Failure 409 {object} core_http_response.ErrorResponse
// @Failure 500 {object} core_http_response.ErrorResponse
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id from path value")
		return
	}
	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}
	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to patch user")
		return
	}
	response := PatchUserResponse(userDtoFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
	log.Debug(fmt.Sprintf("PatchUserRequest fields: FullName: '%+v'\n PhoneNumber:'%+v'", request.FullName, request.PhoneNumber))

}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can not be null")
		}
		fullNameLength := len([]rune(*r.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("`FullName` must be within 3 and 100 symbols")
		}

	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
		}
		if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
			return fmt.Errorf("PhoneNumber must starts with `+` symbol")
		}
	}
	return nil
}
func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {

	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())

}
