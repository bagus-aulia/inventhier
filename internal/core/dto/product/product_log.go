package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductLog is a struct for mongodb product log when its in and out
type ProductLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UUID         string             `bson:"uuid"`
	Product      Product            `bson:"product"`
	Status       string             `bson:"status"`
	Quantity     int                `bson:"quantity"`
	CurrentStock int                `bson:"current_stock"`
	CreatedAt    time.Time          `bson:"created_at"`
}
