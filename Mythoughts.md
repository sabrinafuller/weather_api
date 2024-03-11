1. First set up redis, figured out the go commands to write and read to redis
2. set up  http calls to ensure I'm hitting the endpoint (this was challenging because I didn't realize how go worked at first) getWeatherAlerts
3. Useds the defer keyword, and that means that the `derfer a` is not called until the function surrounding it is done. 
4. to get it working I quered on time stamps, but realisistica

Change: Don't store data based on timestamp?
Change : Want to update on redis x range to store to rank the data
Change: Only update the regions if the alert has changed, requires looking at the most recent call and comparing to current database


For attaching front end: 
Store the updates in redis memory cache
React will read the data in cache and update front end accordingly
React will not writing data only reading data
If we want data from this api to React (ie location weather data) we should use React to talk to web Api for specific location data




https://stackoverflow.com/questions/71982698/testing-with-golang-redis-and-time







