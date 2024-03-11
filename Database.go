package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Database interface {
	write_db(key string, data string) error
	read_db(key string) error 
}
/* Returns reference to redisClient
*/ 
type RedisDatabase struct {
	client *redis.Client
}

/*
param: key string to insert data in db
param: data string of data to insert in db
throw: error if issue marshalling json or setting data in redis
returns: nil
*/
func (db *RedisDatabase) write_db(key string, data string) error {
	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error converting struct to JSON: %v", err)
	}

	// Set the JSON string in Redis
	err = db.client.Set(context.Background(), key, jsonData, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting data in Redis: %v", err)
	}

	return nil
}

/*
param: key string to search in db
throws: error
returns: string
*/
func (db *RedisDatabase) read_db(key string) (string, error) {
	// Get the stored JSON data from Redis
	jsonData, err := db.client.Get(context.Background(), key).Result()
	// db.client.
	if err != nil {
		return jsonData, fmt.Errorf("error getting data from redis")

	}
	return jsonData, err

}
