package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func AddUser(db *sql.DB, username string, password string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(db *sql.DB, username string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

// create a test user to login on the forum and test the app
func SeedUser(db *sql.DB) error {
	err := AddUser(db, "test", "test")
	if err != nil {
		return err
	}

	return nil
}

func AddCategory(db *sql.DB, title string, content string) error {
	_, err := db.Exec("INSERT INTO categories (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		return err
	}

	return nil
}

func InitApp() error {
	db, err := InitDb()
	if err != nil {
		return err
	}
	defer db.Close()

	err = CreateUserTable(db)
	if err != nil {
		return err
	}

	err = SeedUser(db)
	if err != nil {
		return err
	}

	fmt.Println("Database initialized successfully!")
	return nil
}
