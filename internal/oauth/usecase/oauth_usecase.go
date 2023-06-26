package oauth

import (
	dto "e-course-management/internal/oauth/dto"
	entity "e-course-management/internal/oauth/entity"
	repository "e-course-management/internal/oauth/repository"
	userUseCase "e-course-management/internal/user/usecase"
	"e-course-management/pkg/response"
	"e-course-management/pkg/utils"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type OauthUseCase interface {
	Login(dtoLoginRequestBody dto.LoginRequestBody) (*dto.LoginResponse, *response.Error)
}

type oauthUseCase struct {
	oauthClientRepository       repository.OauthClientRepository
	oauthAccessTokenRepository  repository.OauthAccessTokenRepository
	oauthRefreshTokenRepository repository.OauthRefreshTokenRepository
	userUseCase                 userUseCase.UserUseCase
}

// Login implements OauthUseCase.
func (usecase *oauthUseCase) Login(dtoLoginRequestBody dto.LoginRequestBody) (*dto.LoginResponse, *response.Error) {
	oauthClient, err := usecase.oauthClientRepository.FindByClientIDAndClientSecret(
		dtoLoginRequestBody.ClientID,
		dtoLoginRequestBody.ClientSecret,
	)

	if err != nil {
		return nil, err
	}

	var user dto.UserResponse

	dataUser, err := usecase.userUseCase.FindByEmail(dtoLoginRequestBody.Email)

	if err != nil {
		return nil, &response.Error{
			Code: 400,
			Err:  errors.New("username or password is invalid"),
		}
	}

	user.ID = dataUser.ID
	user.Email = dataUser.Email
	user.Name = dataUser.Name
	user.Password = dataUser.Password

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	// Compare password
	errorBcrypt := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(dtoLoginRequestBody.Password),
	)

	if errorBcrypt != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  errorBcrypt,
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour) // 1 day

	claims := &dto.ClaimsResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	// Insert data to oauth access token table
	dataOauthAccessToken := entity.OauthAccessToken{
		OauthClientID: &oauthClient.ID,
		UserID:        user.ID,
		Token:         tokenString,
		Scope:         "*",
		ExpiredAt:     &expirationTime,
	}

	oauthAccessToken, err := usecase.oauthAccessTokenRepository.Create(dataOauthAccessToken)

	if err != nil {
		return nil, err
	}

	expirationTimeOauthAccessToken := time.Now().Add(24 * 366 * time.Hour)

	// Insert data to oauth refresh token table
	dataOauthRefreshToken := entity.OauthRefreshToken{
		OauthAccessTokenID: &oauthAccessToken.ID,
		UserID:             user.ID,
		Token:              utils.RandString(128),
		ExpiredAt:          &expirationTimeOauthAccessToken,
	}

	oauthRefreshToken, err := usecase.oauthRefreshTokenRepository.Create(dataOauthRefreshToken)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  oauthAccessToken.Token,
		RefreshToken: oauthRefreshToken.Token,
		Type:         "Bearer",
		ExpiredAt:    expirationTime.Format(time.RFC3339),
		Scope:        "*",
	}, nil
}

func NewOauthUseCase(
	oauthClientRepository repository.OauthClientRepository,
	oauthAccessTokenRepository repository.OauthAccessTokenRepository,
	oauthRefreshTokenRepository repository.OauthRefreshTokenRepository,
	userUseCase userUseCase.UserUseCase,
) OauthUseCase {
	return &oauthUseCase{
		oauthClientRepository,
		oauthAccessTokenRepository,
		oauthRefreshTokenRepository,
		userUseCase,
	}
}
