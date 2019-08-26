package user

import (
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

// ComparePw compares a plain text password with a hashed one to check if there's a match.
func (u *User) ComparePw(plainTextPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPw))

	if err != nil {
		return false
	}

	return true
}
