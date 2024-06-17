package store

import "database/sql"

type Store struct {
	Db *sql.DB
}
