package database

import (
	"context"
	"fmt"
	"main/package/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnectionMongoDB() *mongo.Collection {
	uri := util.Getenv("MONGODB_URI", "kepo")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	fmt.Println("cooneect")
	dbname := util.Getenv("MONGODB_DB", "db")
	collection := util.Getenv("MONGODB_COLLECTION", "coll")

	return client.Database(dbname).Collection(collection)
}
