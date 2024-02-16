package custom

import (
	"context"

	"github.com/otyang/go-authsvc/dto"
	session_svc "github.com/otyang/go-authsvc/session"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/otp"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/otp/email"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords/session"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/users"
	"golang.org/x/sync/errgroup"
)

// CustomService encapsulates interactions with Stytch sessions API
type CustomService struct {
	client     *stytchapi.API
	sessionSvc *session_svc.SessionService
}

func NewCustomService(client *stytchapi.API) *CustomService {
	return &CustomService{
		client:     client,
		sessionSvc: session_svc.NewSessionService(client),
	}
}

func (s *CustomService) SignIn(ctx context.Context, param SigninParams) (*SigninResponse, error) {
	if param.SessionDurationMinutes == 0 {
		param.SessionDurationMinutes = dto.DefaultSessionDurationMinutes
	}

	sclaims, err := dto.DecodeFromXToX[map[string]any](param.SessionClaims, false)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Passwords.Authenticate(ctx, &passwords.AuthenticateParams{
		Email:                  param.Email,
		Password:               param.Password,
		SessionDurationMinutes: param.SessionDurationMinutes,
		SessionCustomClaims:    *sclaims,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	// Lets clear all sessions
	{
		listOfSessions, err := s.sessionSvc.List(ctx, resp.UserID, resp.Session.SessionID)
		if err != nil {
			return nil, dto.HandleError(err)
		}

		if len(listOfSessions) > 0 {
			var g errgroup.Group

			for _, sesn := range listOfSessions {
				if sesn.CurrentSession {
					continue
				}

				g.Go(func(sc *CustomService) func() error {
					return func() error {
						_, err := sc.client.Sessions.Revoke(ctx, &sessions.RevokeParams{SessionID: sesn.SessionID})
						return err
					}
				}(s))

				if err := g.Wait(); err != nil {
					return nil, dto.HandleError(err)
				}
			}
		}
	}

	return &SigninResponse{
		RequestID:    resp.RequestID,
		SessionID:    resp.Session.SessionID,
		SessionToken: resp.SessionToken,
		SessionJWT:   resp.SessionJWT,
		User:         dto.ConvertStytchUserToUser(resp.User),
	}, err
}

func (s *CustomService) SignupStart(ctx context.Context, p SignupStartParams) (string, error) {
	resp, err := s.client.OTPs.Email.LoginOrCreate(ctx, &email.LoginOrCreateParams{
		Email:               p.Email,
		ExpirationMinutes:   p.CodeExpirationMinutes,
		CreateUserAsPending: true,
	})
	if err != nil {
		return "", dto.HandleError(err)
	}

	return resp.EmailID, nil
}

func (s *CustomService) SignupComplete(ctx context.Context, param SignupCompleteParams) (*dto.User, error) {
	resp, err := s.client.OTPs.Authenticate(ctx, &otp.AuthenticateParams{
		MethodID:               param.ReferenceID,
		Code:                   param.EmailOTPCode,
		SessionDurationMinutes: param.SessionDurationMinutes,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	tm, err := dto.DecodeFromXToX[map[string]any](dto.DefaultTrustedMetadata, true)
	if err != nil {
		return nil, dto.HandleError(err)
	}

	// update user profile
	rsp, err := s.client.Users.Update(ctx, &users.UpdateParams{
		UserID:          resp.UserID,
		Name:            nil,
		TrustedMetadata: *tm,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	// This sets the password for the user account
	if _, err := s.client.Passwords.Sessions.Reset(ctx, &session.ResetParams{
		Password:     param.Password,
		SessionToken: resp.SessionToken,
	}); err != nil {
		return nil, dto.HandleError(err)
	}

	userResponse := dto.ConvertStytchUserToUser(rsp.User)
	return &userResponse, nil
}
