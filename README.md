# weather_api

Weather Datastore updates alerts every 10 minutes form national weather service api

To run: 
1. Clone this repo
2. Start redis client `docker run --name redis-container -p 6379:6379 -d redis`
3. Install`go install`
4. Run tests `go test`
4. `go run .`App will write to the database every 10 minutes until interrupted

## Additional Clean up
1. Create varfile for variables in weatherApp (for deploying)
2. Organize data from api and better format for redis use
4. Create ttl for data in redis so that expired alerts are deleted accordingly
5. `getAlertPerRegion` is not used in main though tests were written
6. Add additional data for tests to check
7. Not sure if this project captured the "spirit" of the Go will need to learn more about this language 
   
## Assumptions

1. This is a barebones application where it just needs to pull data from weather api into redis database
2. Assumed that the data is formatted correctly and there is no additional metadata we want to add
3. Assumed that any front end we add will use the data as is / have more flexability with parsing jsons, although logstically we probably want to localize ie better in redis
4. Added test to ensure api was working and go was parsing things as I went

## Things I learned about go
First time coding in Go, so had to learn things as I went. Still trying to wrap my head around not using objects in the way I was used to, so the design of this project isn't probably the best laid out. I attempted to use an interface was ran into issues with debugging. 
Got myself a bit confused by marshalling/unmarshalling json, and the type checking/ conversion issue, and finally developed a love hate relationship with the compiler (python has made me lazy). Overall I the more barebones feel of golang. 
