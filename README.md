# Flights

**Story**: There are over 100,000 flights a day, with millions of people and cargo being transferred around the world. With so many people and different carrier/agency groups, it can be hard to track where a person might be. In order to determine the flight path of a person, we must sort through all of their flight records.

**Goal**: To create a simple microservice API that can help us understand and track how a particular person's flight path may be queried. The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.

## Required JSON structure

```bash
[["SFO", "EWR"]]                                                 => ["SFO", "EWR"]
[["ATL", "EWR"], ["SFO", "ATL"]]                                 => ["SFO", "EWR"]
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]
````

## Specifications

* Your miscroservice must listen on port 8080 and expose the flight path tracker under the /calculate endpoint.
* Define and document the format of the API endpoint in the README.
* Use Golang and/or any tools that you think will help you best accomplish the task at hand.
 
## API Documentation

The API has following endpoints:

* GET `/v1/health`
* POST `/v1/calculate`

### Basic health check

Useful to check microservice status manually or automatically by load balancers.

**Request example:**

```bash
curl -X 'GET' \
'http://localhost:8080/v1/health' \
-H 'accept: application/json'
```

**Response example**

```json
{
  "data": "pong",
  "hostname": "75823749265host",
  "ts": 1684755273
}
```

See more information at: http://localhost:8080/v1/docs/index.html#/healthCheck/healthCheck-get

### Flight tracking

Given a list of passenger flights, it returns the starting point and travel destination point of given passenger.

**Request example:**

```bash
curl -X 'POST' \
    'http://localhost:8080/v1/calculate' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]'
```

**Response example**

```json
[
  "SFO",
  "EWR"
]
```

See more information at: http://localhost:8080/v1/docs/index.html#/FlightsCalculate/flightsCalculate-GET
