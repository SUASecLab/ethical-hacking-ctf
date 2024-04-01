package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(collection string) (context.Context, context.CancelFunc, *mongo.Client, *mongo.Collection, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoUsername+":"+
		mongoPassword+"@"+mongoHostname+":27017/"+mongoDBName))
	if err != nil {
		log.Println("Can not connect to database: ", err)
		return ctx, cancel, client, nil, false
	}

	return ctx, cancel, client, client.Database(mongoDBName).Collection(collection), true
}
