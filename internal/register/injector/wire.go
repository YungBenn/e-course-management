//go:build wireinject
// +build wireinject

package register

import (
	handler "e-course-management/internal/register/delivery/http"
	registerUseCase "e-course-management/internal/register/usecase"
	userRepository "e-course-management/internal/user/repository"
	userUseCase "e-course-management/internal/user/usecase"
	mail "e-course-management/pkg/mail/sendgrid"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializedService(db *gorm.DB) *handler.RegisterHandler {
	wire.Build(
		registerUseCase.NewRegisterUseCase,
		handler.NewRegisterHandler,
		userRepository.NewUserRepository,
		userUseCase.NewUserUseCase,
		mail.NewMailUseCase,
	)

	return &handler.RegisterHandler{}
}
