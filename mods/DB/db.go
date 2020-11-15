package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection = "mongodb+srv://ImBlue:abubakat2ahmed@charactersheetdb.zre7u.mongodb.net/CharacterSheetDB?retryWrites=true&w=majority"

//GetDBConnection connects to the database and returns the connection
func GetDBConnection() (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	return client, ctx, err
}
