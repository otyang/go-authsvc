package dto

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

func DecodeFromXToX[T any](input any, weaklyTypedInput bool) (*T, error) {
	var result T

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: weaklyTypedInput,
		Result:           &result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err = decoder.Decode(input); err != nil {
		return nil, err
	}

	return &result, nil
}

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	v, ok := err.(stytcherror.Error)
	if !ok {
		return err
	}

	// replace 'support@stytch.com' in error messages to 'us'
	msg := strings.Replace(string(v.ErrorMessage), "support@stytch.com", "us", -1)

	return stytcherror.Error{
		StatusCode:   v.StatusCode,
		RequestID:    v.RequestID,
		ErrorType:    v.ErrorType,
		ErrorMessage: stytcherror.Message(msg),
		ErrorURL:     v.ErrorURL,
	}
}

func handleNames(name string) (*string, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ""
	}
	return &name, name
}
