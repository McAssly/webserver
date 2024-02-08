package main

import (
	"custom/ts"
	"log"
	"net/http"
)

// UDB is the user database for the site
var UDB = new(UserDB)

func main() {
	log.Println("Starting Server on :9090")
	UDB.Connect()          // Connect to the user database
	UDB.UnhandleAll()      // Remember that all users are no longer handled once the server has been restarted
	homeBoard = new(Board) // make a new board
	homeBoard.Load()       // load the home board
	http.Handle("/", new(ts.UHandler))
	http.HandleFunc("/home", HomeHandler) // handle the home page
	http.HandleFunc("/sign-up", SignUpHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/sign-in", SignInHandler)
	http.HandleFunc("/create", CreateUserHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/send", BoardHandler)
	err := http.ListenAndServe(":9090", nil) // start the server
	if err != nil {                          // if for some reason the server cannot be started, log what error that may appear
		log.Fatal(err)
	}
	defer UDB.Disconnect() // disconnect from the user database if need be
}
