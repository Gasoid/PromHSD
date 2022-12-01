# PromHSD


Prometheus http static config discovery service 

PromHSD provides target list for Prometheus through HTTP API. Since version 2.x Prometheus has supported HTTP-based static config.

Official documentation https://prometheus.io/docs/prometheus/latest/http_sd/

## Install
There are various ways to install PromHSD

### Docker
```bash
docker run --name promhsd -d -p 8080:8080 ghcr.io/gasoid/promhsd:latest
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

## API Documentation
Swagger endpoint: /swagger/index.html

Regenerate docs
```
swag init
```
