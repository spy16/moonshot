package user

import (
	"context"
	"strings"
	"time"

	"github.com/spy16/moonshot/errors"
	"github.com/spy16/moonshot/utils"
)

type Service struct {
	Salt  string
	Clock func() time.Time
	Store Store
}

// RegisterUser registers a new user with given email. If a user with same email
// is already registered under same kind, ErrConflict will be returned.
func (svc *Service) RegisterUser(ctx context.Context, kind, email, pwd string) (*User, error) {
	timeNow := svc.Clock()
	email = strings.TrimSpace(email)

	u := User{
		ID:        utils.RandStr(16, utils.CharsetLowerNum),
		Kind:      kind,
		Email:     email,
		Username:  "user" + utils.RandStr(5, utils.CharsetNums),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		AvatarURL: utils.GravatarURL(email, 0),
	}

	if err := u.SetPassword(svc.Salt, pwd, timeNow); err != nil {
		return nil, err
	} else if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := svc.Store.PutUser(ctx, u, true); err != nil {
		if errors.Is(err, errors.ErrConflict) {
			return nil, err
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	}
	return &u, nil
}

// PasswordLogin performs login using username/email and password.
func (svc *Service) PasswordLogin(ctx context.Context, kind, unameOrEmail, pwd string) (*User, error) {
	var errAuth = errors.ErrAuth

	keyType := KeyUsername
	if utils.IsValidEmail(unameOrEmail) {
		keyType = KeyEmail
	}

	u, err := svc.Store.GetUser(ctx, kind, unameOrEmail, keyType)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errAuth
		}
		return nil, errors.ErrInternal.WithCausef(err.Error())
	} else if !u.CheckPassword(svc.Salt, pwd) {
		return nil, errAuth
	}

	return u, nil
}

// IssueTokens issues a refresh and access token for given  Existing
// refresh-token will be invalidated.
func (svc *Service) IssueTokens(ctx context.Context, u User) (access, refresh string, err error) {
	return "", "", nil
}

// VerifyToken verifies the given access token and returns the authenticated
//
//	Returns ErrAuth if token is invalid (signature invalid, expired, etc.)
func (svc *Service) VerifyToken(ctx context.Context, accessToken string) (*User, error) {
	return nil, nil
}

// RefreshTokens refreshes and returns a new access token.
func (svc *Service) RefreshTokens(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}
