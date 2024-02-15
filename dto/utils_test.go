package dto

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

func TestDecodeFromXToX(t *testing.T) {
	type SessionClaims struct {
		UserAgent string
		IPAddress string
		Location  string
		Platform  string // web or mobile
	}

	t.Run("input map", func(t *testing.T) {
		want := &SessionClaims{
			UserAgent: "Chrome",
			IPAddress: "127.0.0.1",
			Location:  "San Francisco",
			Platform:  "mobile",
		}

		input := map[string]any{
			"userAgent": "Chrome",
			"ipAddress": "127.0.0.1",
			"Location":  "San Francisco",
			"Platform":  "mobile",
		}

		got, err := DecodeFromXToX[SessionClaims](input, true)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("input struct", func(t *testing.T) {
		want := map[string]any{
			"UserAgent": "Chrome",
			"IPAddress": "127.0.0.1",
			"Location":  "San Francisco",
			"Platform":  "mobile",
		}

		input := SessionClaims{
			UserAgent: "Chrome",
			IPAddress: "127.0.0.1",
			Location:  "San Francisco",
			Platform:  "mobile",
		}

		got, err := DecodeFromXToX[map[string]any](input, true)

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})
}

func TestHandleError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := HandleError(nil)
		assert.NoError(t, err)
	})

	t.Run("non stytch error", func(t *testing.T) {
		err := errors.New("non-Stytch error")
		got := HandleError(err)
		assert.Equal(t, err, got)
	})

	t.Run("stytch error replacement", func(t *testing.T) {
		originalErr := stytcherror.Error{
			ErrorMessage: stytcherror.Message("Contact support@stytch.com for help."),
			// Other fields as needed
		}

		result := HandleError(originalErr)
		expectedErr := stytcherror.Error{
			ErrorMessage: stytcherror.Message("Contact us for help."),
		}

		assert.ErrorIs(t, result, expectedErr)
	})

	t.Run("stytch error no-replacement", func(t *testing.T) {
		originalErr := stytcherror.Error{
			ErrorMessage: stytcherror.Message("General error message"),
		}

		result := HandleError(originalErr)
		assert.ErrorIs(t, result, originalErr)
	})
}
