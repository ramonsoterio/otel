# Instructions
## Running the application
The command below will start all the services.
It is necessary to provide a valid Weather API key in the environment variable
`WEATHER_API_KEY` to run the application. This variable is defined on the Dockerfile of the 
weather-location project.
```
docker-compose up --build
```

## Testing
Use the file located in cep-receiver/api directory to run requests against the main service. The traces
can be seen on the zipkin service, which can be viewed on http://localhost:9411/zipkin/ once the containers are started.