package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const magiclinkPath = "/magiclink"

// POST /magiclink
//
// DEPRECATED: Use /otp with Email and CreateUser=true instead of /magiclink.
//
// Magic Link. Will deliver a link (e.g.
// /verify?type=magiclink&token=fgtyuf68ddqdaDd) to the user based on email
// address which they can use to redeem an access_token.
//
// By default Magic Links can only be sent once every 60 seconds.
func (c *Client) Magiclink(ctx context.Context, req types.MagiclinkRequest) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, magiclinkPath, http.MethodPost, body)
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
