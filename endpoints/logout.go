package endpoints

import (
	"context"
	"net/http"
)

const logoutPath = "/logout"

// POST /logout
//
// Logout a user (Requires authentication).
//
// This will revoke all refresh tokens for the user. Remember that the JWT
// tokens will still be valid for stateless auth until they expires.
func (c *Client) Logout(ctx context.Context) error {
	r, err := c.newRequest(ctx, logoutPath, http.MethodPost, nil)
	if err != nil {
		return newRequestCreationError(err)
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return newRequestDispatchError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return handleErrorResponse(resp)
	}

	return nil
}
