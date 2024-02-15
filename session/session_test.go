package session

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

const (
	stytchProjectID  = "project-test-6c64090b-e450-4e8f-82ff-4d21360fe41f"
	stytchSecret     = "secret-test-RYtim--PKcjaTRZmbd37nbHnvsWl3GkEwG8="
	testSessionToken = "WJtR5BCy38Szd5AfoDpf0iqFKEt4EE5JhjlWUY7l3FtY"
)

func setUpClient() *stytchapi.API {
	client, err := stytchapi.NewClient(stytchProjectID, stytchSecret)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}
	return client
}

func TestSessionService_Logout(t *testing.T) {
	t.Parallel()

	s := &SessionService{client: setUpClient()}

	err := s.Logout(context.TODO(), SessionLogoutParams{
		SessionID:      "testSessionUserID",
		OrSessionToken: "",
		OrSessionJWT:   "",
	})

	assert.NoError(t, err)
}

func TestSessionService_Authenticate(t *testing.T) {
	t.Parallel()

	s := &SessionService{client: setUpClient()}

	got, err := s.Authenticate(context.TODO(), SessionAuthenticateParams{
		SessionToken:              testSessionToken,
		OrSessionJWT:              "",
		ExtendSessionTTLByMinutes: 200,
	})

	assert.Error(t, err)
	assert.Empty(t, got)
	assert.Equal(t, err, ErrPhoneNumberRequired)
}

func TestSessionService_List(t *testing.T) {
	t.Parallel()

	var (
		userId    = "user-test-960b45dc-bc85-49ea-a255-88403af8bd10"
		client    = &SessionService{client: setUpClient()}
		list, err = client.List(context.TODO(), userId, "")
	)

	assert.Error(t, err)
	assert.Empty(t, list)

	v, ok := err.(stytcherror.Error)
	assert.True(t, ok)
	assert.Equal(t, "user_not_found", string(v.ErrorType))
}

func TestSessionService_List_Detailed(t *testing.T) {
	t.Parallel()

	var (
		sessionId = "session-test-e84b32da-2d6a-41c8-89d2-d3cb54440b41"
		client    = &SessionService{client: setUpClient()}
		list, err = client.List(context.TODO(), "user-test-960b45dc-bc85-49ea-a255-88403af8bd17", sessionId)
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	for _, s := range list {
		if s.SessionID == sessionId {
			assert.True(t, s.CurrentSession)
		} else {
			assert.False(t, s.CurrentSession)
		}
	}
}
