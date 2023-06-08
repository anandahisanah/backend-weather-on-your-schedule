package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"see-weather-on-your-schedule/database"
	"see-weather-on-your-schedule/routers"
)

func main() {
	database.StartDB()

	r := routers.StartServer()
	r.Run(":8080")

	// regionSeeder()
}

func regionSeeder() {
	sqlFile, err := ioutil.ReadFile("./database/data/region.sql")
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}

	fmt.Println(string(sqlFile))
}
