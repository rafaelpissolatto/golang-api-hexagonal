# golang-api-hexagonal
API Template with Golang, Hexagonal Architecture, OpenApi, Prometheus, Postgres, Redis, Kafka and OpenPolice.

# To simply run the application
- Define `resources/.env` parameters (please base the file on `resources/.env.tpl`)
- Before create and load the required environment parameters: `export $(grep -v '^#' .env | xargs -d '\n')`
- docker compose up
- To run the application locally: `go run cmd/main.go`

### Install Docker and docker-compose
- To install Docker follow the official documentation: https://docs.docker.com/engine/install/ubuntu/
- To install docker compose: `sudo apt install docker-compose`
- Give permission to execute: `sudo chmod 666 /var/run/docker.sock`

### How to install an updated dependency
- This example is a dependency to read yaml files: 
`go get gopkg.in/yaml.v3`

### Generate code based on the openapi documentation:
- Install the latest version of "openapi-generator-cli". On linux you can use the npm to install:
- `npm install @openapitools/openapi-generator-cli -g`
- `openapi-generator-cli generate \
  -g go-server -i openapi.yaml \
  -o server/ \
  --additional-properties=outputAsLibrary=true,sourceFolder=openapi`

### Kafka interface
After run the project, you can access the Kafdrop on:
- `localhost:9000` 

On this Kafka interface you can see that the kafka topic was created.

### Run Flyway Database Migration
- Do the database migration
```
  docker run --rm --network="host" -v $(pwd)/resources/database_migrations:/flyway/sql flyway/flyway -url="jdbc:postgresql://127.0.0.1:5432/db?user=root&password=root" -baselineOnMigrate="false" migrate
```
- Check the database status
```
  docker run --rm --network="host" -v $(pwd)/resources/database_migrations:/flyway/sql flyway/flyway -url="jdbc:postgresql://127.0.0.1:5432/db?user=root&password=root" -baselineOnMigrate="false" info
```

#### Endpoints for Health Check:
- GET `http://localhost:8080/health/live`
- GET `http://localhost:8080/health/ready`

#### Prometheus endpoint with Go and Http metrics with custom service_name label:
- GET `http://localhost:8080/metrics`

#### Endpoint to simulate user authentication returning jwt. It is not a full implementation with user database:
- POST `http://localhost:8080/v1/sts/token`

#### Authenticated endpoint to create product:
- POST `http://localhost:8080/v1/product`

#### Authenticated endpoint to get product by id:
- GET `http://localhost:8080/v1/product/{id}`
