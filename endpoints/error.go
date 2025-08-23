package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errCodeUserAlreadyExists          = "user_already_exists"
	errCodeInvalidCredentials         = "invalid_credentials"
	errCodeSessionNotFound            = "session_not_found"
	errCodeBadJWT                     = "bad_jwt"
	errCodeEmailNotConfirmed          = "email_not_confirmed"
	errCodeUnexpectedFailure          = "unexpected_failure"
	errCodeEmailSendRateLimitExceeded = "over_email_send_rate_limit"
	errCodeNoAuthorization            = "no_authorization"
	errCodeNotAdmin                   = "not_admin"
	errCodeValidationFailed           = "validation_failed"
	errCodeRefreshTokenNotFound       = "refresh_token_not_found"

	errMsgErrorSendingConfirmationEmail = "Error sending confirmation email"
	errMsguserIDMustBeUUID              = "user_id must be an UUID"

	ErrUserAlreadyExists              = errors.New("user already exists")
	ErrInvalidCredentials             = errors.New("invalid credentials")
	ErrSessionNotFound                = errors.New("session not found")
	ErrInvalidJWT                     = errors.New("invalid jwt")
	ErrEmailNotConfirmed              = errors.New("email not confirmed")
	ErrRedirectURLNotInResponse       = errors.New("no redirect URL found in response")
	ErrFailedSendingConfirmationEmail = errors.New("failed sending confirmation email")
	ErrEmailSendLimitExceeded         = errors.New("email send limit exceeded")
	ErrNoAuthorization                = errors.New("no authorization")
	ErrNotAdmin                       = errors.New("not admin")
	ErrInvalidUserID                  = errors.New("invalid user id")
	ErrRefreshTokenNotFound           = errors.New("refresh token not found")

	distinctErrors = map[string]error{
		errCodeUserAlreadyExists:          ErrUserAlreadyExists,
		errCodeInvalidCredentials:         ErrInvalidCredentials,
		errCodeSessionNotFound:            ErrSessionNotFound,
		errCodeBadJWT:                     ErrInvalidJWT,
		errCodeEmailNotConfirmed:          ErrEmailNotConfirmed,
		errCodeEmailSendRateLimitExceeded: ErrEmailSendLimitExceeded,
		errCodeNoAuthorization:            ErrNoAuthorization,
		errCodeNotAdmin:                   ErrNotAdmin,
		errCodeRefreshTokenNotFound:       ErrRefreshTokenNotFound,
	}
)

type WeakPasswordReason string

const (
	WeakPasswordReasonLength     = "length"
	WeakPasswordReasonCharacters = "characters"
	WeakPasswordReasonsPwned     = "pwned"
)

type ErrorResponse struct {
	Err          *string       `json:"error"`
	ErrDesc      *string       `json:"error_description"`
	Code         *int          `json:"code"`
	Message      *string       `json:"msg"`
	ErrorCode    *string       `json:"error_code"`
	WeakPassword *WeakPassword `json:"weak_password"`
}

type WeakPassword struct {
	Reasons []WeakPasswordReason `json:"reasons"`
}

func (e ErrorResponse) Error() string {
	if e.Message != nil {
		return *e.Message
	}
	if e.Err != nil {
		return *e.Err
	}

	return ""
}

func (e ErrorResponse) getDistinctError() error {
	if e.ErrorCode == nil {
		if *e.ErrorCode == errCodeUnexpectedFailure &&
			e.Message != nil && *e.Message == errMsgErrorSendingConfirmationEmail {
			return ErrFailedSendingConfirmationEmail
		}

		if *e.ErrorCode == errCodeValidationFailed &&
			e.Message != nil && *e.Message == errMsguserIDMustBeUUID {
			return ErrInvalidUserID
		}
	}

	for k, v := range distinctErrors {
		if e.ErrorCode != nil {
			if k == *e.ErrorCode {
				return v
			}
		}
	}

	return nil
}

func newRequestEncodingError(err error) error {
	return fmt.Errorf("failed to encode request body: %w", err)
}

func newRequestCreationError(err error) error {
	return fmt.Errorf("failed to create request: %w", err)
}

func newRequestDispatchError(err error) error {
	return fmt.Errorf("failed to send request: %w", err)
}

func newResponseDecodingError(err error) error {
	return fmt.Errorf("failed to decode response body: %w", err)
}

func newErrorResponseDecodingError(status string, err error) error {
	return fmt.Errorf("failed to decode error response body (%s): %w", status, err)
}

func wrapError(status string, err error) error {
	return fmt.Errorf("supabase-auth - %s: %w", status, err)
}

func handleErrorResponse(resp *http.Response) error {
	var errRes ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errRes); err != nil {
		return newErrorResponseDecodingError(resp.Status, err)
	}

	if err := errRes.getDistinctError(); err != nil {
		return err
	}

	return wrapError(resp.Status, errRes)
}
