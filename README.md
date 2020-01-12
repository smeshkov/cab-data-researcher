# cab-data-researcher
Welcome to "Cab Data Researcher".

## Requirements

* [Go 1.9+](https://golang.org/dl/);
* [Docker](https://www.docker.com/get-started).

## Startup

To start an application with the DB inside Docker use `make up`.

## Usage

### CLI

### cURL

* Healthcheck `curl http://localhost:8080/health`;
* Get count of trips for medallion(s) on certain date:
    * e.g. cached result: `curl -X POST --data '{ "medallions": ["D7D598CD99978BD012A87A76A7C891B7"], "pickupDate": "2013-12-01" }' http://localhost:8080/api/v1/trip/count`;
    * e.g. not cached result: `curl -X POST --data '{ "medallions": ["D7D598CD99978BD012A87A76A7C891B7"], "pickupDate": "2013-12-01", "noCache": true }' http://localhost:8080/api/v1/trip/count`;
* To clear cache `curl -X POST http://localhost:8080/api/v1/cache/clear`

