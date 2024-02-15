package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserParams_UpdateWith(t *testing.T) {
	var (
		pin            = "0000"
		hashedPin, err = hashPin(pin)
	)
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		inputParams  UpdateUserParams // Input to UpdateWith
		initialMeta  TrustedMetadata  // Initial state of TrustedMetadata
		expectedMeta TrustedMetadata  // Expected state after update
	}{
		{
			name: "update user role",
			inputParams: UpdateUserParams{
				Name:               Name{},
				UserRole:           toPointer("admin"),
				UserPhoneNumber:    toPointer("1234567890"),
				UserProfileImage:   toPointer("https://image.com/logo.jpg"),
				NotificationEmail:  toPointer(true),
				NotificationPush:   toPointer(true),
				NotificationSMS:    toPointer(true),
				NotificationInApp:  toPointer(true),
				TransactionPinHash: toPointer(pin),
				WebhookURL:         toPointer("https://webhook.com/attend"),
			},
			initialMeta: DefaultTrustedMetadata,
			expectedMeta: TrustedMetadata{
				SystemUserID:       "",
				UserRole:           "admin",
				UserPhoneNumber:    toPointer("1234567890"),
				UserProfileImage:   toPointer("https://image.com/logo.jpg"),
				NotificationEmail:  true,
				NotificationPush:   true,
				NotificationSMS:    true,
				NotificationInApp:  true,
				TransactionPinHash: &hashedPin,
				WebhookURL:         toPointer("https://webhook.com/attend"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedMeta, err := tc.inputParams.UpdateWith(tc.initialMeta)

			assert.NoError(t, err)
			assert.Equal(t, updatedMeta.UserRole, tc.expectedMeta.UserRole)
			assert.Equal(t, updatedMeta.UserPhoneNumber, tc.expectedMeta.UserPhoneNumber)
			assert.Equal(t, updatedMeta.UserProfileImage, tc.expectedMeta.UserProfileImage)
			assert.Equal(t, updatedMeta.NotificationEmail, tc.expectedMeta.NotificationEmail)
			assert.Equal(t, updatedMeta.NotificationPush, tc.expectedMeta.NotificationPush)
			assert.Equal(t, updatedMeta.NotificationSMS, tc.expectedMeta.NotificationSMS)
			assert.Equal(t, updatedMeta.NotificationInApp, tc.expectedMeta.NotificationInApp)
			assert.Equal(t, updatedMeta.WebhookURL, tc.expectedMeta.WebhookURL)
			assert.Equal(t, hashedPin, *tc.expectedMeta.TransactionPinHash)
		})
	}
}

// Helper for creating string pointers (optional)
func toPointer[T any](s T) *T {
	return &s
}
