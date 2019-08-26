package auth

import (
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
	"go.mongodb.org/mongo-driver/mongo"
)

// Auth is the representation of all the auth-related middlewares.
type Auth struct {
	c *mongo.Collection
}

// New creates an Auth struct.
func New(c *mongo.Collection) *Auth {
	return &Auth{c}
}

// Messages
var notAuthErrorMessage = &flash.Message{
	Kind:    1,
	Content: "You have to be authenticated.",
}

var dbErrrorMessage = &flash.Message{
	Kind:    1,
	Content: "Error while connecting to the database",
}

// HasToBeAuth checks if the user is authenticated to go to the next http.HandlerFunc.
func (a *Auth) HasToBeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userPayload := user.PayloadFromTokenCookie(r)

		if userPayload == nil {
			flash.Add("/login", w, r, notAuthErrorMessage)

			return
		}

		exist, err := a.userExist(userPayload.ID)

		if err != nil {
			flash.Add("/login", w, r, dbErrrorMessage)

			return
		}

		if !exist {
			flash.Add("/login", w, r, notAuthErrorMessage)

			return
		}

		newR := r.WithContext(user.ContextWithPayload(r.Context(), userPayload))

		next(w, newR)
	}
}

// HasToBeUnauth checks if the user is unauthenticated to go to the next http.HandlerFunc.
func (a *Auth) HasToBeUnauth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userPayload := user.PayloadFromTokenCookie(r)

		if userPayload == nil {
			next(w, r)

			return
		}

		exist, err := a.userExist(userPayload.ID)

		if err != nil {
			flash.Add("/login", w, r, dbErrrorMessage)

			return
		}

		if !exist {
			next(w, r)

			return
		}

		http.Redirect(w, r, "/list", 303)
	}
}
