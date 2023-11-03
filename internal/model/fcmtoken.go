package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FCMToken struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      string             `bson:"user_id"`
	Application string             `bson:"application"`
	Fcm         string             `bson:"fcm"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
