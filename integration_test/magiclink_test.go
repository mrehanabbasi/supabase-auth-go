package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

func TestMagiclink(t *testing.T) {
	assert := assert.New(t)

	email := randomEmail()
	err := client.Magiclink(ctx, types.MagiclinkRequest{
		Email: email,
	})
	assert.NoError(err)

	err = client.Magiclink(ctx, types.MagiclinkRequest{})
	assert.Error(err)
}
