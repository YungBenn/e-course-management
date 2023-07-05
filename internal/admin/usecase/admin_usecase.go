package admin

import (
	dto "e-course-management/internal/admin/dto"
	entity "e-course-management/internal/admin/entity"
	repository "e-course-management/internal/admin/repository"
	"e-course-management/pkg/response"

	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase interface {
	FindAll(offset int, limit int) []entity.Admin
	FindOneById(id int) (*entity.Admin, *response.Error)
	FindOneByEmail(email string) (*entity.Admin, *response.Error)
	Create(dto dto.AdminRequestBody) (*entity.Admin, *response.Error)
	Update(id int, dto dto.AdminRequestBody) (*entity.Admin, *response.Error)
	Delete(id int) *response.Error
	TotalCountAdmin() int64
}

type adminUseCase struct {
	repository repository.AdminRepository
}

// Create implements AdminUseCase.
func (usecase *adminUseCase) Create(dto dto.AdminRequestBody) (*entity.Admin, *response.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, &response.Error{
			Code: 500,
			Err:  err,
		}
	}

	dataAdmin := entity.Admin{
		Email: dto.Email,
		Name: dto.Name,
		Password: string(hashedPassword),
	}

	admin, errCreateAdmin := usecase.repository.Create(dataAdmin)

	if errCreateAdmin != nil {
		return nil, errCreateAdmin
	}

	return admin, nil
}

// Delete implements AdminUseCase.
func (usecase *adminUseCase) Delete(id int) *response.Error {
	admin, err := usecase.repository.FindOneById(id)
	
	if err != nil {
		return err
	}

	if err := usecase.repository.Delete(*admin); err != nil {
		return err
	}

	return nil
}

// FindAll implements AdminUseCase.
func (usecase *adminUseCase) FindAll(offset int, limit int) []entity.Admin {
	return usecase.repository.FindAll(offset, limit)
}

// FindOneByEmail implements AdminUseCase.
func (usecase *adminUseCase) FindOneByEmail(email string) (*entity.Admin, *response.Error) {
	return usecase.repository.FindOneByEmail(email)
}

// FindOneById implements AdminUseCase.
func (usecase *adminUseCase) FindOneById(id int) (*entity.Admin, *response.Error) {
	return usecase.repository.FindOneById(id)
}

// TotalCountAdmin implements AdminUseCase.
func (usecase *adminUseCase) TotalCountAdmin() int64 {
	panic("unimplemented")
}

// Update implements AdminUseCase.
func (usecase *adminUseCase) Update(id int, dto dto.AdminRequestBody) (*entity.Admin, *response.Error) {
	admin, err := usecase.repository.FindOneById(id)

	if err != nil {
		return nil, err
	}

	admin.Name = dto.Name
	admin.Email = dto.Email

	if dto.Password != nil {
		hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)

		if errHashedPassword != nil {
			return nil, &response.Error{
				Code: 500,
				Err:  errHashedPassword,
			}
		}

		admin.Password = string(hashedPassword)
	}

	updateAdmin, err := usecase.repository.Update(*admin)

	if err != nil {
		return nil, err
	}

	return updateAdmin, nil
}

func NewAdminUseCase(repository repository.AdminRepository) AdminUseCase {
	return &adminUseCase{repository}
}
