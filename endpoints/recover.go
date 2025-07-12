package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const recoverPath = "/recover"

// POST /recover
//
// Password recovery. Will deliver a password recovery mail to the user based
// on email address.
//
// By default recovery links can only be sent once every 60 seconds.
func (c *Client) Recover(ctx context.Context, req types.RecoverRequest) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, recoverPath, http.MethodPost, body)
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
