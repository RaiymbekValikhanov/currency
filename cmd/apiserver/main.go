package main

import (
	"currency/internal/api"
)

func main() {
	if err := api.Start(); err != nil {
		panic(err)
	}
}