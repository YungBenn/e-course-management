package oauth

import (
	entity "e-course-management/internal/oauth/entity"
	"e-course-management/pkg/response"

	"gorm.io/gorm"
)

type OauthAccessTokenRepository interface {
	Create(entity entity.OauthAccessToken) (*entity.OauthAccessToken, *response.Error)
	Delete(entity.OauthAccessToken) *response.Error
	FindOneByAccessToken(accessToken string) (*entity.OauthAccessToken, *response.Error)
}

type oauthAccessTokenRepository struct {
	db *gorm.DB
}

// Create implements OauthAccessTokenRepository.
func (repository *oauthAccessTokenRepository) Create(entity entity.OauthAccessToken) (*entity.OauthAccessToken, *response.Error) {
	if err := repository.db.Create(&entity).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &entity, nil
}

// Delete implements OauthAccessTokenRepository.
func (*oauthAccessTokenRepository) Delete(entity.OauthAccessToken) *response.Error {
	panic("unimplemented")
}

// FindOneByAccessToken implements OauthAccessTokenRepository.
func (*oauthAccessTokenRepository) FindOneByAccessToken(accessToken string) (*entity.OauthAccessToken, *response.Error) {
	panic("unimplemented")
}

func NewOauthAccessTokenRepository(db *gorm.DB) OauthAccessTokenRepository {
	return &oauthAccessTokenRepository{db}
}
