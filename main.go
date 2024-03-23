package main

import (
	"project-akhir/database"
	"project-akhir/routers"
)

func main() {
	database.StartDB()
	r := routers.StartRouter()
	r.Run()
}
