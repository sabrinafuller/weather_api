package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestAddAlert(t *testing.T) {
	var currentTime = time.Now().String()
	// Build test database
	db := loginDatabase("localhost:6379", "", 0)
	AddAlert(db, currentTime)
	result, err := CheckAlerts(db, currentTime)
	want := ""

	if result != want && err != nil {
		t.Fail()
	}
}
func TestRead_db_withError(t *testing.T) {
	db := loginDatabase("localhost:6379", "", 0)
	// Would mock redis db, but not able to get it working
	var currentTime = time.Now().String()
	db.write_db(currentTime, "test")
	time.Sleep(time.Second)
	currentTime = time.Now().String()
	result, err := db.read_db(currentTime)
	if result != "" && err != nil {
		t.Fail()
	}
}
func TestRead_db(t *testing.T) {
	var currentTime = time.Now().String()
	// Would mock redis db, but not able to get it working
	db := loginDatabase("localhost:6379", "", 0)
	db.write_db(currentTime, "test")
	result, err := db.read_db(currentTime)
	if err == nil || result == currentTime {
		t.Fail()
	}
}
func TestWriteDB(t *testing.T) {
	var currentTime = time.Now().String()
	// Would mock redis db, but not able to get it working
	db := loginDatabase("localhost:6379", "", 0)
	result := db.write_db(currentTime, "test_data")

	if result != nil {
		t.Fail()
	}
}

func TestWriteAlert(t *testing.T) {
	var currentTime = time.Now().String()
	db := loginDatabase("localhost:6379", "", 0)
	// wri
	result, err := AddAlert(db, currentTime)

	if result != "" && err != nil {
		t.Fail()
	}
}
func TestGetWeatherAlert(t *testing.T) {

	db := loginDatabase("localhost:6379", "", 0)
	query := "https://api.weather.gov/alerts/active/count"
	result, err := getWeatherAlerts(db, query)
	fmt.Print(result.Body)
	if result.Status != string(rune(http.StatusOK)) && err != nil {
		t.Fail()
	}
}

func TestGetWeatherAlertQuery(t *testing.T) {

	db := loginDatabase("localhost:6379", "", 0)
	result, alert_err := getWeatherAlerts(db, "https://api.weather.gov/alerts/active/count")
	bodyBytes, io_err := io.ReadAll(result.Body)
	fmt.Print(string(bodyBytes))
	if bodyBytes != nil && alert_err != nil && io_err != nil {
		t.Fail()
	}

}

func TestGetAlertPerRegion(t *testing.T) {
	test_list := []string{"AL", "AT"}
	db := loginDatabase("localhost:6379", "", 0)
	result, err := getAlertPerRegion(db, test_list)
	resultdata, err_data := db.read_db("AL")
	if result != resultdata && err == err_data {
		t.Fail()
	}

}
