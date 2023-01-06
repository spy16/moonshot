package user

import "context"

const (
	KeyID = iota
	KeyEmail
	KeyUsername
)

// Store implementation is responsible for storage of users.
type Store interface {
	GetUser(ctx context.Context, kind, key string, keyType int) (*User, error)
	PutUser(ctx context.Context, u User, onlyCreate bool) error
	DelUser(ctx context.Context, kind, id string) error
	ListUsers(ctx context.Context, kind string) ([]User, error)
}
