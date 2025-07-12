package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

var userPath = "/user"

// GET /user
//
// Get the JSON object for the logged in user (requires authentication)
func (c *Client) GetUser(ctx context.Context) (*types.UserResponse, error) {
	r, err := c.newRequest(ctx, userPath, http.MethodGet, nil)
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

	var res types.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, newResponseDecodingError(err)
	}

	return &res, nil
}

// PUT /user
//
// Update a user (Requires authentication). Apart from changing email/password,
// this method can be used to set custom user data. Changing the email will
// result in a magiclink being sent out.
func (c *Client) UpdateUser(ctx context.Context, req types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, userPath, http.MethodPut, body)
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

	var res types.UpdateUserResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, newResponseDecodingError(err)
	}

	return &res, nil
}
