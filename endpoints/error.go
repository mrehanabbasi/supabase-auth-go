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

	errMsgErrorSendingConfirmationEmail = "Error sending confirmation email"
	errMsguserIDMustBeUUID              = "user_id must be an UUID"

	ErrUserAlreadyExists              = errors.New("user_already_exists")
	ErrInvalidCredentials             = errors.New("invalid_credentials")
	ErrSessionNotFound                = errors.New("session_not_found")
	ErrInvalidJWT                     = errors.New("invalid_jwt")
	ErrEmailNotConfirmed              = errors.New("email_not_confirmed")
	ErrRedirectURLNotInResponse       = errors.New("no redirect URL found in response")
	ErrFailedSendingConfirmationEmail = errors.New("failed_sending_confirmation_email")
	ErrEmailSendLimitExceeded         = errors.New("email_send_limit_exceeded")
	ErrNoAuthorization                = errors.New("no_authorization")
	ErrNotAdmin                       = errors.New("not_admin")
	ErrInvalidUserID                  = errors.New("invalid_user_id")

	distinctErrors = map[string]error{
		errCodeUserAlreadyExists:          ErrUserAlreadyExists,
		errCodeInvalidCredentials:         ErrInvalidCredentials,
		errCodeSessionNotFound:            ErrSessionNotFound,
		errCodeBadJWT:                     ErrInvalidJWT,
		errCodeEmailNotConfirmed:          ErrEmailNotConfirmed,
		errCodeEmailSendRateLimitExceeded: ErrEmailSendLimitExceeded,
		errCodeNoAuthorization:            ErrNoAuthorization,
		errCodeNotAdmin:                   ErrNotAdmin,
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
