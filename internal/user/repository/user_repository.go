package user

import (
	entity "e-course-management/internal/user/entity"
	"e-course-management/pkg/response"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(offset int, limit int) []entity.User
	FindOneById(id int) (*entity.User, *response.Error)
	FindByEmail(email string) (*entity.User, *response.Error)
	Create(entity entity.User) (*entity.User, *response.Error)
	FindOneByCodeVerified(codeVerified string) (*entity.User, *response.Error)
	Update(entity entity.User) (*entity.User, *response.Error)
	Delete(entity entity.User) (*entity.User, *response.Error)
	TotalCountUser() int64
}

type userRepository struct {
	db *gorm.DB
}

// Create implements UserRepository.
func (repository *userRepository) Create(entity entity.User) (*entity.User, *response.Error) {
	if err := repository.db.Create(&entity).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &entity, nil
}

// Delete implements UserRepository.
func (*userRepository) Delete(entity entity.User) (*entity.User, *response.Error) {
	panic("unimplemented")
}

// FindAll implements UserRepository.
func (*userRepository) FindAll(offset int, limit int) []entity.User {
	panic("unimplemented")
}

// FindByEmail implements UserRepository.
func (repository *userRepository) FindByEmail(email string) (*entity.User, *response.Error) {
	var user entity.User

	if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &user, nil
}

// FindOneByCodeVerified implements UserRepository.
func (*userRepository) FindOneByCodeVerified(codeVerified string) (*entity.User, *response.Error) {
	panic("unimplemented")
}

// FindOneById implements UserRepository.
func (*userRepository) FindOneById(id int) (*entity.User, *response.Error) {
	panic("unimplemented")
}

// TotalCountUser implements UserRepository.
func (*userRepository) TotalCountUser() int64 {
	panic("unimplemented")
}

// Update implements UserRepository.
func (*userRepository) Update(entity entity.User) (*entity.User, *response.Error) {
	panic("unimplemented")
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}
