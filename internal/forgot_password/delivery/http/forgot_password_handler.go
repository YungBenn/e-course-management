package forgot_password

import (
	"net/http"

	dto "e-course-management/internal/forgot_password/dto"
	usecase "e-course-management/internal/forgot_password/usecase"
	"e-course-management/pkg/response"
	"github.com/gin-gonic/gin"
)

type ForgotPasswordHandler struct {
	usecase usecase.ForgotPasswordUseCase
}

func NewForgotPasswordHandler(usecase usecase.ForgotPasswordUseCase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{usecase}
}

func (handler *ForgotPasswordHandler) Route(r *gin.RouterGroup) {
	forgotPasswordRouter := r.Group("/api/v1")

	forgotPasswordRouter.POST("/forgot_passwords", handler.Create)
	forgotPasswordRouter.PUT("/forgot_passwords", handler.Update)
}

func (handler *ForgotPasswordHandler) Create(ctx *gin.Context) {
	var input dto.ForgotPasswordRequestBody

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		))
		ctx.Abort()
		return
	}

	_, err := handler.usecase.Create(input)

	if err != nil {
		ctx.JSON(int(err.Code), response.Response(
			int(err.Code),
			http.StatusText(int(err.Code)),
			err.Err.Error(),
		))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response.Response(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		"Success, please check your email.",
	))
}

func (handler *ForgotPasswordHandler) Update(ctx *gin.Context) {
	var input dto.ForgotPasswordUpdateRequestBody

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		))
		ctx.Abort()
		return
	}

	_, err := handler.usecase.Update(input)

	if err != nil {
		ctx.JSON(int(err.Code), response.Response(
			int(err.Code),
			http.StatusText(int(err.Code)),
			err.Err.Error(),
		))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response.Response(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		"Success change your password",
	))
}
