package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

var settingsPath = "/settings"

// GET /settings
//
// Returns the publicly available settings for this auth instance.
func (c *Client) GetSettings(ctx context.Context) (*types.SettingsResponse, error) {
	r, err := c.newRequest(ctx, settingsPath, http.MethodGet, nil)
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

	var res types.SettingsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, newResponseDecodingError(err)
	}

	return &res, nil
}
