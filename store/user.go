package store

import (
	"database/sql"
	"log"
	"task-manager/models"
)

func GetUsers(db *sql.DB) []models.User {
	// Get all users

	statement := `SELECT * FROM users`

	rows, err := db.Query(statement)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}

	return users
}

func GetUser(db *sql.DB, username string) models.User {
	// Get user by username

	statement := `SELECT * FROM users WHERE username = $1`

	row := db.QueryRow(statement, username)

	user := models.User{}

	if err := row.Scan(&user.Username, &user.Password, &user.Email); err != nil {
		log.Println(err)
	}

	return user
}
