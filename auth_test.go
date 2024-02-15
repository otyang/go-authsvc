package auth

import (
	"context"
	"fmt"

	"github.com/otyang/go-authsvc/custom"
	"github.com/otyang/go-authsvc/session"
)

func Example() {
	auth, err := New("projectID", "projectSecret")
	if err != nil {
		// handle err
	}

	signInResponse, err := auth.Custom.SignIn(context.TODO(), custom.SigninParams{})
	if err != nil {
		// handle err
	}
	fmt.Println("Sign-In Response:", signInResponse)

	// User Retrieval
	user, err := auth.User.Get(context.TODO(), "USER_ID")
	if err != nil {
		// handle err
	}
	fmt.Println("User Data:", user)

	// Session Authentication
	session, err := auth.Session.Authenticate(context.TODO(), session.SessionAuthenticateParams{})
	if err != nil {
		// handle err
	}
	fmt.Println("Session Details:", session)

	// etc...
}
