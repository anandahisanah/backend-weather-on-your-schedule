package create

import (
	"backend-weather-on-your-schedule/database"
	"backend-weather-on-your-schedule/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type data struct {
	Areas []area `json:"areas"`
}

type area struct {
	Names      []name      `json:"name"`
	Parameters []parameter `json:"parameter"`
}

type name struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type parameter struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Timeranges  []timerange `json:"timeranges"`
}

type timerange struct {
	Type     string    `json:"type"`
	Hour     string    `json:"h"`
	Datetime time.Time `json:"datetime"`
	Values   []value   `json:"values"`
}

type value struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

func Forecast(provinceCode string) {
	// Open JSON file
	fileName := fmt.Sprintf("%s.json", provinceCode)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Open JSON file:", err)
	}
	defer file.Close()

	// Decode JSON data
	decoder := json.NewDecoder(file)
	var data data
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to the database
	db := database.GetDB()

	// Disable log output
	log.SetOutput(ioutil.Discard)

	for _, area := range data.Areas {
		// Search city
		for _, name := range area.Names {
			if name.Language == "en_US" {
				city := models.City{}
				resultCity := db.Where("name = ?", name.Value).Preload("Province").First(&city)
				if resultCity.Error != nil {
					log.Println("Search city:", resultCity.Error)
				}

				// Iterate through parameters
				for _, parameter := range area.Parameters {
					if parameter.ID == "t" || parameter.ID == "ws" || parameter.ID == "hu" || parameter.ID == "weather" {
						// Iterate through timeranges
						for _, timerange := range parameter.Timeranges {
							datetime := timerange.Datetime

							// Search existing forecast by datetime and city
							existingForecast := models.Forecast{}
							result := db.Where("datetime = ? AND city_id = ?", datetime, city.ID).First(&existingForecast)
							if result.Error != nil {
								if result.Error != gorm.ErrRecordNotFound {
									log.Println("Search existing forecast:", result.Error)
								}

								// Create forecast if not found
								forecast := models.Forecast{
									ProvinceID: city.ProvinceID,
									CityID:     int(city.ID),
									Datetime:   &datetime,
								}

								// Save forecast to the database
								forecastDB := db.Create(&forecast)
								if forecastDB.Error != nil {
									log.Println("Error saving forecast to database:", forecastDB.Error)
								}

								existingForecast = forecast
							}

							// Update values based on parameter ID
							for _, value := range timerange.Values {
								switch parameter.ID {
								case "t":
									if value.Unit == "C" {
										existingForecast.Temperature = value.Value
									}
								case "ws":
									if value.Unit == "MPH" {
										existingForecast.WindSpeed = value.Value
									}
								case "hu":
									existingForecast.Humidity = value.Value
								case "weather":
									if value.Value == "0" {
										existingForecast.Weather = "Cerah"
									} else if value.Value == "1" {
										existingForecast.Weather = "Cerah Berawan"
									} else if value.Value == "3" {
										existingForecast.Weather = "Berawan"
									} else if value.Value == "4" {
										existingForecast.Weather = "Berawan Tebal"
									} else if value.Value == "45" {
										existingForecast.Weather = "Kabut"
									} else if value.Value == "60" {
										existingForecast.Weather = "Hujan Ringan"
									} else if value.Value == "61" {
										existingForecast.Weather = "Hujan Sedang"
									} else if value.Value == "95" {
										existingForecast.Weather = "Hujan Petir"
									} else {
										existingForecast.Weather = "Tidak diketahui"

									}
								}
							}

							// Save forecast to the database
							forecastDB := db.Save(&existingForecast)
							if forecastDB.Error != nil {
								log.Println("Error updating forecast in database:", forecastDB.Error)
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("%s successfully imported\n\n", fileName)
}
