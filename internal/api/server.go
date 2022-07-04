package api

import (
	"currency/internal/store"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  *store.Store
}

func NewServer(store *store.Store) Server {
	s := Server{
		router: gin.Default(),
		store:  store,
	}

	s.configRouter()
	return s
}

func (s Server) Run() error {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	stopChan := make(chan struct{})
	defer close(stopChan)

	go s.periodicUpdateData(ticker, stopChan)
	return s.router.Run(":8080")
}

func (s Server) periodicUpdateData(ticker *time.Ticker, stopChan chan struct{}) {
	data, err := getCurrenciesData()
	if err != nil {
		return 
	}
	s.store.SaveCurrenciesData(data)
	fmt.Println("data saved in db")

	for {
		select {
		case <-ticker.C:
			data, err := getCurrenciesData()
			if err != nil {
				return 
			}

			s.store.UpdateCurrenciesData(data)
			fmt.Println("db updated")

		case <-stopChan:
			return
		}
	}
}

func (s Server) configRouter() {
	s.router.GET("/currencies", s.getCurrencies)
	s.router.GET("/currencies/:code", s.getCurrency)
	s.router.DELETE("/currencies/:code", s.removeCurrency)
}

func (s Server) getCurrencies(c *gin.Context) {
	data, err := s.store.GetCurrencyList()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s Server) getCurrency(c *gin.Context) {
	data, err := s.store.GetCurrency(c.Param("code"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s Server) removeCurrency(c *gin.Context) {
	err := s.store.RemoveCurrency(c.Param("code"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}
