package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
param: database struct
param: string array of regions
throws: error
returns string nil on success
Writes to the database the alerts per region
*/
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

// ids, err := redisClient.ZRange(context.Background(), "x:123", 0, -1)
func loginDatabase(addr, password string, db_num int) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db_num,
	})

	return &Database{
		client: rdb,
	}
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

// This app writes to the database every 10 minutes to store alerts
// @ToDo make ttl for entries
// @Todo organize redis better...?
func main() {
	// Login to database
	fmt.Println("Running Weather App----->")
	var update = 10 * 60
	db := loginDatabase("localhost:6379", "", 0)

	// Create current time for key
	var currentTime = time.Now().String()

	// go func to run every 10 minutes
	go func() {
		for {
			fmt.Printf("Running Weather App: Last update %s\n", currentTime)
			currentTime = time.Now().String()
			AddAlert(db, currentTime)
			time.Sleep(time.Duration(update) * time.Second)
		}
	}()
	// Keep the main Goroutine running
	select {}

}
