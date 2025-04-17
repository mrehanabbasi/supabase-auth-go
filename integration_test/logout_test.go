package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

func TestLogout(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := autoconfirmClient

	// Logout without a valid token.
	err := client.Logout(ctx)
	assert.Error(err)

	// Create logged in user.
	email := randomEmail()
	password := randomString(10)
	session, err := client.Signup(ctx, types.SignupRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(err)

	// Logout.
	err = client.WithToken(session.AccessToken).Logout(ctx)
	require.NoError(err)

	// Attempt refresh - should fail.
	_, err = client.RefreshToken(ctx, session.RefreshToken)
	assert.Error(err)
}
