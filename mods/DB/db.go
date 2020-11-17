package db

import (
	pages "Pages"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection = "mongodb+srv://ImBlue:abubakar2ahmed@charactersheetdb.zre7u.mongodb.net/CharacterSheetDB?retryWrites=true&w=majority"

var (
	ctx    context.Context
	cancel context.CancelFunc
)

//CheckUser checks if a user exists in the database
func CheckUser(user string, pass string) bool {
	var result struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ok := false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	filter := bson.M{"username": user}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	if (user == result.Username) && (pass == result.Password) {
		ok = true
	}
	return ok
}

//GetSheets gets the characyer sheets from the database
func GetSheets(user string) []string {
	var result struct {
		Sheets []string `json:"sheets"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	filter := bson.M{"username": user}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result.Sheets
}

//GetSheet retrieves character sheet data from the database and fills an object with the data
func GetSheet(username string, sheetname string) pages.Sheet {
	sheet := pages.Sheet{}
	filter := bson.M{"owner": username, "name": sheetname}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("sheets")
	err = collection.FindOne(context.TODO(), filter).Decode(&sheet)
	return sheet
}
