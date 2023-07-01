package register

type EmailVerification struct {
	SUBJECT           string
	EMAIL             string
	VERIFICATION_CODE string
}