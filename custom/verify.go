package custom

import (
	"context"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/otp/email"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/otp"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords/existingpassword"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords/session"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/totps"
	"golang.org/x/crypto/bcrypt"
)

// Compares a plaintext pin with a stored pin hash using the bcrypt algorithm.
// It returns true if the passwords match, false otherwise.
func (s *CustomService) VerifyTransactionPin(pinHash string, pin string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pinHash), []byte(pin)) == nil
}

func (s *CustomService) VerifyTOTP(ctx context.Context, userID string, totpCode string) error {
	_, err := s.client.TOTPs.Authenticate(ctx, &totps.AuthenticateParams{
		UserID:   userID,
		TOTPCode: totpCode,
	})
	return dto.HandleError(err)
}

func (s *CustomService) VerifyPassword(ctx context.Context, email string, password string) error {
	_, err := s.client.Passwords.Authenticate(ctx, &passwords.AuthenticateParams{
		Email:    email,
		Password: password,
	})
	return dto.HandleError(err)
}

func (s *CustomService) ForgotPassword(ctx context.Context, param ForgotPasswordParams) (string, error) {
	resp, err := s.client.OTPs.Email.Send(ctx, &email.SendParams{
		Email:             param.Email,
		ExpirationMinutes: param.CodeExpirationMinutes,
	})
	if err != nil {
		return "", dto.HandleError(err)
	}

	return resp.EmailID, nil
}

func (s *CustomService) ResetPassword(ctx context.Context, param ResetPasswordParams) error {
	resp, err := s.client.OTPs.Authenticate(ctx, &otp.AuthenticateParams{
		MethodID:               param.MethodID,
		Code:                   param.EmailOTPCode,
		SessionDurationMinutes: 5,
	})
	if err != nil {
		return dto.HandleError(err)
	}

	// lets reset the password
	_, err = s.client.Passwords.Sessions.Reset(ctx, &session.ResetParams{
		Password:     param.Password,
		SessionToken: resp.SessionToken,
	})

	return dto.HandleError(err)
}

func (s *CustomService) UpdatePassword(ctx context.Context, email, existingPassword, newPassword string) error {
	_, err := s.client.Passwords.ExistingPassword.Reset(ctx, &existingpassword.ResetParams{
		Email:            email,
		ExistingPassword: existingPassword,
		NewPassword:      newPassword,
	})
	return dto.HandleError(err)
}

func (s *CustomService) ChangeEmailStartSendCode(ctx context.Context, sessionToken string, newEmail string) (string, error) {
	r, err := s.client.OTPs.Email.Send(ctx, &email.SendParams{
		Email:        newEmail,
		SessionToken: sessionToken,
	})
	if err != nil {
		return "", dto.HandleError(err)
	}

	return r.EmailID, nil
}

func (s *CustomService) ChangeEmailCompleteVerifyCode(ctx context.Context, sentOtpMethodID, code string) error {
	_, err := s.client.OTPs.Authenticate(ctx, &otp.AuthenticateParams{
		MethodID: sentOtpMethodID,
		Code:     code,
	})
	return dto.HandleError(err)
}
