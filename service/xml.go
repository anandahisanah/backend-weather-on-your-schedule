package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type value struct {
	Unit  string `xml:"unit,attr" json:"unit"`
	Value string `xml:",chardata" json:"value"`
}

type timerange struct {
	Type     string  `xml:"type,attr" json:"type"`
	H        string  `xml:"h,attr" json:"h"`
	DateTime string  `xml:"datetime,attr" json:"datetime"`
	Values   []value `xml:"value" json:"values"`
}

type parameter struct {
	ID          string      `xml:"id,attr" json:"id"`
	Description string      `xml:"description,attr" json:"description"`
	Type        string      `xml:"type,attr" json:"type"`
	Timeranges  []timerange `xml:"timerange" json:"timeranges"`
}

type name struct {
	Language string `xml:"lang,attr" json:"language"`
	Value    string `xml:",chardata" json:"value"`
}

type area struct {
	Names      []name      `xml:"name" json:"name"`
	Parameters []parameter `xml:"parameter" json:"parameter"`
}

type forecast struct {
	Areas []area `xml:"area" json:"areas"`
}

type data struct {
	Forecast forecast `xml:"forecast" json:"forecast"`
}

type formattedForecast struct {
	Areas []formattedArea `json:"areas"`
}

type formattedArea struct {
	Name       []name               `json:"name"`
	Parameters []formattedParameter `json:"parameter"`
}

type formattedParameter struct {
	ID          string               `json:"id"`
	Description string               `json:"description"`
	Type        string               `json:"type"`
	Timeranges  []formattedTimerange `json:"timeranges"`
}

type formattedTimerange struct {
	Type     string           `json:"type"`
	H        string           `json:"h"`
	DateTime time.Time        `json:"datetime"`
	Values   []formattedValue `json:"values"`
}

type formattedValue struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

func GetCity(provinceName string) ([]byte, error) {
	endpoint := fmt.Sprintf("https://data.bmkg.go.id/DataMKG/MEWS/DigitalForecast/DigitalForecast-%s.xml", provinceName)
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	defer response.Body.Close()

	// Read XML response body
	xmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	// Unmarshal XML data into the Data struct
	var data data
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	// Convert the data to the desired format
	formattedForecast := formattedForecast{
		Areas: make([]formattedArea, 0),
	}

	for _, area := range data.Forecast.Areas {
		formattedArea := formattedArea{
			Name: area.Names,
		}
		for _, areaName := range area.Names {
			if areaName.Language == "en_US" {
				formattedArea.Name = []name{areaName}
				break
			}
		}
		formattedForecast.Areas = append(formattedForecast.Areas, formattedArea)
	}

	// Convert the formatted data to JSON
	jsonData, err := json.MarshalIndent(formattedForecast, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return jsonData, nil
}

func CreateJsonForecastBmkg(provinceCode string, endpoint string) {
	response, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Read XML response body
	xmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Unmarshal XML data into the Data struct
	var data data
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the data to the desired format
	formattedForecast := formattedForecast{
		Areas: make([]formattedArea, 0),
	}

	for _, area := range data.Forecast.Areas {
		formattedArea := formattedArea{
			Name:       area.Names,
			Parameters: make([]formattedParameter, 0),
		}

		if len(area.Parameters) > 0 {
			for _, param := range area.Parameters {
				formattedParam := formattedParameter{
					ID:          param.ID,
					Description: param.Description,
					Type:        param.Type,
					Timeranges:  make([]formattedTimerange, 0),
				}

				for _, tr := range param.Timeranges {
					layout := "200601021504"
					t, err := time.Parse(layout, tr.DateTime)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					formattedValues := make([]formattedValue, 0)
					for _, value := range tr.Values {
						formattedValue := formattedValue{
							Unit:  value.Unit,
							Value: value.Value,
						}
						formattedValues = append(formattedValues, formattedValue)
					}

					formattedTR := formattedTimerange{
						Type:     tr.Type,
						H:        tr.H,
						DateTime: t,
						Values:   formattedValues,
					}

					formattedParam.Timeranges = append(formattedParam.Timeranges, formattedTR)
				}

				formattedArea.Parameters = append(formattedArea.Parameters, formattedParam)
			}

			formattedForecast.Areas = append(formattedForecast.Areas, formattedArea)
		}
	}

	// Convert the formatted data to JSON
	jsonData, err := json.MarshalIndent(formattedForecast, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the JSON data to file
	fileName := fmt.Sprintf("%s.json", provinceCode)
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("fileName saved")
}
