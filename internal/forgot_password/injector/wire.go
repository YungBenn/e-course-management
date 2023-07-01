//go:build wireinject
// +build wireinject

package forgot_password

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	handler "e-course-management/internal/forgot_password/delivery/http"
	repository "e-course-management/internal/forgot_password/repository"
	usecase "e-course-management/internal/forgot_password/usecase"
	userRepository "e-course-management/internal/user/repository"
	userUseCase "e-course-management/internal/user/usecase"
	mail "e-course-management/pkg/mail/sendgrid"
)

func InitializedService(db *gorm.DB) *handler.ForgotPasswordHandler {
	wire.Build(
		handler.NewForgotPasswordHandler,
		repository.NewForgotPasswordRepository,
		usecase.NewForgotPasswordUseCase,
		userRepository.NewUserRepository,
		userUseCase.NewUserUseCase,
		mail.NewMailUseCase,
	)
	return &handler.ForgotPasswordHandler{}
}