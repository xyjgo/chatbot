package mongo

import (
	"chatbot"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	DBNAME        = "chatbot"
	TABLE_REVIEWS = "reviews"
)

func SaveReviewRecord(ctx context.Context, review *chatbot.ReviewRecord) error {
	r := &Record{
		Id:     primitive.NewObjectID(),
		Review: review,
	}
	return Db.InsertOne(ctx, DBNAME, TABLE_REVIEWS, &r)
}

func EnsureIndex(ctx context.Context) {
	if err := Db.CreateIndex(ctx, DBNAME, TABLE_REVIEWS, bson.D{{"review.reviewid", 1}}, options.Index().SetUnique(true)); err != nil {
		log.Fatalln("ensure index err", err)
	}
}
