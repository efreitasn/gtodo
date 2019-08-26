package middlewares

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
	"github.com/hako/branca"
)

// Messages
var notAuthErrorMessage = &flash.Message{
	Kind:    1,
	Content: "You have to be authenticated.",
}

// HasToBeAuth checks if the user is authenticated to go to the next http.HandlerFunc.
func HasToBeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("t")

		if err != nil {
			http.Redirect(w, r, "/login", 302)

			return
		}

		brca := branca.NewBranca(os.Getenv("BRANCA_SECRET"))
		payloadString, err := brca.DecodeToString(token.Value)

		if err != nil {
			flash.Add("/login", w, r, notAuthErrorMessage)

			return
		}

		var userPayload user.Payload

		json.Unmarshal([]byte(payloadString), &userPayload)

		newR := r.WithContext(user.ContextWithPayload(r.Context(), &userPayload))

		next(w, newR)
	}
}
