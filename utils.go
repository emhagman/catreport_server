package main

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("CATREPORT_SESSION_SALT")))

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fn(w, r)
	}
}

func HashPassword(password, salt []byte) string {
	pass := pbkdf2.Key(password, salt, 4096, sha256.Size, sha256.New)
	return fmt.Sprintf("%x", pass)
}

func OldHashPassowrd(password string) string {
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SessionLogin(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	session.Values["access"] = true
	session.Save(req, res)
}

func SessionLogout(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	session.Values["access"] = false
	session.Save(req, res)
}
