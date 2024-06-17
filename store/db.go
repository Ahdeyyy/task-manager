package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func CreateDatabase(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	statement := `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL, 
			email VARCHAR(255) NOT NULL UNIQUE
		);
	`
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stmt := `CREATE TABLE IF NOT EXISTS session_store (
		k  VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		v  BYTEA NOT NULL,
		e  BIGINT NOT NULL DEFAULT '0',
		u  VARCHAR(100)
	);`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func initializeTables(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL, 
			email VARCHAR(255) NOT NULL UNIQUE
		);
	`

	stmt := `CREATE TABLE IF NOT EXISTS session_store (
		k  VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		v  BYTEA NOT NULL,
		e  BIGINT NOT NULL DEFAULT '0',
		u  VARCHAR(100)
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	_, err = db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func CreateStore(connectionUri string) (*Store, error) {
	s := &Store{}
	db, err := sql.Open("postgres", connectionUri)
	if err != nil {
		return nil, err
	}
	s.Db = db
	err = initializeTables(db)

	if err != nil {
		return nil, err
	}

	return s, nil
}
