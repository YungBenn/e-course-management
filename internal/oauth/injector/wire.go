//go:build wireinject
// +build wireinject

package oauth

import (
	handler "e-course-management/internal/oauth/delivery/http"
	oauthRepository "e-course-management/internal/oauth/repository"
	oauthUseCase "e-course-management/internal/oauth/usecase"
	userRepository "e-course-management/internal/user/repository"
	adminRepository "e-course-management/internal/admin/repository"
	userUseCase "e-course-management/internal/user/usecase"
	adminUseCase "e-course-management/internal/admin/usecase"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializedService(db *gorm.DB) *handler.OauthHandler {
	wire.Build(
		handler.NewOauthHandler,
		oauthUseCase.NewOauthUseCase,
		oauthRepository.NewOauthAccessTokenRepository,
		oauthRepository.NewOauthClientRepository,
		oauthRepository.NewOauthRefreshTokenRepository,
		userUseCase.NewUserUseCase,
		userRepository.NewUserRepository,
		adminRepository.NewAdminRepository,
		adminUseCase.NewAdminUseCase,
	)

	return &handler.OauthHandler{}
}
