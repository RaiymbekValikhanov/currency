package api

import (
	"currency/internal/store"
	"database/sql"
)

var (
	DataBaseURL = "host=localhost dbname=halyk sslmode=disable"
)

func Start() error {
	db, err := newDB(DataBaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := store.NewStore(db)
	if err := store.CleanUp(); err != nil {
		return err
	}

	server := NewServer(store)
	return server.Run()
}

func newDB(dburl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return  nil, err
	}

	return db, nil
}