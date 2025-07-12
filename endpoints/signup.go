package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const signupPath = "/signup"

// POST /signup
//
// Register a new user with an email and password.
func (c *Client) Signup(ctx context.Context, req types.SignupRequest) (*types.SignupResponse, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, signupPath, http.MethodPost, body)
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

	var res types.SignupResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newResponseDecodingError(err)
	}

	// To make the response easier to consume, if autoconfirm was enabled, the
	// session user should be populated. Copy that into the embedded user type
	// so it's easier to access.
	//
	// i.e. we can access user fields like res.Email regardless of whether
	// we got back a session or a user.
	if res.Session.User.ID != uuid.Nil {
		res.User = res.Session.User
	}

	return &res, nil
}
