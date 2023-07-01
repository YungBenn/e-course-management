package forgot_password

import (
	entity "e-course-management/internal/forgot_password/entity"
	"e-course-management/pkg/response"

	"gorm.io/gorm"
)

type ForgotPasswordRepository interface {
	Create(entity entity.ForgotPassword) (*entity.ForgotPassword, *response.Error)
	FindOneByCode(code string) (*entity.ForgotPassword, *response.Error)
	Update(entity entity.ForgotPassword) (*entity.ForgotPassword, *response.Error)
}

type forgotPasswordRepository struct {
	db *gorm.DB
}

// Create implements ForgotPasswordRepository.
func (repository *forgotPasswordRepository) Create(entity entity.ForgotPassword) (*entity.ForgotPassword, *response.Error) {
	if err := repository.db.Create(&entity).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &entity, nil
}

// FindOneByCode implements ForgotPasswordRepository.
func (repository *forgotPasswordRepository) FindOneByCode(code string) (*entity.ForgotPassword, *response.Error) {
	var forgotPassword entity.ForgotPassword

	if err := repository.db.Where("code = ?", code).First(&forgotPassword).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &forgotPassword, nil
}

// Update implements ForgotPasswordRepository.
func (repository *forgotPasswordRepository) Update(entity entity.ForgotPassword) (*entity.ForgotPassword, *response.Error) {
	if err := repository.db.Save(&entity).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &entity, nil
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &forgotPasswordRepository{db}
}
