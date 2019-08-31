package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/data/user"
	"github.com/efreitasn/gtodo/internal/utils"
	"github.com/efreitasn/gtodo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Messages
var loginPOSTParsingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while parsing the request.",
}

var loginPOSTFetchingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while fetching the user from the db.",
}

var loginPOSTNoUserFoundErrorMsg = &flash.Message{
	Kind:    1,
	Content: "No user was found with the provided username.",
}

var loginPOSTWrongPasswordErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Wrong password.",
}

var loginPOSTEncodingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while encoding user data.",
}

// LoginGET renders the login form.
func (a *Auth) LoginGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "login"
	tData.Title = "Login"

	utils.WriteTemplates(w, tData, "login")
}

// LoginPOST logs an user in.
func (a *Auth) LoginPOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.ParseForm()

	if err != nil {
		flash.Add("/login", w, r, loginPOSTParsingErrorMsg)

		return
	}

	findRes := a.c.FindOne(
		ctx,
		bson.D{
			{
				Key:   "username",
				Value: r.Form.Get("username"),
			},
		},
	)

	if findRes.Err() == mongo.ErrNoDocuments {
		flash.Add("/login", w, r, loginPOSTNoUserFoundErrorMsg)

		return
	}

	if findRes.Err() != nil {
		flash.Add("/login", w, r, loginPOSTFetchingErrorMsg)

		return
	}

	var u user.User

	err = findRes.Decode(&u)

	if err != nil {
		flash.Add("/login", w, r, loginPOSTFetchingErrorMsg)

		return
	}

	if match := u.ComparePw(r.Form.Get("password")); !match {
		flash.Add("/login", w, r, loginPOSTWrongPasswordErrorMsg)

		return
	}

	userPayload := user.PayloadFromUser(&u)
	err = user.ResponseWithTokenCookie(w, userPayload)

	if err != nil {
		flash.Add("/login", w, r, loginPOSTEncodingErrorMsg)

		return
	}

	http.Redirect(w, r, "/list", 303)
}
