package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/template"
	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/internal/utils"
	"github.com/efreitasn/go-todo/pkg/flash"
)

// Messages
var signupPOSTInsertErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while adding the user.",
}

var signupPOSTSuccessMsg = &flash.Message{
	Kind:    0,
	Content: "User created!",
}

// SignupGET is ...
func (a *Auth) SignupGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
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
		flash.Add("/signup", w, r, signupPOSTInsertErrorMsg)

		return
	}

	flash.Add("/login", w, r, signupPOSTSuccessMsg)
}
