![docker](https://github.com/thegreenwebfoundation/grid-intensity-exporter/workflows/docker-integration-test/badge.svg) ![kubernetes](https://github.com/thegreenwebfoundation/grid-intensity-exporter/workflows/kubernetes-integration-test/badge.svg) ![nomad](https://github.com/thegreenwebfoundation/grid-intensity-exporter/workflows/nomad-integration-test/badge.svg)

# grid-intensity-exporter

A grid intensity exporter for use with [prometheus]. Designed to be used for
understanding carbon footprint of compute.

Uses [grid-intensity-go] library to retrieve data from carbon intensity APIs, to cover the most common use case of getting power exclusively from the grid. If you draw power from offgrid sources or batteries please let us know so we can incorproate this in the design too - we're looking for hosts using a set up like this design for it.

## Usage

The exporter currently supports 2 providers. The metrics are exposed on port
8000.

### carbonintensity.org.uk

This is the default provider as it doesn't require an API key. However it only
supports the UK region.

```sh
go run main.go
```

```sh
curl -s http://localhost:8000/metrics | grep grid

# HELP grid_intensity_carbon_actual Actual carbon intensity for the electricity grid in this region.
# TYPE grid_intensity_carbon_actual gauge
grid_intensity_carbon_actual{provider="carbonintensity.org.uk",region="UK"} 293
```

### electricitymap.org

This provider supports more countries but you will need an API key from
[api.electricitymap.org].

```sh
ELECTRICITY_MAP_API_TOKEN=your-api-token
GRID_INTENSITY_PROVIDER=electricitymap.org
GRID_INTENSITY_REGION=IN-KA
go run main.go
```

## Docker

Build the Docker image.

```sh
CGO_ENABLED=0 GOOS=linux go build .
docker build -t thegreenwebfoundation/grid-intensity-exporter:latest .
```

## Kubernetes

Install via the [helm] chart. Needs the Docker image to be available in the
cluster.

```sh
helm install grid-intensity-exporter helm/grid-intensity-exporter
```

## Nomad

Start the Nomad job. Needs the Docker image to be available in the
cluster.

```sh
nomad run ./nomad/grid-intensity-exporter.nomad
```

## Integration Tests

Build and run the docker container.

```sh
CGO_ENABLED=0 GOOS=linux go build .
docker build -t grid-intensity-exporter:integration-test .
docker run -p 8000:8000 grid-intensity-exporter:integration-test
```

Run the integration test.

```sh
go test -v -tags=dockerrequired ./integration/test/docker
```

[api.electricitymap.org]: https://api.electricitymap.org/
[grid-intensity-go]: https://github.com/thegreenwebfoundation/grid-intensity-go
[helm]: https://helm.sh/
[prometheus]: https://prometheus.io/
