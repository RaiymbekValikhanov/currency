package store

import (
	"currency/internal/models"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store  struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{ 
		db: db,
	}
}

func (s Store) CleanUp() (error) {
	if _, err := s.db.Exec("TRUNCATE currency"); err != nil {
		return err
	}

	return nil
}

func (s Store) GetCurrencyList() ([]models.Currency, error) {
	rows, err := s.db.Query("SELECT code, name, rate FROM currency")
	if err != nil {
		return nil, err
	}

	currs := make([]models.Currency, 0)
	for rows.Next() {
		var c models.Currency
		if err := rows.Scan(&c.Code, &c.Name, &c.Rate); err != nil {
			return nil, err
		}

		currs = append(currs, c)
	}

	return currs, nil
}

func (s Store) GetCurrency(code string) (*models.Currency, error) {
	c := &models.Currency{}

	if err := s.db.QueryRow(
		"SELECT code, name, rate FROM currency WHERE code=$1", 
		code,
	).Scan(&c.Code, &c.Name, &c.Rate); err != nil {
		return nil, err
	}

	return c, nil
}

func (s Store) RemoveCurrency(code string) (error) {
	if _, err := s.db.Exec(
		"DELETE FROM currency WHERE code=$1", 
		code,
	); err != nil {
		return err
	}

	return nil
}

func (s Store) SaveCurrenciesData(data models.NBKData) (error) {
	for _, c := range data.Currencies {
		if _, err := s.db.Exec(
			"INSERT INTO currency (code, name, rate) VALUES ($1, $2, $3) RETURNING id",
			c.Code, c.Name, c.Rate,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s Store) UpdateCurrenciesData(data models.NBKData) (error) {
	for _, c := range data.Currencies {
		if _, err := s.db.Exec(
			"UPDATE currency SET rate=$1 WHERE code=$2",
			c.Rate, c.Code,
		); err != nil {
			return err
		}
	}

	return nil
}