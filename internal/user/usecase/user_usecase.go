package user

import (
	dto "e-course-management/internal/user/dto"
	entity "e-course-management/internal/user/entity"
	repository "e-course-management/internal/user/repository"
	"e-course-management/pkg/response"
	"e-course-management/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	FindAll(offset int, limit int) []entity.User
	FindByEmail(email string) (*entity.User, *response.Error)
	FindOneById(id int) (*entity.User, *response.Error)
	Create(dto dto.UserRequestBody) (*entity.User, *response.Error)
	FindOneByCodeVerified(codeVerified string) (*entity.User, *response.Error)
	Update(id int, dto dto.UserUpdateRequestBody) (*entity.User, *response.Error)
	Delete(id int) *response.Error
	TotalCountUser() int64
}

type userUseCase struct {
	repository repository.UserRepository
}

// Create implements UserUseCase.
func (usecase *userUseCase) Create(dto dto.UserRequestBody) (*entity.User, *response.Error) {
	checkUser, err := usecase.repository.FindByEmail(dto.Email)

	if err != nil && !errors.Is(err.Err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if checkUser != nil {
		return nil, &response.Error{
			Code: 409,
			Err:  errors.New("email sudah terdaftar"),
		}
	}

	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword(
		[]byte(dto.Password),
		bcrypt.DefaultCost,
	)

	if errHashedPassword != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  errHashedPassword,
		}
	}

	user := entity.User{
		Name:         dto.Name,
		Email:        dto.Email,
		Password:     string(hashedPassword),
		CodeVerified: utils.RandString(32),
	}

	dataUser, err := usecase.repository.Create(user)

	if err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  errHashedPassword,
		}
	}

	return dataUser, nil
}

// Delete implements UserUseCase.
func (*userUseCase) Delete(id int) *response.Error {
	panic("unimplemented")
}

// FindAll implements UserUseCase.
func (*userUseCase) FindAll(offset int, limit int) []entity.User {
	panic("unimplemented")
}

// FindByEmail implements UserUseCase.
func (usecase *userUseCase) FindByEmail(email string) (*entity.User, *response.Error) {
	return usecase.repository.FindByEmail(email)
}

// FindOneByCodeVerified implements UserUseCase.
func (*userUseCase) FindOneByCodeVerified(codeVerified string) (*entity.User, *response.Error) {
	panic("unimplemented")
}

// FindOneById implements UserUseCase.
func (*userUseCase) FindOneById(id int) (*entity.User, *response.Error) {
	panic("unimplemented")
}

// TotalCountUser implements UserUseCase.
func (*userUseCase) TotalCountUser() int64 {
	panic("unimplemented")
}

// Update implements UserUseCase.
func (usecase *userUseCase) Update(id int, dto dto.UserUpdateRequestBody) (*entity.User, *response.Error) {
	// Cari user berdasarkan id
	user, err := usecase.repository.FindOneById(id)

	if err != nil {
		return nil, err
	}

	if dto.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)

		if err != nil {
			return nil, &response.Error{
				Code: 500,
				Err:  err,
			}
		}

		user.Password = string(hashedPassword)
	}

	updateUser, err := usecase.repository.Update(*user)

	if err != nil {
		return nil, err
	}

	return updateUser, nil
}

func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}
