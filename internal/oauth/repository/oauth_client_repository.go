package oauth

import (
	entity "e-course-management/internal/oauth/entity"
	"e-course-management/pkg/response"

	"gorm.io/gorm"
)

type OauthClientRepository interface {
	FindByClientIDAndClientSecret(clientID string, clientSecret string) (*entity.OauthClient, *response.Error)
}

type oauthClientRepository struct {
	db *gorm.DB
}

// FindByClientIDAndClientSecret implements OauthClientRepository.
func (repository *oauthClientRepository) FindByClientIDAndClientSecret(clientID string, clientSecret string) (*entity.OauthClient, *response.Error) {
	var oauthClient entity.OauthClient

	if err := repository.db.Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&oauthClient).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}	
	}

	return &oauthClient, nil
}

func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{db}
}
