package govpsie

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// APIError is the error returned for a non-2xx API answer. Its Error text is
// the API's message, or "vpsie: unexpected response: <status>" when the
// answer carried none (an HTML page from a proxy, say). Use errors.As to read
// the HTTP status and the API's error code:
//
//	var apiErr *govpsie.APIError
//	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict { ... }
type APIError struct {
	// StatusCode is the HTTP status of the answer.
	StatusCode int
	// Code is the API's machine-readable error code when it sent one, for
	// example "bareMetalBusy" or "private_network_busy". The API sends some
	// codes as numbers; those keep their textual form.
	Code string
	// Message is the API's human-readable message.
	Message string
}

// Error implements error.
func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("vpsie: unexpected response: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// newAPIError builds the error for a non-2xx answer from its status and body.
// A body that is not the JSON error envelope keeps only the status.
func newAPIError(status int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: status}

	var envelope struct {
		Code    json.RawMessage `json:"code"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return apiErr
	}

	apiErr.Message = envelope.Message
	if len(envelope.Code) > 0 {
		var code FlexString
		if err := json.Unmarshal(envelope.Code, &code); err == nil {
			apiErr.Code = code.String()
		}
	}

	return apiErr
}

// IsNotFound reports whether err is an API answer with HTTP status 404.
func IsNotFound(err error) bool {
	var apiErr *APIError

	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// ErrorCode returns the API error code carried by err, or "" when err is not
// an *APIError or the answer had no code.
func ErrorCode(err error) string {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}

	return ""
}
