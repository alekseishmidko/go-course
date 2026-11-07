package users_transport_http

import (
	"net/http"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_logger "github.com/alekseishmidko/go-course/cmd/internal/core/logger"
	core_http_request "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/request"
	core_http_response "github.com/alekseishmidko/go-course/cmd/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" example:"Алексей Иванов" minLength:"3" maxLength:"100" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" example:"+79997654321" minLength:"10" maxLength:"15" extensions:"x-nullable" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDtoResponse

// CreateUser godoc
// @Summary Создать пользователя
// @Tags users
// @Produce json
// @Accept json
// @Param request body CreateUserRequest true "Поля новой записи"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} core_http_response.ErrorResponse
// @Failure 500 {object} core_http_response.ErrorResponse
// @Router /users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {

		responseHandler.ErrorResponse(err, "validation Failed")
		return
	}

	userDomain := domainFromDto(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "CreateUser Failed")
		return
	}
	response := CreateUserResponse(userDtoFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDto(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
