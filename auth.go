package auth

import (
	"log"

	"github.com/otyang/go-authsvc/custom"
	"github.com/otyang/go-authsvc/session"
	"github.com/otyang/go-authsvc/user"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
)

type Auth struct {
	Custom       *custom.CustomService
	User         *user.UserService
	Session      *session.SessionService
	StytchClient *stytchapi.API
}

func New(stytchProjectID string, stytchSecret string) (*Auth, error) {
	client, err := stytchapi.NewClient(stytchProjectID, stytchSecret)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}

	return &Auth{
		Custom:       custom.NewCustomService(client),
		User:         user.NewUserService(client),
		Session:      session.NewSessionService(client),
		StytchClient: client,
	}, nil
}
