package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/users"
)

type (
	Email struct {
		Id       string
		Address  string
		Verified bool
	}

	OAuthAccount struct {
		OauthUserRegistrationID string `json:"oauth_user_registration_id"`
		ProviderSubject         string `json:"provider_subject"`
		ProviderType            string `json:"provider_type"`
		ProfilePictureURL       string `json:"profile_picture_url"`
		Locale                  string `json:"locale"`
	}

	User struct {
		UserID            string
		SytstemUserID     string
		FirstName         *string
		MiddleName        *string
		LastName          *string
		FullName          string
		Email             *string
		PhoneNumber       *string
		TotpId            *string
		EmailIsVerified   bool
		PhoneIsVerified   bool
		PasswordIsEnabled bool
		TotpIsEnabled     bool
		IsAccountActive   bool
		ReferralCode      string
		EmailAddresses    []Email
		OAuthAccounts     []OAuthAccount
		TrustedMetadata   TrustedMetadata
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}
)

func ConvertStytchUserToUser(user users.User) User {
	var (
		converter                                 = converter{stytch: &user}
		metadata                                  = converter.getTrustedMetadata()
		totpIsEnabled, totpIdOrHash               = converter.getTOTP()
		email, emailIsVerified                    = converter.getPrimaryEmail()
		nameFirst, nameMiddle, nameLast, nameFull = converter.getNames()
	)

	return User{
		UserID:            user.UserID,
		SytstemUserID:     metadata.SystemUserID,
		FirstName:         nameFirst,
		MiddleName:        nameMiddle,
		LastName:          nameLast,
		FullName:          nameFull,
		Email:             email,
		PhoneNumber:       metadata.UserPhoneNumber,
		TotpId:            totpIdOrHash,
		EmailIsVerified:   emailIsVerified,
		PhoneIsVerified:   converter.getIsPhoneVerifiedFromMetadata(metadata.UserPhoneNumber),
		PasswordIsEnabled: converter.getIsPasswordEnabled(),
		TotpIsEnabled:     totpIsEnabled,
		IsAccountActive:   converter.getIsAccountActive(),
		ReferralCode:      metadata.SystemUserID,
		EmailAddresses:    converter.getEmailAddresses(),
		OAuthAccounts:     converter.getOAuth(),
		TrustedMetadata:   converter.getTrustedMetadata(),
		CreatedAt:         converter.getCreatedAt(),
		UpdatedAt:         converter.getUpdatedAt(),
	}
}

type converter struct {
	stytch *users.User
}

func (e *converter) getIsAccountActive() bool {
	return strings.EqualFold(e.stytch.Status, "active")
}

func (e *converter) getIsPhoneVerifiedFromMetadata(phone *string) bool {
	return !(phone == nil)
}

func (e *converter) getIsPasswordEnabled() bool {
	return !(e.stytch.Password == nil)
}

// stytch doc states: there can be only one OTP array
// Ref: https://stytch.com/docs/api/errors/400#active_totp_exists
func (e *converter) getTOTP() (bool, *string) {
	if len(e.stytch.TOTPs) == 0 {
		return false, nil
	}
	return e.stytch.TOTPs[0].Verified, &e.stytch.TOTPs[0].TOTPID
}

func (e *converter) getPrimaryEmail() (*string, bool) {
	if len(e.stytch.Emails) == 0 {
		return nil, false
	}

	if len(e.stytch.Emails) == 1 {
		return &e.stytch.Emails[0].Email, e.stytch.Emails[0].Verified
	}

	for i := range e.stytch.Emails {
		if e.stytch.Emails[i].Verified {
			return &e.stytch.Emails[i].Email, e.stytch.Emails[i].Verified
		}
	}

	return &e.stytch.Emails[0].Email, e.stytch.Emails[0].Verified
}

func (e *converter) getTrustedMetadata() TrustedMetadata {
	if len(e.stytch.TrustedMetadata) == 0 || e.stytch.TrustedMetadata == nil {
		return TrustedMetadata{}
	}

	metadata, err := DecodeFromXToX[TrustedMetadata](e.stytch.TrustedMetadata, true)
	if err != nil {
		panic(err)
	}

	return *metadata
}

func (e *converter) getCreatedAt() time.Time {
	if e.stytch.CreatedAt == nil {
		return time.Now()
	}
	return *e.stytch.CreatedAt
}

func (e *converter) getUpdatedAt() time.Time {
	return time.Now()
}

func (e *converter) getEmailAddresses() []Email {
	if len(e.stytch.Emails) == 0 {
		return nil
	}

	var newEmail []Email

	for i := range e.stytch.Emails {
		newEmail = append(newEmail, Email{
			Id:       e.stytch.Emails[i].EmailID,
			Address:  e.stytch.Emails[i].Email,
			Verified: e.stytch.Emails[i].Verified,
		})
	}

	return newEmail
}

func (e *converter) getOAuth() []OAuthAccount {
	if len(e.stytch.Providers) == 0 {
		return nil
	}

	var oauths []OAuthAccount

	for i := range e.stytch.Providers {
		oauths = append(oauths, OAuthAccount{
			OauthUserRegistrationID: e.stytch.Providers[i].OAuthUserRegistrationID,
			ProviderSubject:         e.stytch.Providers[i].ProviderSubject,
			ProviderType:            e.stytch.Providers[i].ProviderType,
			ProfilePictureURL:       e.stytch.Providers[i].ProfilePictureURL,
			Locale:                  e.stytch.Providers[i].Locale,
		})
	}

	return oauths
}

func (e *converter) getNames() (firstName, middleName, lastName *string, fullName string) {
	if e.stytch.Name == nil {
		return nil, nil, nil, ""
	}

	var (
		fnPtr, fn = handleNames(e.stytch.Name.FirstName)
		mPtr, m   = handleNames(e.stytch.Name.MiddleName)
		lPtr, l   = handleNames(e.stytch.Name.LastName)
	)

	fullName = fmt.Sprintf("%s %s %s", fn, m, l)

	return fnPtr, mPtr, lPtr, strings.TrimSpace(fullName)
}
