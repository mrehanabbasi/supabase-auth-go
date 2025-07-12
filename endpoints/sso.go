package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const ssoPath = "/sso"

// POST /sso
//
// Initiate an SSO session with the given provider.
//
// If successful, the server returns a redirect to the provider's authorization
// URL. The client will follow it and return the final HTTP response.
//
// Auth allows you to skip following the redirect by setting SkipHTTPRedirect
// on the request struct. In this case, the URL to redirect to will be returned
// in the response.
func (c *Client) SSO(ctx context.Context, req types.SSORequest) (*types.SSOResponse, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, http.MethodPost, ssoPath, body)
	if err != nil {
		return nil, newRequestCreationError(err)
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, newRequestDispatchError(err)
	}

	if !req.SkipHTTPRedirect {
		// If the client is following redirects, we can return the response
		// directly.
		return &types.SSOResponse{
			HTTPResponse: resp,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		return nil, handleErrorResponse(resp)
	}

	// If the client is not following redirects, we can unmarshal the response from
	// the server to get the URL.
	var res types.SSOResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, newResponseDecodingError(err)
	}
	return &res, nil
}
