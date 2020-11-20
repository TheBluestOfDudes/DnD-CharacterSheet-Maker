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

//User represents a user in the database
type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Sheets   []string `json:"sheets"`
}

//CheckUser checks if a user exists in the database
func CheckUser(user string, pass string) (bool, error) {
	ok := false         //Boolean that determines if the match was valid
	var result struct { //Help struct that saves the data retrieved from the database
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	filter := bson.M{"username": user} //Query filter we use to fetch from the database
	defer func() {                     //We disconnect the client after the method is done
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End method if we failed to connect
		return ok, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //We get the mongodb collection
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil { //End method if we failed to get collection
		return ok, err
	}
	pwd := []byte(pass)
	dpwd := []byte(result.Password)
	passErr := bcrypt.CompareHashAndPassword(dpwd, pwd) //Check if the hashed password in the database matches the given password
	if passErr != nil {
		return ok, passErr //End the method if the hashes didnt match
	}
	if user == result.Username { //If the username was correct we set our boolean to sucess
		ok = true
	}
	return ok, nil
}

//CheckUserName checks the database if the name can be found
func CheckUserName(user string) (bool, error) {
	ok := false         //Boolean that determines if the match was valid
	var result struct { //Help struct that saves the data retrieved from the database
		Username string `json:"username"`
	}
	filter := bson.M{"username": user}                                       //Query filter we use to fetch from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we failed to connect
		return false, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Getting collection and try to find the user
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil { //End if we failed. Error is set to nil, because we expect a failure
		return false, nil
	}
	if user == result.Username { //If we found a match, the ok is set to true
		ok = true
	}
	return ok, nil
}

//GetSheets gets the list of character sheet names from a given user in the database
func GetSheets(user string) ([]string, error) {
	var result struct { //Help struct that saves the data retrieved from the database
		Sheets []string `json:"sheets"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	filter := bson.M{"username": user} //Query filter we use to fetch from the database
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we failed to connect. Return empty array
		result.Sheets = []string{}
		return result.Sheets, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Getting collection and query for the sheets
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil { //End if we failed to connect. Return empty array
		result.Sheets = []string{}
		return result.Sheets, err
	}
	return result.Sheets, nil
}

//GetSheet retrieves a given user's character sheet from the database
func GetSheet(username string, sheetname string) (pages.Sheet, error) {
	sheet := pages.Sheet{}                                                   //Object that holds
	filter := bson.M{"owner": username, "name": sheetname}                   //Query filter we use to fetch from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we failed to connect.
		return sheet, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get the collection and retrieve the sheet data
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("sheets")
	err = collection.FindOne(context.TODO(), filter).Decode(&sheet)
	return sheet, err
}

//RegisterUser registers a new user in the database
func RegisterUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we fail
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get collection and try to insert the new user
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, user)
	return err
}

//RegisterSheet registers a new character sheet with a user
func RegisterSheet(user string, sheet pages.Sheet) error {
	filterUser := bson.M{"username": user}                                   //Query filter we use to select the user we want to update
	update := bson.M{"$push": bson.M{"sheets": sheet.Name}}                  //Update query for the given user's "sheets" array
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connect to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we fail
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get the user collection and try to update the user's "sheets" array
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, update)
	if err != nil { //End if we fail
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get the sheets collection and try to insert the new sheet
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, sheet)
	if err != nil { //End if we fail
		return err
	}
	return err
}

//DeleteSheet deletes a sheet from the database
func DeleteSheet(user string, sheet string) error {
	filterUser := bson.M{"username": user}                                   //Query filter to select the right user
	filterSheets := bson.M{"owner": user, "name": sheet}                     //Query filter to select the correct sheet
	updateUser := bson.M{"$pull": bson.M{"sheets": sheet}}                   //Update query to remove the sheet from the given user's "sheets" array
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Connecting to the database
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil { //End if we fail
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get the users collection and try to delete the sheet from the user's array
	defer cancel()
	collection := client.Database("CharacterSheets").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filterUser, updateUser)
	if err != nil { //End if we fail
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //Get the sheets collection and try to delete the sheet
	defer cancel()
	collection = client.Database("CharacterSheets").Collection("sheets")
	_, err = collection.DeleteOne(ctx, filterSheets)
	if err != nil { //End if we fail
		return err
	}
	return err
}
