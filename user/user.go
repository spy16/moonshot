package user

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/spy16/moonshot/errors"
	"github.com/spy16/moonshot/utils"
)

const (
	minPwdLen   = 10
	pwdHashCost = 13
)

const (
	KindAdmin = "admin"
	KindUser  = "user"
)

var (
	idPattern    = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	unamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// User represents a registered user.
type User struct {
	ID           string    `json:"id"`
	Kind         string    `json:"kind"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PasswordHash string    `json:"password_hash"`
	LastResetAt  time.Time `json:"last_reset_at"`

	// Profile data.
	Name      string `json:"name"`
	Gender    string `json:"gender,omitempty"`
	Locale    string `json:"locale,omitempty"`
	Location  string `json:"location,omitempty"`
	AvatarURL string `json:"avatar_url"`
}

// Updates to be applied on a user.
type Updates struct {
	Password string `json:"password"`
}

// Validate returns error if user value is not valid.
func (u *User) Validate() error {
	if u == nil {
		return errors.ErrInvalid.WithMsgf("user is empty/nil")
	}

	if !IsValidKind(u.Kind) {
		return errors.ErrInvalid.WithMsgf("valid kind must be set")
	}

	if !idPattern.MatchString(u.ID) {
		return errors.ErrInvalid.WithMsgf("valid id must be set")
	}

	if !unamePattern.MatchString(u.Username) {
		return errors.ErrInvalid.WithMsgf("valid username must be set")
	}

	if u.Email == "" || !utils.IsValidEmail(u.Email) {
		return errors.ErrInvalid.WithMsgf("valid email must be set")
	}

	return nil
}

// SetPassword hashes and stores the password for this user.
func (u *User) SetPassword(salt, pwd string, now time.Time) error {
	if len(pwd) < minPwdLen {
		return errors.ErrInvalid.WithMsgf("password must be at-least %d characters", minPwdLen)
	}

	b, err := bcrypt.GenerateFromPassword([]byte(salt+pwd), pwdHashCost)
	if err != nil {
		return errors.ErrInternal.WithCausef(err.Error())
	}
	u.PasswordHash = string(b)
	u.LastResetAt = now
	return nil
}

// CheckPassword checks if the given password is valid.
func (u *User) CheckPassword(salt, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(salt+pwd))
	return err == nil
}

// IsValidKind returns true if given string is valid user-kind.
func IsValidKind(kind string) bool { return kind == KindUser || kind == KindAdmin }
