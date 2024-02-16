package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/users"
)

var testMetadata, testMetadataMap = func() (TrustedMetadata, map[string]any) {
	testTrustedMetadata := TrustedMetadata{
		SystemUserID:      "12345TRE",
		UserRole:          "customer",
		UserPhoneNumber:   toPointer("+1234567890"),
		UserProfileImage:  nil,
		NotificationEmail: false,
		NotificationPush:  false,
		NotificationSMS:   false,
		NotificationInApp: false,
		PinHash:           nil,
		WebhookURL:        nil,
	}

	m, _ := DecodeFromXToX[map[string]any](testTrustedMetadata, false)
	return testTrustedMetadata, *m
}()

var testUser = users.User{
	UserID: "userid-1234567890",
	Emails: []users.Email{
		{EmailID: "test-email-id-1", Email: "test@example.com", Verified: false},
		{EmailID: "test-email-id-2", Email: "test2@example.com", Verified: true},
	},
	Status: "active",
	Providers: []users.OAuthProvider{
		{
			ProviderType:            "Github",
			ProviderSubject:         "subject",
			ProfilePictureURL:       "https://google.com/picture.jpg",
			Locale:                  "en",
			OAuthUserRegistrationID: "oauth-registration-id",
		},
	},
	TOTPs: []users.TOTP{
		{
			TOTPID:   "totp-id-1",
			Verified: true,
		},
	},
	Name: &users.Name{
		FirstName:  "first",
		MiddleName: "middle",
		LastName:   "last",
	},
	CreatedAt: &time.Time{},
	Password: &users.Password{
		PasswordID:    "password-id",
		RequiresReset: false,
	},
	TrustedMetadata: testMetadataMap,
}

func TestConvertStytchUserToUser(t *testing.T) {
	want := User{
		UserID:            testUser.UserID,
		SytstemUserID:     testMetadata.SystemUserID,
		FirstName:         toPointer(testUser.Name.FirstName),
		MiddleName:        toPointer(testUser.Name.MiddleName),
		LastName:          toPointer(testUser.Name.LastName),
		FullName:          "first middle last",
		Email:             toPointer(testUser.Emails[1].Email),
		PhoneNumber:       toPointer("+1234567890"),
		TotpId:            toPointer(testUser.TOTPs[0].TOTPID),
		EmailIsVerified:   true,
		PhoneIsVerified:   true,
		PasswordIsEnabled: true,
		TotpIsEnabled:     true,
		IsAccountActive:   true,
		ReferralCode:      testMetadata.SystemUserID,
		EmailAddresses: []Email{
			{Id: testUser.Emails[0].EmailID, Address: testUser.Emails[0].Email, Verified: false},
			{Id: testUser.Emails[1].EmailID, Address: testUser.Emails[1].Email, Verified: true},
		},
		OAuthAccounts: []OAuthAccount{
			{
				OauthUserRegistrationID: "oauth-registration-id",
				ProviderSubject:         "subject",
				ProviderType:            "Github",
				ProfilePictureURL:       "https://google.com/picture.jpg",
				Locale:                  "en",
			},
		},
		TrustedMetadata: testMetadata,
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Now(),
	}

	got := ConvertStytchUserToUser(testUser)

	// unit testing time has to be duration lets overide
	want.UpdatedAt = got.UpdatedAt
	assert.Equal(t, got, want)
}
