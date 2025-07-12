package endpoints

import (
	"context"
	"net/http"
)

const reauthenticatePath = "/reauthenticate"

// GET /reauthenticate
//
// Sends a nonce to the user's email (preferred) or phone. This endpoint
// requires the user to be logged in / authenticated first. The user needs to
// have either an email or phone number for the nonce to be sent successfully.
func (c *Client) Reauthenticate(ctx context.Context) error {
	r, err := c.newRequest(ctx, reauthenticatePath, http.MethodGet, nil)
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
