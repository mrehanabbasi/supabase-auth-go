package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

var healthPath = "/health"

// GET /health
//
// Check the health of the Auth server.
func (c *Client) HealthCheck(ctx context.Context) (*types.HealthCheckResponse, error) {
	r, err := c.newRequest(ctx, healthPath, http.MethodGet, nil)
	if err != nil {
		return nil, newRequestCreationError(err)
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, newRequestDispatchError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleErrorResponse(resp)
	}

	var res types.HealthCheckResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, newResponseDecodingError(err)
	}

	return &res, nil
}
