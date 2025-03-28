package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func checkUsername(username string) (bool, error) {
	db, err := sql.Open("mysql", "iqro:iqroDB@tcp(localhost:3306)/WebProjects")
	if err != nil {
		fmt.Println("Can't connect to the database: %w", err)
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Login_page_Accounts WHERE username=?", username).Scan(&count)
	if err != nil {
		fmt.Println("failed to execute query: %w", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func registerUser(username, password string) string {
	db, err := sql.Open("mysql", "iqro:iqroDB@tcp(localhost:3306)/WebProjects")
	if err != nil {
		fmt.Println("Can't connect to the database: %w", err)
		return ""
	}
	defer db.Close()

	sessionToken := generateToken(32)
	_, err = db.Exec("INSERT INTO Login_page_Accounts (username, password, session_token) VALUES(?, ?, ?)", username, password, sessionToken)
	if err != nil {
		fmt.Println("Cannot insert new data to the table: %w", err)
		return ""
	}

	return sessionToken
}

func CheckSession(userSession string) (bool, error) {
	db, err := sql.Open("mysql", "iqro:iqroDB@tcp(localhost:3306)/WebProjects")
	if err != nil {
		fmt.Println("Error, failed to open database", http.StatusInternalServerError)
		return false, errors.New("failed to connect to the database")
	}
	defer db.Close()

	var count int
	data := db.QueryRow("SELECT COUNT(*) FROM Login_page_Accounts WHERE session_token=?", userSession).Scan(&count)
	if data != nil {
		fmt.Println("failed to execute query: %w", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
