package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User is an user in the database.
type User struct {
	ID        *primitive.ObjectID `bson:"_id"`
	Name      string              `bson:"name"`
	Username  string              `bson:"username"`
	Email     string              `bson:"email"`
	Password  string              `bson:"password"`
	CreatedAt time.Time           `bson:"createdAt"`
}

// InsertUser is an user to be inserted into the database.
type InsertUser struct {
	Name      string    `bson:"name"`
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"createdAt"`
}

// HashPw hashes the password.
func (iU *InsertUser) HashPw() error {
	res, err := bcrypt.GenerateFromPassword([]byte(iU.Password), 10)

	if err != nil {
		return err
	}

	iU.Password = string(res)

	return nil
}

// Payload is the user data that can used in a token.
type Payload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ctxKey string

var payloadCtxKey ctxKey = "payload"

// PayloadFromContext is ...
func PayloadFromContext(ctx context.Context) *Payload {
	if payload, ok := ctx.Value(payloadCtxKey).(*Payload); ok {
		return payload
	}

	return nil
}

// AddToContext a
func (p *Payload) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, payloadCtxKey, p)
}

// ComparePw compares a plain text password with a hashed one to check if they match.
func (u *User) ComparePw(plainTextPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPw))

	if err != nil {
		return false
	}

	return true
}

// Payload returns the user payload.
func (u *User) Payload() *Payload {
	return &Payload{
		ID:       u.ID.Hex(),
		Username: u.Username,
	}
}
