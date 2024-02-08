package main

import (
	"custom/ts"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// HomeHandler will be the http handler for the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.Method + ": On Home")
	homeBoard.Current = TypeUser{Current: current, LoggedIn: loggedIn} // create the current user for the board
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, homeBoard) // execute the template
}

// LoginHandler will all logins
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ": On Sign Up")
	http.ServeFile(w, r, "templates/login.html") //
}

// SignInHandler will handle all sign-ins
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ": On Sign in")
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		r.ParseForm()
		user, err := UDB.GetUser(
			strings.Join(r.Form["email"], ""),
			strings.Join(r.Form["password"], ""),
		)
		if err != nil {
			log.Fatal(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		current = user
		loggedIn = true
		current.Handle()
		http.Redirect(w, r, "/"+current.Username+"-"+current.WebID, http.StatusSeeOther)
	}
}

// SignUpHandler will handle all sign-ups
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ": On Sign Up")
	http.ServeFile(w, r, "templates/signup.html") //
}

// CreateUserHandler handles each newly create user from the signup handler
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ": On Create")
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		r.ParseForm()
		log.Println(r.Form["password"])
		log.Println(r.Form["cpass"])
		created, err := NewUser(
			strings.Join(r.Form["username"], ""),
			strings.Join(r.Form["email"], ""),
			strings.Join(r.Form["password"], ""),
			strings.Join(r.Form["cpass"], ""),
		)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		} else {
			current = created
			loggedIn = true
			current.Handle()
			UDB.InsertUser(current)
			http.Redirect(w, r, "/"+current.Username+"-"+current.WebID, http.StatusSeeOther)
		}
	}
}

// LogoutHandler will handle each time the user log's out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	current = CurrentUser{}
	loggedIn = false
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// UploadHandler will handle all file uploads to the website
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		r.ParseMultipartForm(32 << 20)                 // Parse the upload form
		file, handler, err := r.FormFile("uploadfile") // get the file that was uploaded
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		defer file.Close()                                                                         // close the file when it is no longer used
		fileID := ts.GenerateID(24)                                                                // create a file ID
		suffix := ts.GetSuffix(handler.Filename)                                                   // get the file's suffix
		filename := ts.ReplaceWhitespace(ts.RemoveSuffix(handler.Filename)+"-"+fileID+suffix, '_') // make the filename
		location := "private/img/" + filename                                                      // set the location
		f, err := os.OpenFile(location, os.O_WRONLY|os.O_CREATE, 0666)                             // open the created file, WRITE | CREATE
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		defer f.Close()                                 // close said file, when no longer needed
		io.Copy(f, file)                                // copy the file into the server's storage
		current.Files = append(current.Files, filename) // add the file to the current user's file array
		err = UDB.InsertFile(current, filename)         // insert the file into the current user's file array (WITHIN THE DATABASE)
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		current.HandleFile(filename)
		http.Redirect(w, r, "/"+current.Username+"-"+current.WebID, http.StatusSeeOther) // Redirect the user back to their user-page
	}
}

// FileHandler will handle file deletions for each file
type FileHandler struct {
	Filename string
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/"+current.Username+"-"+current.WebID, http.StatusSeeOther)
	} else {
		var e error
		current.Files, e = current.RemoveFile(h.Filename)
		if e != nil {
			log.Fatal(e)
			panic(1)
		}
		err := UDB.DeleteFile(current, h.Filename)
		if err != nil {
			log.Fatal(err)
			panic(1)
		}
		ts.Execute("rm private/img/" + h.Filename)
		http.Redirect(w, r, "/"+current.Username+"-"+current.WebID, http.StatusSeeOther)
	}
}
