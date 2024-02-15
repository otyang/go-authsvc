package custom

import (
	"github.com/otyang/go-authsvc/dto"
)

type (
	ForgotPasswordParams struct {
		Email                 string
		CodeExpirationMinutes int32
	}

	ResetPasswordParams struct {
		// methodID from forgot password stage
		MethodID     string
		EmailOTPCode string
		Password     string
	}

	SignupStartParams struct {
		Email                 string
		CodeExpirationMinutes int32
	}

	SignupCompleteParams struct {
		ReferenceID            string
		Password               string
		EmailOTPCode           string
		SessionDurationMinutes int32
	}

	SigninParams struct {
		Email                  string
		Password               string
		SessionDurationMinutes int32
		SessionClaims          dto.SessionClaims
	}

	SigninResponse struct {
		RequestID    string
		SessionID    string
		SessionToken string
		SessionJWT   string
		User         dto.User
	}
)
