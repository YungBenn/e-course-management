package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	forgotPasswordDto "e-course-management/internal/forgot_password/dto"
	registerDto "e-course-management/internal/register/dto"
)

type Mail interface {
	SendVerification(toEmail string, data registerDto.EmailVerification)
	SendForgotPassword(toEmail string, data forgotPasswordDto.ForgotPasswordEmailRequestBody)
}

type mailUsecase struct {
}

// SendForgotPassword implements Mail
func (usecase *mailUsecase) SendForgotPassword(toEmail string, data forgotPasswordDto.ForgotPasswordEmailRequestBody) {
	cwd, _ := os.Getwd()
	templateFile := filepath.Join(cwd, "/templates/emails/forgot_password.html")

	result, err := ParseTemplate(templateFile, data)

	if err != nil {
		fmt.Println(err)
	} else {
		usecase.sendMail(toEmail, result, data.SUBJECT)
	}
}

// SendVerification implements Mail
func (usecase *mailUsecase) SendVerification(toEmail string, data registerDto.EmailVerification) {
	cwd, _ := os.Getwd()
	templateFile := filepath.Join(cwd, "/templates/emails/verification_email.html")

	result, err := ParseTemplate(templateFile, data)

	if err != nil {
		fmt.Println(err)
	} else {
		usecase.sendMail(toEmail, result, data.SUBJECT)
	}
}

func (usecase *mailUsecase) sendMail(toEmail string, result string, subject string) {
	from := mail.NewEmail(os.Getenv("MAIL_SENDER_NAME"), os.Getenv("MAIL_SENDER_NAME"))
	to := mail.NewEmail(toEmail, toEmail)

	message := mail.NewSingleEmail(from, subject, to, "", result)

	client := sendgrid.NewSendClient(os.Getenv("MAIL_KEY"))
	resp, err := client.Send(message)

	if err != nil {
		fmt.Println(err)
	} else if resp.StatusCode != 200 {
		fmt.Println(resp)
	} else {
		fmt.Printf("success send email to %s", toEmail)
	}
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func NewMailUseCase() Mail {
	return &mailUsecase{}
}
