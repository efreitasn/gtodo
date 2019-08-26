package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser is an user to be inserted in the database.
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
