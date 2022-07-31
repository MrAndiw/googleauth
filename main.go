package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"log"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {

	key := "login-google-go" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30     // 30 days
	isProd := false          // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("1088820804908-46qpbgsqolb8v62ksd4tcuc4psodeknn.apps.googleusercontent.com", "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", "http://127.0.0.1:3000/auth/google/callback", "email", "profile"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	port := os.Getenv("PORT")
	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, p))
}
