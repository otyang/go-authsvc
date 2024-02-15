package session

import (
	"context"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

// SessionService encapsulates interactions with Stytch sessions API
type SessionService struct {
	client *stytchapi.API
}

func NewSessionService(client *stytchapi.API) *SessionService {
	return &SessionService{
		client: client,
	}
}

// Logout logs out of a specific session using its ID, or token, or JWT.
func (s *SessionService) Logout(ctx context.Context, p SessionLogoutParams) error {
	_, err := s.client.Sessions.Revoke(ctx, &sessions.RevokeParams{
		SessionID:    p.SessionID,
		SessionToken: p.OrSessionToken,
		SessionJWT:   p.OrSessionJWT,
	})

	if v, ok := err.(stytcherror.Error); ok {
		if string(v.ErrorType) == "invalid_session_id" {
			return nil
		}
	}

	return dto.HandleError(err)
}

// Lists active sessions for a user.
func (s *SessionService) List(ctx context.Context, userID string, currentSessionId string) ([]dto.SessionListResponse, error) {
	resp, err := s.client.Sessions.Get(ctx, &sessions.GetParams{
		UserID: userID,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	var sessList []dto.SessionListResponse
	{
		for _, session := range resp.Sessions {
			sessionClaims, err := dto.DecodeFromXToX[dto.SessionClaims](session.CustomClaims, false)
			if err != nil {
				return nil, err
			}

			sessList = append(sessList, dto.SessionListResponse{
				SessionID:        session.SessionID,
				CurrentSession:   currentSessionId == session.SessionID,
				LastAccessedAt:   *session.LastAccessedAt,
				StartedAt:        *session.StartedAt,
				ExpiresAt:        *session.ExpiresAt,
				DeviceIPAddress:  sessionClaims.DeviceIPAddress,
				DeviceUserAgent:  sessionClaims.DeviceUserAgent,
				DeviceType:       sessionClaims.DeviceType,
				IPAddressCity:    sessionClaims.IPAddressCity,
				IPAddressCountry: sessionClaims.IPAddressCountry,
			})
		}
	}
	return sessList, nil
}

// Authenticates a session and returns user and session details.
func (s *SessionService) Authenticate(ctx context.Context, p SessionAuthenticateParams) (*dto.Session, error) {
	resp, err := s.client.Sessions.Authenticate(ctx, &sessions.AuthenticateParams{
		SessionToken:           p.SessionToken,
		SessionJWT:             p.OrSessionJWT,
		SessionDurationMinutes: p.ExtendSessionTTLByMinutes,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	// Convert Stytch user and session data to internal DTOs
	user := dto.ConvertStytchUserToUser(resp.User)
	sessionClaims, err := dto.DecodeFromXToX[dto.SessionClaims](resp.Session.CustomClaims, false)
	if err != nil {
		return nil, err
	}

	sn := dto.Session{
		UserID:           user.UserID,
		ID:               resp.Session.SessionID,
		Jwt:              resp.SessionJWT,
		Token:            resp.SessionToken,
		StartedAt:        *resp.Session.StartedAt,
		LastAccessedAt:   *resp.Session.LastAccessedAt,
		ExpiresAt:        *resp.Session.ExpiresAt,
		DeviceIPAddress:  sessionClaims.DeviceIPAddress,
		DeviceUserAgent:  sessionClaims.DeviceUserAgent,
		DeviceType:       sessionClaims.DeviceType,
		IPAddressCity:    sessionClaims.IPAddressCity,
		IPAddressCountry: sessionClaims.IPAddressCountry,
		User:             user,
	}

	if isTwoFARequiredForThisSession(sn, resp.Session.AuthenticationFactors) {
		return nil, ErrTwoFARequired
	}

	if isPhoneNumberRequiredForThisSession(sn) {
		return nil, ErrPhoneNumberRequired
	}

	return &sn, nil
}
