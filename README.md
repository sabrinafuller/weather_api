# weather_api

Weather Datastore updates alerts every 10 minutes form national weather service api

To run: 
1. Clone this repo
2. Start redis client `docker run --name redis-container -p 6379:6379 -d redis`
3. Install`go install`
4. Run tests `go test`
4. `go run .`App will start updating every 10 minutes until interrupted 