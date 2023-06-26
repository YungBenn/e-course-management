package register

import (
	userDto "e-course-management/internal/user/dto"
	userUseCase "e-course-management/internal/user/usecase"
	"e-course-management/pkg/response"
)

type RegisterUseCase interface {
	Register(dto userDto.UserRequestBody) *response.Error
}

type registerUseCase struct {
	userUseCase userUseCase.UserUseCase
}

// Register implements RegisterUseCase.
func (usecase *registerUseCase) Register(dto userDto.UserRequestBody) *response.Error {
	_, err := usecase.userUseCase.Create(dto)
	if err != nil {
		return err
	}

	return nil
}

func NewRegisterUseCase(userUseCase userUseCase.UserUseCase) RegisterUseCase {
	return &registerUseCase{userUseCase}
}
