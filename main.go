package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to weather tracker</h1>"))
}

func query(city string) (weatherData, error) {
	fmt.Println("Querying for city:", city)
	apiConfig, err := loadApiCOnfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	response, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=" + apiConfig.OpenWeatherMapApiKey)
	if err != nil {
		return weatherData{}, err
	}
	// fmt.Println("Response:", response)
	defer response.Body.Close()

	var d weatherData
	if err := json.NewDecoder(response.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	fmt.Println(d)
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2] // 3 is substrings
			// fmt.Println("City:", city)
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8000", nil)
}

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`

	Main struct {
		Kelvin float32 `json:"temp"`
	} `json:"main"`
}

func loadApiCOnfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}
