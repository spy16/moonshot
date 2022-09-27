package user

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/spy16/moonshot/pkg/errors"
)

// User represents a registered user in moonshot.
type User struct {
	// Account data.
	ID         string            `json:"id"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	VerifiedAt *time.Time        `json:"verified_at,omitempty"`
	Attributes map[string]string `json:"-"`

	// Profile data.
	Name      string `json:"name,omitempty"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// IsVerified returns true if the user has been verified. OAuth2/OpenID
// based logins are considered verified by default.
func (u *User) IsVerified() bool {
	return u.VerifiedAt != nil && !u.VerifiedAt.IsZero()
}

// Sanitize sanitizes the user value by setting defaults where possible.
// Returns ErrInvalid if required fields are missing/invalid.
func (u *User) Sanitize(isCreate bool) error {
	if u.Email == "" {
		return errors.ErrInvalid.WithMsgf("email must be present")
	}

	if isCreate {
		u.ID = randomString(8)
		u.CreatedAt = time.Now()
		u.UpdatedAt = u.CreatedAt
		if u.AvatarURL == "" {
			u.AvatarURL = gravatarURL(u.Email, 128)
		}
	}

	return nil
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func gravatarURL(email string, size int) string {
	if size <= 0 || size >= 2048 {
		size = 128
	}
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}
