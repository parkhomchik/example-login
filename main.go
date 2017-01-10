package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const indexPage = `<doctype html>
<html>
    <head>
        <title>Example Login</title>
    </head>
    <body>
        <h1>Example Login</h1>
		%s
    </body>
</html>`

const loginBlock = `<form action="/login" method="post">
    Username:<input type="text" name="username">
    Password:<input type="password" name="password">
    <input type="submit" value="Login">
</form>`

const logoffBlock = `<form action="/logoff" method="post">
    <input type="hidden" name="csrf" value="">    
	<h1>%s</h1>
    <input type="submit" value="Logoff">
</form>`

//SessionManager -
type SessionManager struct {
	Sessions map[string]string
}

//NewSessions -
func NewSessions() *SessionManager {
	s := new(SessionManager)
	s.Sessions = make(map[string]string)
	return s
}

//Add -
func (s *SessionManager) Add(userName string, sessionID string) {
	if _, ok := s.Sessions[sessionID]; !ok {
		s.Sessions[sessionID] = userName
	}
}

//GetUser -
func (s *SessionManager) GetUser(sessionID string) (User string) {
	return s.Sessions[sessionID]
}

//DeleteSession -
func (s *SessionManager) DeleteSession(sessionID string) {
	if _, ok := s.Sessions[sessionID]; ok {
		delete(s.Sessions, sessionID)
	}
}

//CheckSession -
func (s *SessionManager) CheckSession(sessionID string) bool {
	_, ok := s.Sessions[sessionID]
	return ok
}

//SESSIONS END

var manager = NewSessions()

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logoff", logoff)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	cookie, _ := req.Cookie("sessionid")
	var html string
	if manager.CheckSession(cookie.Value) {
		html = fmt.Sprintf(indexPage, fmt.Sprintf(logoffBlock, manager.GetUser(cookie.Value)))
	} else {
		html = fmt.Sprintf(indexPage, loginBlock)
	}

	fmt.Println("index page")
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		html,
	)
}

func login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		http.Redirect(w, r, "/", 200)
	} else {
		r.ParseForm()

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		fmt.Println(len(username))

		if (username == "") || (password == "") {

			http.Redirect(w, r, "/", 301)
			return
		}

		fmt.Println("username:", username, "password:", password)
		sessionID := sessionID()
		manager.Add(username, sessionID)
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "sessionid", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", 301)
		return
	}

}

func logoff(w http.ResponseWriter, r *http.Request) {

	fmt.Println("logoff")
	cookie, _ := r.Cookie("sessionid")
	fmt.Println("session = ", cookie)
	if manager.CheckSession(cookie.Value) {
		manager.DeleteSession(cookie.Value)
	}

	fmt.Println("Sessions = ", manager.Sessions)
	http.Redirect(w, r, "/", 301)
}

func sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
