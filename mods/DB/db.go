package db

import (
	pages "Pages"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//Conncetion link to the mongodb database.
var connection = os.Getenv("CONNECTION")

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
func CheckUser(user string, pass string) (bool, error) {
	ok := false
	var result struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
		return ok, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return ok, err
	}
	pwd := []byte(pass)
	dpwd := []byte(result.Password)
	passErr := bcrypt.CompareHashAndPassword(dpwd, pwd)
	if passErr != nil {
		return ok, passErr
	}
	if user == result.Username {
		ok = true
	}
	return ok, nil
}

//CheckUserName checks the database if the name can be found
func CheckUserName(user string) (bool, error) {
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
		return false, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false, nil
	}
	if user == result.Username {
		ok = true
	}
	return ok, nil
}

//GetSheets gets the characyer sheets from the database
func GetSheets(user string) ([]string, error) {
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
		result.Sheets = []string{}
		return result.Sheets, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		result.Sheets = []string{}
		return result.Sheets, err
	}
	return result.Sheets, nil
}

//GetSheet retrieves character sheet data from the database and fills an object with the data
func GetSheet(username string, sheetname string) (pages.Sheet, error) {
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
		return sheet, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("sheets")
	err = collection.FindOne(context.TODO(), filter).Decode(&sheet)
	return sheet, err
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
		return err
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
func RegisterSheet(user string, sheet pages.Sheet) error {
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
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, update)
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, sheet)
	if err != nil {
		return err
	}
	return err
}

//DeleteSheet deletes a sheet from the database
func DeleteSheet(user string, sheet string) error {
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
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, updateUser)
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	_, err = collection.DeleteOne(ctx, filterSheets)
	if err != nil {
		return err
	}
	return err
}
