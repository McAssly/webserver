package main

import (
	"custom/ts"
	"errors"
	"html/template"
	"log"
	"net/http"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var current CurrentUser
var loggedIn bool

// CurrentUser is the current user structure
type CurrentUser struct {
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password []byte             `bson:"password"`
	ID       primitive.ObjectID `bson:"_id"`
	WebID    string             `bson:"webid"`
	Files    []string           `bson:"files"`
	Handled  bool               `bson:"handled"`
}

// TypeUser will determine whether the user is logged in or not
type TypeUser struct {
	LoggedIn bool
	Current  CurrentUser
}

// NewUser will create a new user after the user signs-up for an account
func NewUser(name, email, pass, cpass string) (CurrentUser, error) {
	if pass != cpass {
		log.Printf("%v != %v\n", pass, cpass)
		return CurrentUser{}, errors.New("Passwords did not match")
	}
	return CurrentUser{Username: name, Email: email, Password: HashPassword(pass), ID: primitive.NewObjectID(), Files: []string{}, Handled: false, WebID: ts.GenerateID(24)}, nil
}

// HandleFileDeletors will do as it says
func (C *CurrentUser) HandleFileDeletors() {
	for _, b := range C.Files {
		C.HandleFile(b)
	}
}

// HandleFile will handle a file deletor for the given file
func (C *CurrentUser) HandleFile(filename string) {
	fh := new(FileHandler)
	fh.Filename = filename
	http.Handle("/r="+filename, fh)
}

// Handle will create a user handler for the current user
func (C *CurrentUser) Handle() {
	if !C.Handled {
		http.HandleFunc("/"+C.Username+"-"+C.WebID, func(w http.ResponseWriter, r *http.Request) {
			log.Println("Handling User")
			t, err := template.ParseFiles("templates/user.html")
			if err != nil {
				log.Fatal(err)
				panic(1)
			}
			t.Execute(w, C) // execute the template
		})
		C.Handled = true
		UDB.SetHandled(C, true)
		C.HandleFileDeletors()
	}
}

// RemoveFile will simple remove a file from the current user's file array
func (C *CurrentUser) RemoveFile(filename string) ([]string, error) {
	result := []string{}
	for _, b := range C.Files {
		if b != filename {
			result = append(result, b)
		}
	}
	if reflect.DeepEqual(result, C.Files) {
		return result, errors.New("Nothing was removed")
	}
	return result, nil
}
