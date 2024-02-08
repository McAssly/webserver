package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var homeBoard *Board

// Board is a message board structure
type Board struct {
	Current      TypeUser
	Incoming     []Message
	Messages     string
	HTMLMessages template.HTML
}

// Message is a structure that contains who sent it, and what they sent
type Message struct {
	Username string
	Message  string
}

// SendMessage will send in a message into the message board
func (B *Board) SendMessage(msg string, usrname string) {
	B.Incoming = append(B.Incoming, Message{Message: msg, Username: usrname})
}

// ParseIncomingMessages will parse the incoming messages into a string format
func (B *Board) ParseIncomingMessages() {
	result := []string{}
	for _, m := range B.Incoming {
		result = append(result, m.Parse())
	}
	B.Messages = B.Messages + strings.Join(result, "")
	B.Incoming = []Message{}
}

// Save will save the messaging history
func (B *Board) Save() {
	// log.Println(B.Messages)
	msgs, err := Encrypt(B.Messages)
	err = ioutil.WriteFile("public/messages/home.txt", msgs, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// Load will load the messaging history
func (B *Board) Load() {
	hash, err := ioutil.ReadFile("public/messages/home.txt")
	if len(hash) == 0 {
		return
	}
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	msgs, e := Decrypt(hash)
	if e != nil {
		log.Fatal(e)
		panic(1)
	}
	B.Messages = msgs
	B.Parse()
}

// Parse will parse the messages from string, to an html format onto the webpage
func (B *Board) Parse() {
	B.HTMLMessages = template.HTML(B.Messages)
}

// Parse will parse the message into a readable format
func (M *Message) Parse() string {
	return "<p class='message'><b class='usrname'>" + M.Username + "</b>: " + M.Message + "</p>"
}

// BoardHandler will handle the HTML link required for the board
func BoardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		r.ParseForm()
		message := strings.Join(r.Form["message"], "")
		homeBoard.SendMessage(message, current.Username)
		homeBoard.ParseIncomingMessages()
		homeBoard.Parse()
		homeBoard.Save()
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
