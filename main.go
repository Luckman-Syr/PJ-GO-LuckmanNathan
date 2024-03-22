package main

import (
	"project-akhir/database"
	"project-akhir/routers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load env file")
	}
	database.StartDB()
	r := routers.StartRouter()
	r.Run()
}
