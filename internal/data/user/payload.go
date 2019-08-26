package user

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/hako/branca"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ctxKey string

const payloadCtxKey ctxKey = "payload"
const payloadTokenCookieName = "t"

// Payload is the user data that can used in a token.
type Payload struct {
	ID       *primitive.ObjectID `json:"id"`
	Username string              `json:"username"`
}

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
		ID:       u.ID,
		Username: u.Username,
	}
}

// PayloadFromToken gets a payload from a token.
// If the token is invalid, or any other error occurs, a nil pointer is returned.
func PayloadFromToken(t string) *Payload {
	brca := branca.NewBranca(os.Getenv("BRANCA_SECRET"))
	payloadString, err := brca.DecodeToString(t)

	if err != nil {
		return nil
	}

	var payload Payload

	json.Unmarshal([]byte(payloadString), &payload)

	return &payload
}

// PayloadFromTokenCookie gets a payload from a request's token cookie.
// If the cookie doesn't exist, or any other error occurs, a nil pointer is returned.
func PayloadFromTokenCookie(r *http.Request) *Payload {
	token, err := r.Cookie(payloadTokenCookieName)

	if err != nil {
		return nil
	}

	return PayloadFromToken(token.Value)
}

// Token returns a token from a payload.
func (p *Payload) Token() (string, error) {
	stringPayload, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	brca := branca.NewBranca(os.Getenv("BRANCA_SECRET"))
	token, err := brca.EncodeToString(string(stringPayload))

	if err != nil {
		return "", err
	}

	return token, nil
}

// ContextWithPayload adds a payload to a context.
func ContextWithPayload(ctx context.Context, p *Payload) context.Context {
	return context.WithValue(ctx, payloadCtxKey, p)
}

// ResponseWithTokenCookie adds a cookie to a response with the provided payload's token.
func ResponseWithTokenCookie(w http.ResponseWriter, p *Payload) error {
	token, err := p.Token()

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     payloadTokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(time.Hour * 24),
	})

	return nil
}
