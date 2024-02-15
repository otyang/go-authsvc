package session

import (
	"errors"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
)

// Define custom types for authentication and logout parameters
type SessionAuthenticateParams struct {
	SessionToken              string
	OrSessionJWT              string
	ExtendSessionTTLByMinutes int32 // Minutes to extend session TTL (optional)
}

type SessionLogoutParams struct {
	SessionID      string // ID of the session to logout
	OrSessionToken string // Token of the session to logout (alternative to SessionID)
	OrSessionJWT   string // JWT of the session to logout (alternative to SessionID or SessionToken)
}

var (
	ErrPhoneNumberRequired = errors.New("phone number required")
	ErrTwoFARequired       = errors.New("two factor auth required")
)

func isTwoFARequiredForThisSession(session dto.Session, stytchAuthFactors []sessions.AuthenticationFactor) bool {
	// user have no two_fa activated or user initiated
	// two_fa usage on his account but havent completed initiation
	if !session.User.TotpIsEnabled {
		return false
	}

	if len(stytchAuthFactors) == 0 {
		return true
	}

	for _, stytchAuthFactor := range stytchAuthFactors {
		if stytchAuthFactor.DeliveryMethod == sessions.AuthenticationFactorDeliveryMethodAuthenticatorApp {
			return false
		}
	}

	return true
}

func isPhoneNumberRequiredForThisSession(session dto.Session) bool {
	return !session.User.PhoneIsVerified
}
