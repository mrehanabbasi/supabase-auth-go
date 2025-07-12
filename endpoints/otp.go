package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const otpPath = "/otp"

// POST /otp
// One-Time-Password. Will deliver a magiclink or SMS OTP to the user depending
// on whether the request contains an email or phone key.
//
// If CreateUser is true, the user will be automatically signed up if the user
// doesn't exist.
func (c *Client) OTP(ctx context.Context, req types.OTPRequest) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, otpPath, http.MethodPost, body)
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
