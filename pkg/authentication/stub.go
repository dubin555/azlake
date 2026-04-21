package authentication

import "errors"

var (
	ErrNotImplemented          = errors.New("not implemented")
	ErrSessionExpired          = errors.New("session expired")
	ErrInvalidTokenFormat      = errors.New("invalid token format")
	ErrInvalidRequest          = errors.New("invalid request")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)

type Service interface {
	// Login returns a token for the given credentials
	Login(ctx interface{}, loginRequest LoginRequest) (*LoginResponse, error)
	// Logout invalidates the token
	Logout(ctx interface{}, token string) error
	// GetUserByToken returns the user for the given token
	GetUserByToken(ctx interface{}, token string) (*UserData, error)
	// SupportsSessions returns true if the service supports sessions
	SupportsSessions() bool
}

type LoginRequest struct {
	AccessKeyID     string
	SecretAccessKey  string
}

type LoginResponse struct {
	Token               string
	TokenExpiration     int64
	InternalAuthSession bool
}

type UserData struct {
	Username string
}
