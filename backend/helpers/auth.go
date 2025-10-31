package helpers

import (
	"github.com/pocketbase/pocketbase/core"
)

// GetAuthenticatedUser extracts and validates the authorization token from request headers
// Returns the authenticated user record or throws an error
func GetAuthenticatedUser(e *core.RequestEvent) (*core.Record, error) {
	// Extract and validate auth token
	authHeader := e.Request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, e.UnauthorizedError("You are not authorized to access this resource", nil)
	}

	// Find authenticated user by token
	user, err := e.App.FindAuthRecordByToken(authHeader, core.TokenTypeAuth)
	if err != nil {
		return nil, e.UnauthorizedError("Invalid or expired token", err)
	}

	return user, nil
}
