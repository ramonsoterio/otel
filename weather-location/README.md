# Instructions
## Running the application
The command below will start the webserver in the port 8080. 
It is necessary to provide a valid Weather API key in the environment variable 
`WEATHER_API_KEY` to run the application. This variable is defined on the Dockerfile.
```
docker-compose up --build
```

## Testing local server
Use the files located in /api directory to run requests against the local webserver and
the deployed webserver.

## Testing webserver
Run the command below to test the deployed webserver:
```curl
curl --location https://weather-location-639493986028.us-central1.run.app/location/88960000
```