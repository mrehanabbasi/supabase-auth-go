package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const tokenPath = "/token"

// Sign in with email and password
//
// This is a convenience method that calls Token with the password grant type
func (c *Client) SignInWithEmailPassword(ctx context.Context, email, password string) (*types.TokenResponse, error) {
	return c.Token(ctx, types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  password,
	})
}

// Sign in with phone and password
//
// This is a convenience method that calls Token with the password grant type
func (c *Client) SignInWithPhonePassword(ctx context.Context, phone, password string) (*types.TokenResponse, error) {
	return c.Token(ctx, types.TokenRequest{
		GrantType: "password",
		Phone:     phone,
		Password:  password,
	})
}

// Sign in with refresh token
//
// This is a convenience method that calls Token with the refresh_token grant type
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*types.TokenResponse, error) {
	return c.Token(ctx, types.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
}

// Sign in with id token
//
// This is a convenience method that calls Token with the id_token grant type
func (c *Client) SignInWithIdToken(
	ctx context.Context,
	provider, idToken, nonce, accessToken, captchaToken string,
) (*types.TokenResponse, error) {
	return c.Token(ctx, types.TokenRequest{
		GrantType:   "id_token",
		IdToken:     idToken,
		Nonce:       nonce,
		Provider:    provider,
		AccessToken: accessToken,
		SecurityEmbed: types.SecurityEmbed{
			Security: types.GoTrueMetaSecurity{
				CaptchaToken: captchaToken,
			},
		},
	})
}

// POST /token
//
// This is an OAuth2 endpoint that currently implements the password,
// refresh_token, and PKCE grant types
func (c *Client) Token(ctx context.Context, req types.TokenRequest) (*types.TokenResponse, error) {
	switch req.GrantType {
	case "password":
		if (req.Email == "" && req.Phone == "") || req.Password == "" || req.RefreshToken != "" {
			return nil, types.ErrInvalidTokenRequest
		}
	case "refresh_token":
		if req.RefreshToken == "" || req.Email != "" || req.Phone != "" || req.Password != "" {
			return nil, types.ErrInvalidTokenRequest
		}
	case "pkce":
		if req.Code == "" || req.CodeVerifier == "" {
			return nil, types.ErrInvalidTokenRequest
		}
	case "id_token":
		if req.IdToken == "" {
			return nil, types.ErrInvalidTokenRequest
		}

		if req.Provider == "" ||
			(req.Provider != "github" &&
				req.Provider != "apple" &&
				req.Provider != "kakao" &&
				req.Provider != "keycloak") {
			return nil, types.ErrInvalidProviderRequest
		}
	default:
		return nil, types.ErrInvalidTokenRequest
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, tokenPath+"?grant_type="+req.GrantType, http.MethodPost, body)
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

	var res types.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
