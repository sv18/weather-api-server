package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const openWeatherAPIKey = "YOUR_API_KEY"

// MainData represents the main temperature data in the weather response.
type MainData struct {
	Temp float64 `json:"temp"`
}

// WeatherData represents the weather description data in the weather response.
type WeatherData struct {
	Description string `json:"description"`
}

// WeatherResponse represents the entire response from the API
type WeatherResponse struct {
	Main    MainData      `json:"main"`
	Weather []WeatherData `json:"weather"`
}

// fetchWeatherData fetches weather data from the OpenWeather API using lat and long coordinates
func fetchWeatherData(lat, long string) (WeatherResponse, error) {
	// construct API URL
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, long, openWeatherAPIKey)

	// perform API request
	resp, err := http.Get(apiURL)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the response
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to decode weather data: %v", err)
	}
	return weatherResp, nil

}

// weatherHandler handles the HTTP request for fetching weather information.
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// parse latitude and longitude from query parameters
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")

	log.Printf("Latitude: %s, Longitude: %s", lat, long)

	// validate latitude and longitude
	if lat == "" || long == "" {
		http.Error(w, "Latitude and longitude are required", http.StatusBadRequest)
		return
	}

	// fetch weather data
	weatherResp, err := fetchWeatherData(lat, long)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// determine weather condition
	weatherCondition := getWeatherCondition(weatherResp)
	log.Printf("Weather Condition: %s\n", weatherCondition)

	// parse temperature from response
	tempCelsius := getTemperatureCelsius(weatherResp)
	log.Printf("Temperature in Celsius: %.2f\n", tempCelsius)

	// convert temperature to Fahrenheit
	tempFahrenheit := celsiusToFahrenheit(tempCelsius)
	log.Printf("Temperature in Fahrenheit: %.2f\n", tempFahrenheit)

	// determine temperature category
	temperature := getTemperatureCategory(tempCelsius)

	// prepare response
	response := fmt.Sprintf("Weather condition: %s\nTemperature: %.2fºC (%.2fºF) (%s)", weatherCondition, tempCelsius, tempFahrenheit, temperature)

	// write response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, response)
}

// getWeatherCondition extracts the weather condition description from the weather response.
func getWeatherCondition(resp WeatherResponse) string {
	if len(resp.Weather) > 0 {
		return resp.Weather[0].Description
	}
	return ""
}

// getTemperatureCelsius extracts the temperature in Celsius from the weather response.
func getTemperatureCelsius(resp WeatherResponse) float64 {
	return resp.Main.Temp - 273.15
}

// getTemperatureCategory determines the temperature category based on Celsius temperature.
func getTemperatureCategory(tempCelsius float64) string {
	if tempCelsius < 5 {
		return "cold"
	} else if tempCelsius > 25 {
		return "hot"
	} else {
		return "moderate"
	}
}

// celsiusToFahrenheit converts temp from Celsius to Fahrenheit.
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

func main() {
	http.HandleFunc("/weather", weatherHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
