# PromHSD


Prometheus http static config discovery service 

PromHSD provides target list for Prometheus through HTTP API. Since version 2.x Prometheus has supported HTTP-based static config.

Official documentation https://prometheus.io/docs/prometheus/latest/http_sd/

## Install
There are various ways to install PromHSD

### Docker
```bash
docker run --name promhsd -d -p 8080:8080 --env PROMHSD_STORAGE="filedb" --env PROMHSD_FILEDB_ARGS="db.json" --env ghcr.io/gasoid/promhsd:latest
```

### Run from source
```bash
go generate assets.go
go run ./
```

<!-- ### Helm chart
```bash
helm install promhsd 
``` -->

## Prometheus configuration
```yaml
scrape_configs:
  - job_name: httpsd
    http_sd_configs:
      - url: "http://promhsd:8080/prom-target/qwe"

```

### Storages
Now PromHSD supports 2 databases:
- AWS DynamoDB
- file (simple json file)
<!--
- Azure CosmosDB
- Google
-->

## Configuration
| Variable Name  | Default value | Description |
| ------------- | ------------- | ------------- |
| PROMHSD_STORAGE | "" | You should choose storage engine where data will be stored. Possible values: "filedb", "dynamodb"  |
| PROMHSD_FILEDB_ARGS | "" | Filepath, e.g. "temp.json", "/opt/db/file.json". File will be created automatically. |
| PROMHSD_DYNAMODB_ARGS | "" | Table Name, Table will be created automatically. You need to provide usual AWS credentials (env variables, profile and etc) |

## API Documentation
Swagger endpoint: /swagger/index.html

Regenerate docs
```
swag init
```
