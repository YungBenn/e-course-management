package main

import (
	"github.com/gin-gonic/gin"
	userRepository "e-course-management/internal/user/repository"
	userUsecase "e-course-management/internal/user/usecase"
	mysql "e-course-management/pkg/db/mysql"

	registerHandler "e-course-management/internal/register/delivery/http"
	registerUseCase "e-course-management/internal/register/usecase"

	oauthHandler "e-course-management/internal/oauth/delivery/http"
	oauthRepository "e-course-management/internal/oauth/repository"
	oauthUsecase "e-course-management/internal/oauth/usecase"
)

func main() {
	r := gin.Default()
	db := mysql.DB()

	userRepository := userRepository.NewUserRepository(db)
	userUsecase := userUsecase.NewUserUseCase(userRepository)

	registerUsecase := registerUseCase.NewRegisterUseCase(userUsecase)
	registerHandler.NewRegisterHandler(registerUsecase).Route(&r.RouterGroup)

	oauthClientRepository := oauthRepository.NewOauthClientRepository(db)
	oauthAccessTokenRepository := oauthRepository.NewOauthAccessTokenRepository(db)
	oauthRefreshTokenRepository := oauthRepository.NewOauthRefreshTokenRepository(db)

	oauthUsecase := oauthUsecase.NewOauthUseCase(
		oauthClientRepository,
		oauthAccessTokenRepository,
		oauthRefreshTokenRepository,
		userUsecase,
	)

	oauthHandler.NewOauthHandler(oauthUsecase).Route(&r.RouterGroup)

	r.Run()
}