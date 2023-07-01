package forgot_password

import (
	dto "e-course-management/internal/forgot_password/dto"
	entity "e-course-management/internal/forgot_password/entity"
	repository "e-course-management/internal/forgot_password/repository"
	userDto "e-course-management/internal/user/dto"
	userUseCase "e-course-management/internal/user/usecase"
	mail "e-course-management/pkg/mail/sendgrid"
	"e-course-management/pkg/response"
	"e-course-management/pkg/utils"
	"errors"
	"time"
)

type ForgotPasswordUseCase interface {
	Create(dtoForgotPassword dto.ForgotPasswordRequestBody) (*entity.ForgotPassword, *response.Error)
	Update(dto dto.ForgotPasswordUpdateRequestBody) (*entity.ForgotPassword, *response.Error)
}

type forgotPasswordUseCase struct {
	repository  repository.ForgotPasswordRepository
	userUseCase userUseCase.UserUseCase
	mail        mail.Mail
}

// Create implements ForgotPasswordUseCase.
func (usecase *forgotPasswordUseCase) Create(dtoForgotPassword dto.ForgotPasswordRequestBody) (*entity.ForgotPassword, *response.Error) {
	// Check email
	user, err := usecase.userUseCase.FindByEmail(dtoForgotPassword.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &response.Error{
			Code: 200,
			Err:  errors.New("success, please check your email"),
		}
	}

	dateTime := time.Now().Add(24 * 1 * time.Hour)

	forgotPassword := entity.ForgotPassword{
		UserID:    &user.ID,
		Valid:     true,
		Code:      utils.RandString(32),
		ExpiredAt: &dateTime,
	}

	dataForgotPassword, err := usecase.repository.Create(forgotPassword)

	// Send email
	dataEmailForgotPassword := dto.ForgotPasswordEmailRequestBody{
		SUBJECT: "Code Forgot Password",
		EMAIL:   user.Email,
		CODE:    forgotPassword.Code,
	}

	go usecase.mail.SendForgotPassword(user.Email, dataEmailForgotPassword)

	if err != nil {
		return nil, err
	}

	return dataForgotPassword, nil
}

// Update implements ForgotPasswordUseCase.
func (usecase *forgotPasswordUseCase) Update(dto dto.ForgotPasswordUpdateRequestBody) (*entity.ForgotPassword, *response.Error) {
	// Check code
	code, err := usecase.repository.FindOneByCode(dto.Code)

	if err != nil || !code.Valid {
		return nil, &response.Error{
			Code: 400,
			Err:  errors.New("code is invalid"),
		}
	}

	// Search user
	user, err := usecase.userUseCase.FindOneById(int(*code.UserID))

	if err != nil {
		return nil, err
	}

	dataUser := userDto.UserUpdateRequestBody{
		Password: &dto.Password,
	}

	_, err = usecase.userUseCase.Update(int(user.ID), dataUser)

	if err != nil {
		return nil, err
	}

	code.Valid = false

	usecase.repository.Update(*code)

	return code, nil
}

func NewForgotPasswordUseCase(
	repository repository.ForgotPasswordRepository,
	userUseCase userUseCase.UserUseCase,
	mail mail.Mail,
) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		repository, userUseCase, mail,
	}
}
