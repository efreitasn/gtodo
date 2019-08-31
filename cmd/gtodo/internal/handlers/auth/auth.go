package auth

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Auth is the representation of all the auth-related handlers.
type Auth struct {
	c *mongo.Collection
}

// New creates an Auth struct.
func New(c *mongo.Collection) *Auth {
	return &Auth{c}
}
