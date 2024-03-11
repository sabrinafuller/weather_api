package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

/*
param: database struct
param: queryString
throws: error
returns response weather alerts
*/
func getWeatherAlerts(query string) (*http.Response, error) {
	response, err := http.Get(query)
	if err != nil {
		return response, fmt.Errorf("error with api call %s", err)
	}
	defer response.Body.Close()
	return response, err
}

/*
param: database struct
param: string array of regions
throws: error
returns string nil on success
Writes to the database the alerts per region
*/
func getAlertPerRegion(db *RedisDatabase, regions []string) (string, error) {
	var query = "https://api.weather.gov/alerts/active/region/"

	for _, value := range regions {
		response, err := getWeatherAlerts(query + value)
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
func AddAlert(db *RedisDatabase, currentTime string) (string, error) {
	response, alert_err := getWeatherAlerts("https://api.weather.gov/alerts/active/count")
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
func CheckAlerts(db *RedisDatabase, currentTime string) (string, error) {
	resp, err := db.client.Get(context.Background(), currentTime).Result()
	if err != nil {
		return "", err

	}
	return resp, nil

}

/*
Returns the reference to the Redis database
*/
func loginDatabase(addr, password string, db_num int) *RedisDatabase {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db_num,
	})

	return &RedisDatabase{
		client: rdb,
	}
}

// This app writes to the database every 10 minutes to store alerts
// @ToDo make ttl for entries
// @Todo organize redis better...?
func main() {
	// Login to database
	err := godotenv.Load()

	if err != nil {
		fmt.Errorf("Error loading .env file: %v", err)
	}

	// Access the loaded environment variables
	update, _ := strconv.Atoi(os.Getenv("UPDATE"))
	database := os.Getenv("DATABASE")
	password := os.Getenv("PASSWORD")

	fmt.Println("Running Weather App----->")
	db := loginDatabase(database, password, 0)

	// Create current time for key
	var currentTime = time.Now().String()

	// for loop gets the updates, writes to db and then sleeps
	go func() {
		for {
			fmt.Printf("Running Weather App: Last update %s\n", currentTime)
			currentTime = time.Now().String()
			_, err = AddAlert(db, currentTime)
			if err != nil {
				fmt.Errorf("Error adding alert to db %w", err)
			}
			time.Sleep(time.Duration(update) * time.Second)
		}
	}()
	// Keep the main Goroutine running
	select {}

}
