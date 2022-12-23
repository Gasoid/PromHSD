# PromHSD

[![codecov](https://codecov.io/gh/Gasoid/PromHSD/branch/main/graph/badge.svg?token=HXLQV248WC)](https://codecov.io/gh/Gasoid/PromHSD)
[![CI](https://github.com/Gasoid/PromHSD/actions/workflows/ci.yml/badge.svg)](https://github.com/Gasoid/PromHSD/actions/workflows/ci.yml)

Prometheus http static config discovery service.

PromHSD provides target list for Prometheus through HTTP API. Since version 2.x Prometheus has supported HTTP-based static config.

Official documentation https://prometheus.io/docs/prometheus/latest/http_sd/

Main purpose of the project is to take advantage of http_sd and to allow devops engineers to use UI instead of static_sd.

## Use cases
In all situations and cases where you need to use either static_sd or file_sd you can use promHSD instead.
- Blackbox exporter targets
- Multiple prometheus instances
- On-premise virtual machines
- Load-balancers (haproxy, nginx)


## Storages
Now PromHSD supports 2 databases:
- MongoDB (so that you can use Atlas, Azure CosmosDB, etc)
- AWS DynamoDB
- file (simple json file)
<!--
- Google
-->

![screen](screen.webp)

## Install
There are various ways to install PromHSD

### Docker
```bash
docker run --name promhsd -d -p 8080:8080 --env PROMHSD_STORAGE="filedb" --env PROMHSD_FILEDB_ARGS="db.json" --env ghcr.io/gasoid/promhsd:latest
```

#### Docker-compose example
```yaml
version: "3.4"
services:
  prometheus:
    depends_on:
      - promhsd
    image: prom/prometheus:v2.40.1
    ports:
      - "9090:9090"
    volumes:
      - ./example/prometheus.yml:/etc/prometheus/prometheus.yml

  promhsd:
    image: ghcr.io/gasoid/promhsd:latest
    environment:
      - PROMHSD_FILEDB_ARGS=/tmp/promhsd.json
      - PROMHSD_STORAGE=filedb
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/"]
      interval: 1m30s
      timeout: 10s
      retries: 3
      start_period: 40s

```


### Run from source
```bash
go generate assets.go
go run ./
```

### Helm chart
```bash
helm install promhsd --set PROMHSD_STORAGE="dynamodb" --set PROMHSD_DYNAMODB_ARGS="tableName" https://github.com/Gasoid/PromHSD/releases/download/v0.0.1/promhsd-0.1.0.tgz
```


## Prometheus configuration

### Simple config to collect metrics from hosts
```yaml
scrape_configs:
  - job_name: httpsd
    http_sd_configs:
      - url: "http://promhsd:8080/prom-target/db1"

```
`/prom-target/%ID%` entrypoint is intended for prometheus, `%ID%` is target id created in promHSD.


### Blackbox config to check host availability
```yaml
scrape_configs:
  - job_name: 'blackbox'
    metrics_path: /probe
    params:
      module: [http_2xx]  # Look for a HTTP 200 response.
    http_sd_configs:
      - url: "http://promhsd:8080/prom-target/websites" # websites is id created in promHSD
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9115  # The blackbox exporter's real hostname:port.
```


## Configuration
| Variable Name  | Default value | Description |
| ------------- | ------------- | ------------- |
| PROMHSD_STORAGE | "" | You should choose storage engine where data will be stored. Possible values: "filedb", "dynamodb", "mongodb"  |
| PROMHSD_FILEDB_ARGS | "" | Filepath, e.g. "temp.json", "/opt/db/file.json". File will be created automatically. |
| PROMHSD_DYNAMODB_ARGS | "" | Table Name, Table will be created automatically. You need to provide usual AWS credentials (env variables, profile and etc) |
| PROMHSD_MONGODB_ARGS | "" | it is MONGODB URI, e.g. mongodb+srv://user:pAssw0rd@cluster0.1tivu8s.mongodb.net/DatabaseName?retryWrites=true&w=majority |

## API Documentation
Swagger endpoint: /swagger/index.html

Regenerate docs
```
swag init
```
