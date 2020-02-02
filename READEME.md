Run Prometheus
```bash
docker run \
    -p 9090:9090 \
    -v $(pwd)/config/prometheus.yml:/etc/prometheus/prometheus.yml \
    --rm \
    prom/prometheus
```

Run Git Status Exporter
```bash
go build
./git-status-exporter my-repo textfile_collector
```