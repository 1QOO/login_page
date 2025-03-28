package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "../index.html") })
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/register/success", registerSuccess)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/protected", protected)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))

	fmt.Print("Server runs at localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := http.StatusMethodNotAllowed
		http.Error(w, "HTTP method is not allowed.", err)
	}

	http.ServeFile(w, r, "../static/register.html")
}

// SIGNUP LOGIC
func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Method is not allowed: ", err)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		err := http.StatusBadRequest
		http.Error(w, "Error, failed to parse form:", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmedPassword := r.FormValue("confirmed-password")
	type res struct {
		Status  string
		Message string
	}

	if username != "" {
		userExist, err := checkUsername(username)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if userExist {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res{Status: "UsernameIsTaken", Message: "Username already taken"})
			return
		}
	}

	if password != confirmedPassword {
		err := http.StatusBadRequest
		fmt.Println("Error, password didn't match:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res{Status: "PasswordDidNotMatch", Message: "Password didn't match"})
		return
	}

	if len(password) < 8 {
		err := http.StatusBadRequest
		fmt.Println("Error, password must contain 8 or more characters:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res{Status: "PasswordIsTooShort", Message: "Password must contain 8 or more characters"})
		return
	}

	HashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error, failed to hash password:", err)
		return
	}

	sessionToken := registerUser(username, HashedPassword)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res{Status: "Registered", Message: "/register/success"})
}

func registerSuccess(w http.ResponseWriter, r *http.Request) {
	userSession, err := r.Cookie("session_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Redirect(w, r, "../index.html", http.StatusSeeOther)
			return

		default:
			log.Println(err)
			http.Error(w, "Server error.", http.StatusInternalServerError)
			return
		}
	}
	registrationSucces, err := CheckSession(userSession.Value)
	if err != nil {
		fmt.Println("Something is wrong", err)
		return
	}

	if !registrationSucces {
		http.Redirect(w, r, "../index.html", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "../static/registered.html")

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Method is not allowed.", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("passwd")

	user, ok := users[username]
	if !ok || checkPasswordHash(password, user.HashedPassword) {
		err := http.StatusUnauthorized
		http.Error(w, "username or password wrong.", err)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: false,
	})

	user.SessionToken = sessionToken
	users[username] = user

	fmt.Fprintf(w, "Login successfuly.")
}

func protected(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodPost {
		err := http.StatusUnauthorized
		http.Error(w, "Invalid request method.", err)
		return
	}*/

	if err := Authorize(r); err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "Unauthorized.", er)
		return
	}

	//	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation successful! Welcome ")
}
