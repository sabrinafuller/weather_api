package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
param: database struct
param: string array of regions
throws: error
returns string nil on success
Writes to the database the alerts per region
*/
/// current time should be Time Object
type App struct {
	weatherApp *WeatherApp
}
type WeatherApp interface {
	getRegions(db *Database)
	getAlertPerRegion(db *Database, regions []string)
	AddAlert(db *Database, currentTime string)
	CheckAlerts(db *Database, currentTime string)
	getWeatherAlerts(db *Database, query string)
}

func getAreaAlert(db *Database, area_id string) (string, error) {
	query := "https://api.weather.gov/alerts/active/area/" + area_id
	response, err := getWeatherAlerts(db, query)
	if err != nil {
		return response.Status, fmt.Errorf("error getting alert with query : %v", query)
	}
	bodyBytes, io_err := io.ReadAll(response.Body)
	if io_err != nil {
		return response.Status, fmt.Errorf("error reading data: %v", bodyBytes)
	}
	db.write_db(area_id, string(bodyBytes))

	return "", nil
}
func getRegions(db *Database) ([]string, error) {
	query := "https://api.weather.gov/alerts/active/count"
	response, err := getWeatherAlerts(db, query)
	var alert_count alert_count
	if err != nil {
		return nil, fmt.Errorf("error with api call %v", err)
	}
	bodyBytes, io_err := io.ReadAll(response.Body)
	// jsonData, json_err := json.Marshal(bodyBytes)
	json_err := json.Unmarshal([]byte(bodyBytes), &alert_count)
	if json_err != nil && io_err != nil {
		return nil, fmt.Errorf("error unmarshalling json %v", json_err)
	}
	// for regionCode, count := range alert_count.Regions {
	// 	fmt.printf("  %s: %d\n", regionCode, count)

	// }
	return alert_count.Regions, nil

}

func getAlertPerRegion(db *Database, regions []string) (string, error) {
	var query = "https://api.weather.gov/alerts/active/region/"

	for _, value := range regions {
		response, err := getWeatherAlerts(db, query+value)
		if err != nil {
			return response.Status, fmt.Errorf("error with api call %v", err)
		}
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return response.Status, fmt.Errorf("error with api call %v", err)
		}
		db.write_db(value, string(bodyBytes))
	}
	return "", nil
}

/*
param: database struct
param: currentTime string
throws: error
returns string nil on success
*/
func AddAlert(db *Database, currentTime string) (string, error) {
	response, alert_err := getWeatherAlerts(db, "https://api.weather.gov/alerts/active/count")
	bodyBytes, io_err := io.ReadAll(response.Body)
	if alert_err != nil && io_err != nil {
		return "", fmt.Errorf("error with api call %s", alert_err)
	}

	db.write_db(currentTime, string(bodyBytes))
	return string(bodyBytes), nil

}

/*
CheckAlerts retrieves current time entry from
param: database struct
param: currentTime string
throws: error
returns string nil on success
*/
func CheckAlerts(db *Database, currentTime string) (string, error) {
	resp, err := db.client.Get(context.Background(), currentTime).Result()
	if err != nil {
		return "", err

	}
	return resp, nil

}

/*
param: database struct
param: queryString
throws: error
returns response weather alerts
*/
func getWeatherAlerts(db *Database, query string) (*http.Response, error) {
	response, err := http.Get(query)
	if err != nil {
		return response, fmt.Errorf("error with api call %s", err)
	}
	defer response.Body.Close()
	return response, err
}

func runApp(db *Database, update int) *WeatherApp { // go func to run every 10 minutes
	var currentTime = time.Now().String()
	go func() {
		for {
			fmt.Printf("Running Weather App: Last update %s\n", currentTime)
			currentTime = time.Now().String()

			activeRegions, err := getRegions(db)
			if err != nil {
				fmt.Errorf("error finding active regions %v", err)
			}
			getAlertPerRegion(db, activeRegions)
			time.Sleep(time.Duration(update) * time.Second)
		}
	}()
	// Keep the main Goroutine running
	select {}

}

// This app writes to the database every 10 minutes to store alerts
// @ToDo make ttl for entries
// @Todo organize redis better...?
func main() {
	// Login to database
	fmt.Println("Running Weather App----->")
	db := loginDatabase("localhost:6379", "", 0)
	app := runApp(db, 10)
	fmt.Printf("app: %v\n", app)

	if app != nil {
		defer fmt.Printf("App error: %v\n", app)
	}

	// Create current time for key

}
