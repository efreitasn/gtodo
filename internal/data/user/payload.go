package user

import "context"

// Payload is the user data that can used in a token.
type Payload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ctxKey string

var payloadCtxKey ctxKey = "payload"

// PayloadFromContext gets a payload from a context.
func PayloadFromContext(ctx context.Context) *Payload {
	if payload, ok := ctx.Value(payloadCtxKey).(*Payload); ok {
		return payload
	}

	return nil
}

// PayloadFromUser gets a payload from an user.
func PayloadFromUser(u *User) *Payload {
	return &Payload{
		ID:       u.ID.Hex(),
		Username: u.Username,
	}
}

// ContextWithPayload adds a payload to a context.
func ContextWithPayload(ctx context.Context, p *Payload) context.Context {
	return context.WithValue(ctx, payloadCtxKey, p)
}
