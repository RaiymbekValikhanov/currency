package api

import (
	"currency/internal/models"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	API_URL = "https://nationalbank.kz/rss/get_rates.cfm?fdate="
)

func getCurrenciesData() (data models.NBKData, err error) {
	curTime := time.Now()
	curDate := fmt.Sprintf("%02d.%02d.%04d", curTime.Day(), curTime.Month(), curTime.Year())

	FULL_URL := API_URL + curDate

	res, err := http.Get(FULL_URL)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	
	if err = xml.Unmarshal(body, &data); err != nil {
		return 
	}

	return
}