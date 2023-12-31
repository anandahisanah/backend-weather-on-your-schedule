package main

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"backend-weather-on-your-schedule/routers"
	"backend-weather-on-your-schedule/service"
	"backend-weather-on-your-schedule/service/create"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type FormattedForecast struct {
	Areas []FormattedArea `json:"areas"`
}

type FormattedArea struct {
	Name []Name `json:"name"`
}

type Name struct {
	Language string `xml:"lang,attr" json:"language"`
	Value    string `xml:",chardata" json:"value"`
}

var envSeed = os.Getenv("SEED")

func main() {
	database.StartDB()

	if envSeed == "true" {
		seederProvince()
		seederCity()
		createForecastFromJson()
	}

	go runForecastJob()

	router := routers.StartServer()
	router.Run()
}

func runForecastJob() {
    duration := 400 * time.Minute

    for {
        createForecastFromJson()

        // waited for 10 minutes before running again
        time.Sleep(duration)
    }
}

func createForecastFromJson() {
	fmt.Println("Executing Goroutine")

	db := database.GetDB()

	// get province
	provinces := []models.Province{}
	err := db.Find(&provinces).Preload("Cities").Error
	if err != nil {
		log.Fatalln(err)
	}

	for _, province := range provinces {
		// create json by province
		service.CreateJsonForecastBmkg(province.Code, province.Endpoint)
		create.Forecast(province.Code)
	}
	fmt.Println("All forecast successfully saved to Database")
}

func seederProvince() {
	// find json file
	filePath := "./database/data/provinces.json"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Error opening JSON file:", err)
	}

	// parse json to slice of struct
	var provinces []models.Province
	err = json.Unmarshal(file, &provinces)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	for _, item := range provinces {
		code := strings.Replace(item.Name, " ", "", -1)
		endpoint := fmt.Sprintf("https://data.bmkg.go.id/DataMKG/MEWS/DigitalForecast/DigitalForecast-%s.xml", code)

		province := models.Province{
			Code:     code,
			Name:     item.Name,
			Endpoint: endpoint,
		}

		// create
		err := database.GetDB().Create(&province).Error
		if err != nil {
			log.Fatal("Error saving to database:", err)
		}
	}

	fmt.Println("Seeding Province complete")
}

func seederCity() {
	provinces := []models.Province{}
	err := database.GetDB().Find(&provinces).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, province := range provinces {
		jsonData, err := service.GetCity(province.Code)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Decode the JSON data
		var formattedForecast FormattedForecast
		err = json.Unmarshal(jsonData, &formattedForecast)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Extract city names from the formatted forecast
		cityNames := make([]string, 0)
		for _, formattedArea := range formattedForecast.Areas {
			for _, name := range formattedArea.Name {
				cityNames = append(cityNames, name.Value)
			}
		}

		// Save city names to the database
		for _, cityName := range cityNames {
			city := models.City{
				ProvinceID: int(province.ID),
				Name:       cityName,
			}
			err := database.GetDB().Create(&city).Error
			if err != nil {
				fmt.Println("Error saving to database:", err)
				return
			}
		}

		fmt.Printf("City names saved to the database for province: %d. %s\n", index+1, province.Name)
	}
}
