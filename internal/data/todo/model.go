package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo is an item of the list.
type Todo struct {
	ID        *primitive.ObjectID `bson:"_id"`
	Title     string              `bson:"title"`
	Done      bool                `bson:"done"`
	CreatedAt time.Time           `bson:"createdAt"`
}

// InsertTodo is an item to be inserted into the database.
type InsertTodo struct {
	UserID    *primitive.ObjectID `bson:"_user"`
	Title     string              `bson:"title"`
	Done      bool                `bson:"done"`
	CreatedAt time.Time           `bson:"createdAt"`
}
