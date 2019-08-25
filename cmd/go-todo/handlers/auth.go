package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/internal/utils"
	"github.com/efreitasn/go-todo/pkg/flash"
	"github.com/hako/branca"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 * Messages
 */

// Login
var authLoginPOSTParsingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while parsing the request.",
}

var authLoginPOSTFetchingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while fetching the user from the db.",
}

var authLoginPostNoUserFoundErrorMsg = &flash.Message{
	Kind:    1,
	Content: "No user was found with the provided username.",
}

var authLoginPOSTEncodingErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while encoding user data.",
}

// Signup
var signupPOSTInsertErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while adding the user.",
}

var signupPOSTSuccessMsg = &flash.Message{
	Kind:    0,
	Content: "User created!",
}

// Auth is aa
type Auth struct {
	c *mongo.Collection
}

/*
 * Login
 */

// LoginGET is ...
func (a *Auth) LoginGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "login"

	utils.WriteTemplates(w, tData, "login")
}

// LoginPOST is ....
func (a *Auth) LoginPOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.ParseForm()

	if err != nil {
		flash.Add(
			"/login",
			w,
			r,
			authLoginPOSTParsingErrorMsg,
			true,
		)

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
		flash.Add(
			"/login",
			w,
			r,
			authLoginPostNoUserFoundErrorMsg,
			true,
		)

		return
	}

	if findRes.Err() != nil {
		flash.Add(
			"/login",
			w,
			r,
			authLoginPOSTFetchingErrorMsg,
			true,
		)

		return
	}

	var user user.User

	err = findRes.Decode(&user)

	if err != nil {
		flash.Add(
			"/login",
			w,
			r,
			authLoginPOSTFetchingErrorMsg,
			true,
		)

		return
	}

	if match := user.ComparePw(r.Form.Get("password")); !match {
		flash.Add(
			"/login",
			w,
			r,
			&flash.Message{
				Kind:    1,
				Content: "Wrong password.",
			},
			true,
		)

		return
	}

	stringPayload, err := json.Marshal(user.Payload())

	if err != nil {
		flash.Add(
			"/login",
			w,
			r,
			authLoginPOSTEncodingErrorMsg,
			true,
		)

		return
	}

	brca := branca.NewBranca(os.Getenv("BRANCA_SECRET"))
	token, err := brca.EncodeToString(string(stringPayload))

	http.SetCookie(w, &http.Cookie{
		Name:     "t",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(time.Hour * 24),
	})

	http.Redirect(w, r, "/list", 303)
}

/*
 * Signup
 */

// SignupGET is ...
func (a *Auth) SignupGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "signup"

	utils.WriteTemplates(w, tData, "signup")
}

// SignupPOST is ...
func (a *Auth) SignupPOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	userToInsert := user.InsertUser{
		Name:      r.Form.Get("name"),
		Username:  r.Form.Get("username"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
		CreatedAt: time.Now(),
	}

	userToInsert.HashPw()

	_, err := a.c.InsertOne(
		ctx,
		userToInsert,
	)

	if err != nil {
		flash.Add(
			"/signup",
			w,
			r,
			signupPOSTInsertErrorMsg,
			true,
		)

		return
	}

	flash.Add(
		"/login",
		w,
		r,
		signupPOSTSuccessMsg,
		true,
	)
}

/*
 * Middleware
 */

// SetUpAuth is ...
func SetUpAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("t")

		if err != nil {
			http.Redirect(w, r, "/login", 302)

			return
		}

		brca := branca.NewBranca(os.Getenv("BRANCA_SECRET"))
		payloadString, err := brca.DecodeToString(token.Value)

		if err != nil {
			flash.Add(
				"/login",
				w,
				r,
				&flash.Message{
					Kind:    1,
					Content: "You have to be authenticated.",
				},
				true,
			)

			return
		}

		var userPayload user.Payload

		json.Unmarshal([]byte(payloadString), &userPayload)

		newR := r.WithContext(userPayload.AddToContext(r.Context()))

		next(w, newR)
	}
}
