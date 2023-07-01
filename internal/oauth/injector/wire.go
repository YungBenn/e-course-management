//go:build wireinject
// +build wireinject

package oauth

import (
	handler "e-course-management/internal/oauth/delivery/http"
	oauthUseCase "e-course-management/internal/oauth/usecase"
	oauthRepository "e-course-management/internal/oauth/repository"
	userUseCase "e-course-management/internal/user/usecase"
	userRepository "e-course-management/internal/user/repository"

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
	)

	return &handler.OauthHandler{}
}