package oauth

import (
	dto "e-course-management/internal/oauth/dto"
	entity "e-course-management/internal/oauth/entity"
	repository "e-course-management/internal/oauth/repository"
	userUseCase "e-course-management/internal/user/usecase"
	adminUseCase "e-course-management/internal/admin/usecase"
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
	Refresh(dtoRefreshToken dto.RefreshTokenRequestBody) (*dto.LoginResponse, *response.Error)
}

type oauthUseCase struct {
	oauthClientRepository       repository.OauthClientRepository
	oauthAccessTokenRepository  repository.OauthAccessTokenRepository
	oauthRefreshTokenRepository repository.OauthRefreshTokenRepository
	userUseCase                 userUseCase.UserUseCase
	adminUseCase				adminUseCase.AdminUseCase
}

// Refresh implements OauthUseCase.
func (usecase *oauthUseCase) Refresh(dtoRefreshToken dto.RefreshTokenRequestBody) (*dto.LoginResponse, *response.Error) {
	oauthRefreshToken, err := usecase.oauthRefreshTokenRepository.FindOneByToken(dtoRefreshToken.RefreshToken)
	
	if err != nil {
		return nil, err
	}

	if oauthRefreshToken.ExpiredAt.Before(time.Now()) {
		return nil, &response.Error{
			Code: 500,
			Err:  errors.New("your refresh token is already expired"),
		}
	}

	var user dto.UserResponse

	expirationTime := time.Now().Add(24 * 365 * time.Hour) // 1 year

	if *oauthRefreshToken.OauthAccessToken.OauthClientID == 2 {
		admin, _ := usecase.adminUseCase.FindOneById(int(oauthRefreshToken.UserID))

		user.ID = admin.ID
		user.Name = admin.Name
		user.Email = admin.Email
	} else {
		dataUser, _ := usecase.userUseCase.FindOneById(int(oauthRefreshToken.UserID))

		user.ID = dataUser.ID
		user.Name = dataUser.Name
		user.Email = dataUser.Email
	}

	claims := &dto.ClaimsResponse{
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		IsAdmin:          false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	if *oauthRefreshToken.OauthAccessToken.OauthClientID == 2 {
		claims.IsAdmin = true
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errSignedString := token.SignedString(jwtKey)

	if errSignedString != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  errSignedString,
		}
	}

	dataOauthAccessToken := entity.OauthAccessToken{
		OauthClientID: oauthRefreshToken.OauthAccessToken.OauthClientID,
		UserID:        oauthRefreshToken.UserID,
		Token:         tokenString,
		Scope:         "*",
		ExpiredAt:     &expirationTime,
	}

	saveOauthAccessToken, err := usecase.oauthAccessTokenRepository.Create(dataOauthAccessToken)

	if err != nil {
		return nil, err
	}

	expirationTimeOauthRefreshToken := time.Now().Add(24 * 366 * time.Hour)

	dataOauthRefreshToken := entity.OauthRefreshToken{
		OauthAccessTokenID: &saveOauthAccessToken.ID,
		UserID:             oauthRefreshToken.UserID,
		Token:              utils.RandString(128),
		ExpiredAt:          &expirationTimeOauthRefreshToken,
	}

	saveOauthRefreshToken, err := usecase.oauthRefreshTokenRepository.Create(dataOauthRefreshToken)

	if err != nil {
		return nil, err
	}

	err = usecase.oauthRefreshTokenRepository.Delete(*oauthRefreshToken)

	if err != nil {
		return nil, err
	}
	
	err = usecase.oauthAccessTokenRepository.Delete(*oauthRefreshToken.OauthAccessToken)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokenString,
		RefreshToken: saveOauthRefreshToken.Token,
		Type:         "Bearer",
		ExpiredAt:    expirationTime.Format(time.RFC3339),
		Scope:        "*",
	}, nil
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

	if oauthClient.Name == "web-admin" {
		dataAdmin, err := usecase.adminUseCase.FindOneByEmail(dtoLoginRequestBody.Email)
	
		if err != nil {
			return nil, &response.Error{
				Code: 400,
				Err:  errors.New("username or password is invalid"),
			}
		}
	
		user.ID = dataAdmin.ID
		user.Email = dataAdmin.Email
		user.Name = dataAdmin.Name
		user.Password = dataAdmin.Password
	} else {
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
	}


	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	// Compare password
	errorBcrypt := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(dtoLoginRequestBody.Password),
	)

	if errorBcrypt != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  errors.New("username or password is invalid"),
		}
	}

	expirationTime := time.Now().Add(24 * 365 * time.Hour) // 1 year

	claims := &dto.ClaimsResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	if oauthClient.Name == "web-admin" {
		claims.IsAdmin = true
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
	adminUseCase adminUseCase.AdminUseCase,
) OauthUseCase {
	return &oauthUseCase{
		oauthClientRepository,
		oauthAccessTokenRepository,
		oauthRefreshTokenRepository,
		userUseCase,
		adminUseCase,
	}
}
