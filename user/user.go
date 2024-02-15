package user

import (
	"context"

	"github.com/otyang/go-authsvc/dto"

	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/users"
)

type UserService struct {
	client *stytchapi.API
}

func NewUserService(client *stytchapi.API) *UserService {
	return &UserService{
		client: client,
	}
}

func (u *UserService) Get(ctx context.Context, userID string) (*dto.User, error) {
	resp, err := u.client.Users.Get(context.Background(), &users.GetParams{
		UserID: userID,
	})
	if err != nil {
		return nil, dto.HandleError(err)
	}

	uv := users.User{
		UserID:                 resp.UserID,
		Emails:                 resp.Emails,
		Status:                 resp.Status,
		PhoneNumbers:           resp.PhoneNumbers,
		WebAuthnRegistrations:  resp.WebAuthnRegistrations,
		Providers:              resp.Providers,
		TOTPs:                  resp.TOTPs,
		CryptoWallets:          resp.CryptoWallets,
		BiometricRegistrations: resp.BiometricRegistrations,
		Name:                   resp.Name,
		CreatedAt:              resp.CreatedAt,
		Password:               resp.Password,
		TrustedMetadata:        resp.TrustedMetadata,
		UntrustedMetadata:      resp.UntrustedMetadata,
	}

	user := dto.ConvertStytchUserToUser(uv)

	return &user, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, userID string, param *dto.UpdateUserParams) (*dto.User, error) {
	user, err := u.Get(ctx, userID)
	if err != nil {
		return nil, dto.HandleError(err)
	}

	if param == nil {
		return user, nil
	}

	updatedUserTrustedMetadata, err := param.UpdateWith(user.TrustedMetadata)
	if err != nil {
		return nil, err
	}

	data, err := dto.DecodeFromXToX[map[string]any](updatedUserTrustedMetadata, true)
	if err != nil {
		return nil, err
	}

	// update the user personal datas
	_, err = u.client.Users.Update(ctx, &users.UpdateParams{
		UserID: userID,
		Name: &users.Name{
			FirstName:  param.Name.FirstName,
			MiddleName: param.Name.MiddleName,
			LastName:   param.Name.LastName,
		},
		TrustedMetadata: *data,
	})

	return user, dto.HandleError(err)
}
