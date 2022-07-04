package store

import "currency/internal/models"

type Repository interface {
	GetCurrencyList() ([]models.Currency, error)
	GetCurrencyRate(string) (*models.Currency, error)
}
