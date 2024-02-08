package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// URI is the database URI access
const URI = `mongodb+srv://McAssly:5zIXhG8756789ibj@testing-lugq4.gcp.mongodb.net/test?retryWrites=true&w=majority`

// UserDB is the user database structure
type UserDB struct {
	Collection *mongo.Collection
	DB         *mongo.Database
	Client     *mongo.Client
	Context    context.Context
}

// Connect will simply connect the user database client to the user database
func (DB *UserDB) Connect() {
	log.Println("Connecting to User Database")
	DB.Context = context.Background()
	log.Println("Creating new mognodb client")
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	DB.Client = client
	log.Println("Connecting to the mongodb client")
	err = DB.Client.Connect(DB.Context)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	log.Println("Pinging the mongodb client to check connection")
	err = DB.Client.Ping(DB.Context, nil)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	DB.DB = DB.Client.Database("user-database")
	DB.Collection = DB.DB.Collection("users")
	log.Println("Finished connecting to the mongodb client\nUser Database Connection Secured")
}

// UnhandleAll will unhandle all the users that are in the user database
func (DB *UserDB) UnhandleAll() {
	result, err := DB.Collection.UpdateMany(DB.Context, bson.M{}, bson.D{{"$set", bson.D{{"handled", false}}}})
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	log.Println(result.UpsertedID)
}

// InsertUser will append a user into the User Database collection
func (DB *UserDB) InsertUser(user CurrentUser) {
	log.Println("Attempting to insert user into User Database")
	result, err := DB.Collection.InsertOne(DB.Context, user)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	log.Printf("The resulting ID: %v\n", result.InsertedID)
}

// GetUser will obtain a user from the collection based on the given [email & password]
func (DB *UserDB) GetUser(email, password string) (CurrentUser, error) {
	log.Println("Attempting to obtain a user from the UDB")
	var result CurrentUser
	err := DB.Collection.FindOne(DB.Context, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return CurrentUser{}, err
	}
	err = CheckPassword(password, result.Password)
	if err != nil {
		return CurrentUser{}, err
	}
	log.Println("Obtained User")
	return result, nil
}

// SetHandled will set whether the user is handled
func (DB *UserDB) SetHandled(User *CurrentUser, to bool) {
	result, err := DB.Collection.UpdateOne(DB.Context, bson.M{"_id": User.ID}, bson.D{{"$set", bson.D{{"handled", to}}}})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.UpsertedID)
}

// Disconnect will disconnect from the user database
func (DB *UserDB) Disconnect() {
	err := DB.Client.Disconnect(DB.Context)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
}

// InsertFile will insert a file into the give User's file array
func (DB *UserDB) InsertFile(user CurrentUser, filename string) error {
	result, err := DB.Collection.UpdateOne(
		DB.Context,
		bson.M{"_id": user.ID},
		bson.D{{"$push", bson.D{{"files", filename}}}},
	)
	if err != nil {
		return err
	}
	log.Println(result.UpsertedID)
	return nil
}

// DeleteFile will delete a file from the given User's file array
func (DB *UserDB) DeleteFile(user CurrentUser, filename string) error {
	result, err := DB.Collection.UpdateOne(
		DB.Context,
		bson.M{"_id": user.ID},
		bson.D{{"$pull", bson.D{{"files", filename}}}},
	)
	if err != nil {
		return err
	}
	log.Println(result.UpsertedID)
	return nil
}
