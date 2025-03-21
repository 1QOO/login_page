package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", login)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/welcome", welcome)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))

	fmt.Println("Server runs at port 8080.")
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "admin" && password == "admin" {
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			return
		}
	}

	http.ServeFile(w, r, "../index.html")
}

func register(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/register.html")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/welcome.html")
}
