package integration_test

import (
	"testing"

	supaAuth "github.com/mrehanabbasi/supabase-auth-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	client := supaAuth.New(projectReference, apiKey).WithCustomAuthURL("http://localhost:9999")
	health, err := client.HealthCheck()
	require.NoError(err)
	assert.Equal(health.Name, "GoTrue")
}
