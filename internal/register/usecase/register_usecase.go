package register

import (
	registerDto "e-course-management/internal/register/dto"
	userDto "e-course-management/internal/user/dto"
	userUseCase "e-course-management/internal/user/usecase"
	mail "e-course-management/pkg/mail/sendgrid"
	"e-course-management/pkg/response"
)

type RegisterUseCase interface {
	Register(dto userDto.UserRequestBody) *response.Error
}

type registerUseCase struct {
	userUseCase userUseCase.UserUseCase
	mail        mail.Mail
}

// Register implements RegisterUseCase.
func (usecase *registerUseCase) Register(dto userDto.UserRequestBody) *response.Error {
	user, err := usecase.userUseCase.Create(dto)
	if err != nil {
		return err
	}

	// Melakukan pengiriman melalui email dengan sendgrid
	data := registerDto.EmailVerification{
		SUBJECT:           "Verification Account",
		EMAIL:             dto.Email,
		VERIFICATION_CODE: user.CodeVerified,
	}

	go usecase.mail.SendVerification(dto.Email, data)

	return nil
}

func NewRegisterUseCase(userUseCase userUseCase.UserUseCase, mail mail.Mail) RegisterUseCase {
	return &registerUseCase{
		userUseCase: userUseCase,
		mail:        mail,
	}
}
