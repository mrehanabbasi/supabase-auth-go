package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.SettingsResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
