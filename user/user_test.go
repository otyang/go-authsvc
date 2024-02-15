package user

import (
	"context"
	"log"
	"testing"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

const (
	stytchProjectID = "project-test-098574c2-9deb-4dac-b8e7-bb3b1b3d2224"
	stytchSecret    = "secret-test-G1CewToEM0_XIrN4xLqDVVGlb0A2NXSV9oQ="
)

func setUpClient() *stytchapi.API {
	client, err := stytchapi.NewClient(stytchProjectID, stytchSecret)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}
	return client
}

func TestUserService_Get(t *testing.T) {
	t.Parallel()

	u := &UserService{
		client: setUpClient(),
	}

	got, err := u.Get(context.TODO(), "user-test-d832d3af-e34b-45c4-8c90-99d90681236d")

	assert.Error(t, err)
	assert.Empty(t, got)

	v, ok := err.(stytcherror.Error)
	assert.True(t, ok)
	assert.Equal(t, "user_not_found", string(v.ErrorType))
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Parallel()

	userId := "user-test-b5c555a2-98ca-40b1-b45b-ebf8d01"

	u := &UserService{
		client: setUpClient(),
	}

	got, err := u.UpdateProfile(context.TODO(), userId, &dto.UpdateUserParams{
		Name:               dto.Name{FirstName: "Xxa", MiddleName: "X", LastName: "X"},
		UserRole:           toPointer("1role"),
		UserPhoneNumber:    toPointer("1_phone_number"),
		UserProfileImage:   toPointer("1_profile_image"),
		NotificationEmail:  toPointer(true),
		NotificationPush:   toPointer(false),
		NotificationSMS:    toPointer(true),
		NotificationInApp:  toPointer(false),
		TransactionPinHash: toPointer("1pin"),
	})

	assert.Error(t, err)
	assert.Empty(t, got)

	v, ok := err.(stytcherror.Error)
	assert.True(t, ok)
	assert.Equal(t, "invalid_user_id", string(v.ErrorType))
}

func toPointer[T any](t T) *T {
	return &t
}
