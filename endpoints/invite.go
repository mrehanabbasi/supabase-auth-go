package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const invitePath = "/invite"

// POST /invite
//
// Invites a new user with an email.
// This endpoint requires the service_role or supabase_admin JWT set using WithToken.
func (c *Client) Invite(ctx context.Context, req types.InviteRequest) (*types.InviteResponse, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, invitePath, http.MethodPost, body)
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

	var res types.InviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newResponseDecodingError(err)
	}

	return &res, nil
}
