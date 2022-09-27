package user

import "context"

type ctxKey string

var curUserKey = ctxKey("cur-user")

// From returns current authenticated user from context. Returns nil
// if no user is present.
func From(ctx context.Context) *User {
	u, ok := ctx.Value(curUserKey).(User)
	if !ok {
		return nil
	}
	return &u
}

func newContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, curUserKey, u)
}
