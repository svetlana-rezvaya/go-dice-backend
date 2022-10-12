# go-dice-backend

[![GoDoc](https://godoc.org/github.com/svetlana-rezvaya/go-dice-backend?status.svg)](https://godoc.org/github.com/svetlana-rezvaya/go-dice-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/svetlana-rezvaya/go-dice-backend)](https://goreportcard.com/report/github.com/svetlana-rezvaya/go-dice-backend)
[![Build Status](https://app.travis-ci.com/svetlana-rezvaya/go-dice-backend.svg?branch=master)](https://app.travis-ci.com/svetlana-rezvaya/go-dice-backend)
[![codecov](https://codecov.io/gh/svetlana-rezvaya/go-dice-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/svetlana-rezvaya/go-dice-backend)

The web service that implements dice rolling.

## Installation

```
$ go install github.com/svetlana-rezvaya/go-dice-backend@latest
```

## Usage

```
$ go-dice-backend
```

Environment variables:

- `PORT` &mdash; service port (default: `8080`).

## Testing

Running the unit tests:

```
$ go test -race -cover ./...
```

Running the load and integration tests:

1.  Run the service:

    ```
    $ docker-compose up -d
    ```

2.  Run [InfluxDB](https://www.influxdata.com/) for storing and [Grafana](https://grafana.com/) for displaying test results (better on the second server):

    ```
    $ docker-compose -f docker-compose-for-load-testing.yml up -d influxdb grafana
    ```

3.  Run [k6](https://k6.io/) on the server from step 2 (replace `localhost:8080` with the address of the service from step 1):

    ```
    $ SERVICE_ADDRESS=localhost:8080 docker-compose -f docker-compose-for-load-testing.yml up k6
    ```

4.  Open http://localhost:3000/ on the server from step 2.

## Docs

[Swagger](https://swagger.io/) specification of the service API: [docs/swagger.yaml](docs/swagger.yaml)

## Output Example

```
2022/10/02 12:23:06 GET /api/v1/dice?throws=51&faces=91 80.332µs
2022/10/02 12:23:06 GET /api/v1/dice?throws=70&faces=15 52.361µs
2022/10/02 12:23:06 GET /api/v1/dice?throws=94&faces=88 94.807µs
2022/10/02 12:23:06 GET /api/v1/dice?throws=89&faces=81 115.173µs
2022/10/02 12:23:07 GET /api/v1/dice?throws=3&faces=13 28.316µs
2022/10/02 12:23:07 GET /api/v1/dice?throws=71&faces=27 40.163µs
2022/10/02 12:23:07 GET /api/v1/dice?throws=50&faces=85 58.52µs
2022/10/02 12:23:07 GET /api/v1/dice?throws=18&faces=49 49.17µs
2022/10/02 12:23:08 GET /api/v1/dice?throws=8&faces=59 29.693µs
2022/10/02 12:23:08 GET /api/v1/dice?throws=56&faces=16 108.716µs
```

## License

The MIT License (MIT)

Copyright &copy; 2022 svetlana-rezvaya
