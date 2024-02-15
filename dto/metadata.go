package dto

import (
	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
)

type TrustedMetadata struct {
	SystemUserID       string  `mapstructure:"system_user_id"`
	UserRole           string  `mapstructure:"user_role"`
	UserPhoneNumber    *string `mapstructure:"user_phone_number"`
	UserProfileImage   *string `mapstructure:"user_profile_image"`
	NotificationEmail  bool    `mapstructure:"notification_email"`
	NotificationPush   bool    `mapstructure:"notification_push"`
	NotificationSMS    bool    `mapstructure:"notification_sms"`
	NotificationInApp  bool    `mapstructure:"notification_in_app"`
	TransactionPinHash *string `mapstructure:"transaction_pin_hash"`
	WebhookURL         *string `mapstructure:"webhook_url"`
}

var DefaultTrustedMetadata = TrustedMetadata{
	SystemUserID:       gonanoid.MustGenerate("0123456789abcdefghijklmnopqrstuvwxyz", 10),
	UserRole:           "customer",
	UserPhoneNumber:    nil,
	UserProfileImage:   nil,
	NotificationEmail:  true,
	NotificationPush:   true,
	NotificationSMS:    true,
	NotificationInApp:  true,
	TransactionPinHash: nil,
	WebhookURL:         nil,
}

type (
	Name struct {
		FirstName  string
		MiddleName string
		LastName   string
	}

	UpdateUserParams struct {
		// If name is empty it wouldnt be updated
		Name               Name
		UserRole           *string
		UserPhoneNumber    *string
		UserProfileImage   *string
		NotificationEmail  *bool
		NotificationPush   *bool
		NotificationSMS    *bool
		NotificationInApp  *bool
		TransactionPinHash *string
		WebhookURL         *string
	}
)

func (u *UpdateUserParams) UpdateWith(p TrustedMetadata) (TrustedMetadata, error) {
	if u.UserRole != nil {
		p.UserRole = *u.UserRole
	}

	if u.UserPhoneNumber != nil {
		p.UserPhoneNumber = u.UserPhoneNumber
	}

	if u.UserProfileImage != nil {
		p.UserProfileImage = u.UserProfileImage
	}

	if u.NotificationEmail != nil {
		p.NotificationEmail = *u.NotificationEmail
	}

	if u.NotificationPush != nil {
		p.NotificationPush = *u.NotificationPush
	}

	if u.NotificationSMS != nil {
		p.NotificationSMS = *u.NotificationSMS
	}

	if u.NotificationInApp != nil {
		p.NotificationInApp = *u.NotificationInApp
	}

	if u.TransactionPinHash != nil {
		hashed, err := hashPin(*u.TransactionPinHash)
		if err != nil {
			return TrustedMetadata{}, err
		}
		p.TransactionPinHash = &hashed
	}

	if u.WebhookURL != nil {
		p.WebhookURL = u.WebhookURL
	}

	return p, nil
}

// HashPin takes a plain text pin and returns a hashed version using bcrypt.
func hashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(bytes), err
}
