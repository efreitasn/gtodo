package todo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FetchDone fetches the done todos.
func FetchDone(ctx context.Context, c *mongo.Collection, uID *primitive.ObjectID) (todos []Todo, err error) {
	cursor, err := c.Find(
		ctx,
		bson.D{
			{
				Key:   "_user",
				Value: uID,
			},
			{
				Key:   "done",
				Value: true,
			},
		},
		&options.FindOptions{
			Sort: bson.D{
				{
					Key:   "createdAt",
					Value: 1,
				},
			},
		},
	)

	if err != nil {
		return todos, err
	}

	err = cursor.All(
		ctx,
		&todos,
	)

	if err != nil {
		return todos, err
	}

	return todos, nil
}

// FetchNotDone fetches the not done todos.
func FetchNotDone(ctx context.Context, c *mongo.Collection, uID *primitive.ObjectID) (todos []Todo, err error) {
	cursor, err := c.Find(
		ctx,
		bson.D{
			{
				Key:   "_user",
				Value: uID,
			},
			{
				Key:   "done",
				Value: false,
			},
		},
		&options.FindOptions{
			Sort: bson.D{
				{
					Key:   "createdAt",
					Value: 1,
				},
			},
		},
	)

	if err != nil {
		return todos, err
	}

	err = cursor.All(
		context.Background(),
		&todos,
	)

	if err != nil {
		return todos, err
	}

	return todos, nil
}
