package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Value struct {
	Unit  string `xml:"unit,attr" json:"unit"`
	Value string `xml:",chardata" json:"value"`
}

type Timerange struct {
	Type     string  `xml:"type,attr" json:"type"`
	H        string  `xml:"h,attr" json:"h"`
	DateTime string  `xml:"datetime,attr" json:"datetime"`
	Values   []Value `xml:"value" json:"values"`
}

type Parameter struct {
	ID          string      `xml:"id,attr" json:"id"`
	Description string      `xml:"description,attr" json:"description"`
	Type        string      `xml:"type,attr" json:"type"`
	Timeranges  []Timerange `xml:"timerange" json:"timeranges"`
}

type Name struct {
	Language string `xml:"lang,attr" json:"language"`
	Value    string `xml:",chardata" json:"value"`
}

type Area struct {
	Names      []Name      `xml:"name" json:"name"`
	Parameters []Parameter `xml:"parameter" json:"parameter"`
}

type Forecast struct {
	Areas []Area `xml:"area" json:"areas"`
}

type Data struct {
	Forecast Forecast `xml:"forecast" json:"forecast"`
}

type FormattedForecast struct {
	Areas []FormattedArea `json:"areas"`
}

type FormattedArea struct {
	Name       []Name               `json:"name"`
	Parameters []FormattedParameter `json:"parameter"`
}

type FormattedParameter struct {
	ID          string               `json:"id"`
	Description string               `json:"description"`
	Type        string               `json:"type"`
	Timeranges  []FormattedTimerange `json:"timeranges"`
}

type FormattedTimerange struct {
	Type     string           `json:"type"`
	H        string           `json:"h"`
	DateTime time.Time        `json:"datetime"`
	Values   []FormattedValue `json:"values"`
}

type FormattedValue struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

func main() {
	GetForecastBmkg()
}

func GetForecastBmkg() {
	response, err := http.Get("https://data.bmkg.go.id/DataMKG/MEWS/DigitalForecast/DigitalForecast-KalimantanTimur.xml")
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
	var data Data
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the data to the desired format
	formattedForecast := FormattedForecast{
		Areas: make([]FormattedArea, 0),
	}

	for _, area := range data.Forecast.Areas {
		formattedArea := FormattedArea{
			Name:       area.Names,
			Parameters: make([]FormattedParameter, 0),
		}

		for _, param := range area.Parameters {
			formattedParam := FormattedParameter{
				ID:          param.ID,
				Description: param.Description,
				Type:        param.Type,
				Timeranges:  make([]FormattedTimerange, 0),
			}

			for _, tr := range param.Timeranges {
				layout := "200601021504"
				t, err := time.Parse(layout, tr.DateTime)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				formattedValues := make([]FormattedValue, 0)
				for _, value := range tr.Values {
					formattedValue := FormattedValue{
						Unit:  value.Unit,
						Value: value.Value,
					}
					formattedValues = append(formattedValues, formattedValue)
				}

				formattedTR := FormattedTimerange{
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

	// Convert the formatted data to JSON
	jsonData, err := json.MarshalIndent(formattedForecast, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the JSON data to file
	err = ioutil.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Data JSON berhasil disimpan ke file output.json")
}
