package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mrehanabbasi/supabase-auth-go/types"
)

const factorsPath = "/factors"

// POST /factors
//
// Enroll a new factor.
func (c *Client) EnrollFactor(ctx context.Context, req types.EnrollFactorRequest) (*types.EnrollFactorResponse, error) {
	if req.FactorType == "" {
		req.FactorType = types.FactorTypeTOTP
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := c.newRequest(ctx, factorsPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	var res types.EnrollFactorResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// POST /factors/{factor_id}/challenge
//
// Challenge a factor.
func (c *Client) ChallengeFactor(ctx context.Context, req types.ChallengeFactorRequest) (*types.ChallengeFactorResponse, error) {
	url := fmt.Sprintf("%s/%s/challenge", factorsPath, req.FactorID)
	r, err := c.newRequest(ctx, url, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	type decodeResp struct {
		ID     uuid.UUID `json:"id"`
		Expiry int64     `json:"expires_at"`
	}
	res := decodeResp{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	expiresAt := time.Unix(res.Expiry, 0)
	return &types.ChallengeFactorResponse{
		ID:        res.ID,
		ExpiresAt: expiresAt,
	}, nil
}

// POST /factors/{factor_id}/verify
//
// Verify the challenge for an enrolled factor.
func (c *Client) VerifyFactor(ctx context.Context, req types.VerifyFactorRequest) (*types.VerifyFactorResponse, error) {
	url := fmt.Sprintf("%s/%s/verify", factorsPath, req.FactorID)

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, newRequestEncodingError(err)
	}

	r, err := c.newRequest(ctx, url, http.MethodPost, body)
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

	var res types.VerifyFactorResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newResponseDecodingError(err)
	}
	return &res, nil
}

// DELETE /factors/{factor_id}
//
// Unenroll an enrolled factor.
func (c *Client) UnenrollFactor(ctx context.Context, req types.UnenrollFactorRequest) (*types.UnenrollFactorResponse, error) {
	url := fmt.Sprintf("%s/%s", factorsPath, req.FactorID)

	r, err := c.newRequest(ctx, url, http.MethodDelete, nil)
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

	var res types.UnenrollFactorResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newResponseDecodingError(err)
	}
	return &res, nil
}
