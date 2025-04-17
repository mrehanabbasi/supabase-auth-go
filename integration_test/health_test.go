package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	supaAuth "github.com/mrehanabbasi/supabase-auth-go"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := supaAuth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9999")
	health, err := client.HealthCheck(ctx)
	require.NoError(err)
	assert.Equal(health.Name, "GoTrue")
}
