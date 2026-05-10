package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	From      string             `bson:"from" json:"from"`
	To        string             `bson:"to" json:"to"`
	Note      string             `bson:"note" json:"note"`
	Status    string             `bson:"status" json:"status"` // pending, approved, declined
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
