package custom

import (
	"context"
	"log"
	"testing"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
)

const (
	stytchProjectID    = "project-test-6c64090b-e450-4e8f-82ff-4d21360fe41f"
	stytchSecret       = "secret-test-RYtim--PKcjaTRZmbd37nbHnvsWl3GkEwG8="
	stytchTestEmail    = "sandbox@stytch.com"
	stytchTestOTPCode  = "000000"
	stytchTestMethodID = "email-test-23873e89-d4ed-4e92-b3b9-e5c7198fa286"
)

const (
	testEmail    = "khh4of3h@spicy.homes"
	testPassword = "X~g:)6h]A7(?`Da}q8UkPx"
)

func setUpClient() *stytchapi.API {
	client, err := stytchapi.NewClient(stytchProjectID, stytchSecret)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}
	return client
}

func TestCustomService_Signin(t *testing.T) {
	t.Parallel()

	s := NewCustomService(setUpClient())
	s.client = setUpClient()

	// lets sign in a user
	rsp, err := s.SignIn(context.TODO(), SigninParams{
		Email:                  testEmail,
		Password:               testPassword,
		SessionDurationMinutes: 30,
		SessionClaims: dto.SessionClaims{
			DeviceIPAddress:  "127.0.0.2",
			DeviceUserAgent:  "brave-browser",
			DeviceType:       "desktop",
			IPAddressCity:    "moscow",
			IPAddressCountry: "russia",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, rsp)

	r1, err := s.client.Sessions.Get(context.TODO(), &sessions.GetParams{UserID: rsp.User.UserID})
	assert.NoError(t, err)

	// ensuring other sessions were logged out & only active session is on
	assert.Equal(t, len(r1.Sessions), 1)
	assert.Equal(t, rsp.SessionID, r1.Sessions[0].SessionID)
}

func TestCustomService_SignupStart(t *testing.T) {
	t.Parallel()

	s := &CustomService{
		client: setUpClient(),
	}

	refID, err := s.SignupStart(context.TODO(), SignupStartParams{
		Email:                 testEmail,
		CodeExpirationMinutes: 10,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, refID)
	assert.Equal(t, stytchTestMethodID, refID)
}

func TestCustomService_SignupComplete(t *testing.T) {
	t.Parallel()

	s := &CustomService{client: setUpClient()}

	_, err := s.SignupComplete(context.TODO(), SignupCompleteParams{
		ReferenceID:            "email-test-50d8a6c9-2b99-40c7-8553-0816927d71ef", // stytchTestMethodID,
		Password:               testPassword,
		EmailOTPCode:           "241627", // stytchTestOTPCode,
		SessionDurationMinutes: 10,
	})

	assert.Error(t, err)

	// v, ok := err.(stytcherror.Error)
	// assert.True(t, ok)
	// assert.Equal(t, "user_not_found", string(v.ErrorType))
}
