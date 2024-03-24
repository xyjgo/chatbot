package mongo

import (
	"chatbot"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	Id     primitive.ObjectID `bson:"_id"`
	Review *chatbot.ReviewRecord
}
