package db

import (
	pages "Pages"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const connection = "mongodb+srv://ImBlue:abubakar2ahmed@charactersheetdb.zre7u.mongodb.net/CharacterSheetDB?retryWrites=true&w=majority"

var (
	ctx    context.Context
	cancel context.CancelFunc
)

//User represents a user in the database
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Sheets   []string `json:"sheets"`
}

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
		ok = false
	}
	pwd := []byte(pass)
	dpwd := []byte(result.Password)
	passErr := bcrypt.CompareHashAndPassword(dpwd, pwd)
	if passErr != nil {
		ok = false
	} else {
		ok = true
	}
	if (user == result.Username) && (pass == result.Password) {
		ok = true
	}
	return ok
}

//CheckUserName checks the database if the name can be found
func CheckUserName(user string) bool {
	ok := false
	var result struct {
		Username string `json:"username"`
	}
	filter := bson.M{"username": user}
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
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		ok = false
	}
	if user == result.Username {
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

//RegisterUser registers a new user in the database
func RegisterUser(user User) error {
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
	collection := client.Database("CharacterSheets").Collection("users")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, user)
	return err
}

//RegisterSheet registers a new character sheet with a user
func RegisterSheet(user string, sheet pages.Sheet) bool {
	var ok = false
	filterUser := bson.M{"username": user}
	update := bson.M{"$push": bson.M{"sheets": sheet.Name}}
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
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, update)
	if err != nil {
		ok = false
	} else {
		ok = true
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, sheet)
	if err != nil {
		ok = false
	}
	return ok
}

//DeleteSheet deletes a sheet from the database
func DeleteSheet(user string, sheet string) {
	filterUser := bson.M{"username": user}
	filterSheets := bson.M{"owner": user, "name": sheet}
	updateUser := bson.M{"$pull": bson.M{"sheets": sheet}}
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
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, updateUser)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	_, err = collection.DeleteOne(ctx, filterSheets)
}
