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

//SESSIONS

type SessionManager struct {
	Sessions map[string]string
}

func NewSessions() *SessionManager {
	s := new(SessionManager)
	s.Sessions = make(map[string]string)
	return s
}

func (s *SessionManager) Add(userName string, sessionID string) {
	if _, ok := s.Sessions[sessionID]; !ok {
		s.Sessions[sessionID] = userName
	}
}

func (s *SessionManager) GetUser(sessionID string) (User string) {
	return s.Sessions[sessionID]
}

func (s *SessionManager) DeleteSession(sessionID string) {
	if _, ok := s.Sessions[sessionID]; ok {
		delete(s.Sessions, sessionID)
	}
}

//SESSIONS END

var manager = NewSessions()

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	//http.HandleFunc("/logoff", notFound)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("sessionid")
	var html string
	if err == nil {
		html = fmt.Sprintf(indexPage, "<a href='/logoff'> Logoff "+manager.GetUser(cookie.Value)+"</a>")
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

/*func logoff(w http.ResponseWriter, r *http.Request) {

}*/

func sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
