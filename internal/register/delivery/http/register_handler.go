package register

import (
	registerUseCase "e-course-management/internal/register/usecase"
	userDto "e-course-management/internal/user/dto"
	"e-course-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	registerUseCase registerUseCase.RegisterUseCase
}

func NewRegisterHandler(registerUseCase registerUseCase.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{registerUseCase}
}

func (handler *RegisterHandler) Route(r *gin.RouterGroup) {
	r.POST("/api/v1/registers", handler.Register)
}

func (handler *RegisterHandler) Register(ctx *gin.Context) {
	var registerRequestInput userDto.UserRequestBody

	if err := ctx.ShouldBindJSON(&registerRequestInput); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		))
		ctx.Abort()
		return
	}

	err := handler.registerUseCase.Register(registerRequestInput)

	if err != nil {
		ctx.JSON(int(err.Code), response.Response(
			int(err.Code),
			http.StatusText(int(err.Code)),
			err.Err.Error(),
		))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, response.Response(
		http.StatusCreated,
		http.StatusText(http.StatusCreated),
		"Success, please check your email",
	))
}
