package main

import (
	"github.com/FoolVPN-ID/megalodon-api/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	api.StartApi()
}
