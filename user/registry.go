package user

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Registry provides user management functions.
type Registry struct {
	Store     Store
	Providers []Provider
}

// Provider represents a OAuth2/OpenID login provider.
type Provider struct {
	Name   string        `json:"name"`
	Config oauth2.Config `json:"config"`
}

// Store implementations provide persistence for user data.
type Store interface {
	// Get should return user by id or email (if isEmail=true). If not found
	// ErrNotFound should be returned.
	Get(ctx context.Context, idOrEmail string, isEmail bool) (*User, error)

	// Put should create/update the given user. If onlyUpdate and user with
	// given id or email already exists, then ErrConflict must be returned.
	Put(ctx context.Context, u User, onlyUpdate bool) error
}

// Get fetches a user by id/email. Returns ErrUnauthorized if current user from
// context is not authorized.
func (reg *Registry) Get(ctx context.Context, idOrEmail string, isEmail bool) (*User, error) {
	// TODO: authorize this action (only admins & user himself should be able to do this).

	u, err := reg.Store.Get(ctx, idOrEmail, isEmail)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Register validates 'u' and registers the user if possible.
func (reg *Registry) Register(ctx context.Context, u User) (*User, error) {
	if err := u.Sanitize(true); err != nil {
		return nil, err
	}

	if err := reg.Store.Put(ctx, u, false); err != nil {
		return nil, err
	}

	return &u, nil
}

// Authenticate performs validation of session-cookie/access-token sent in the
// request and returns a new request with user injected into context. Boolean
// flag indicates whether user was authenticated (i.e., a valid token found).
// An error is returned only in case of some internal failures.
func (reg *Registry) Authenticate(r *http.Request) (*http.Request, bool, error) {
	tokenStr := findToken(r)
	if tokenStr == "" {
		return r, false, nil
	}

	// TODO: validate token and populate user.
	authenticatedUser := User{}
	authCtx := newContext(r.Context(), authenticatedUser)

	return r.WithContext(authCtx), true, nil
}

func findToken(r *http.Request) string {
	const (
		bearerPrefix = "Bearer "
		cookieName   = "moonshot:session"
		authzHeader  = "Authorization"
	)

	var tokenStr string
	if value := r.Header.Get(authzHeader); strings.HasPrefix(value, bearerPrefix) {
		tokenStr = strings.TrimPrefix(value, bearerPrefix)
	} else {
		cookie, _ := r.Cookie(cookieName)
		tokenStr = cookie.Value
	}
	return strings.TrimSpace(tokenStr)
}
