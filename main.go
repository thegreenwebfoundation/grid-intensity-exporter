package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	gridintensity "github.com/thegreenwebfoundation/grid-intensity-go"
	"github.com/thegreenwebfoundation/grid-intensity-go/carbonintensity"
	"github.com/thegreenwebfoundation/grid-intensity-go/electricitymap"
)

const (
	carbonIntensityProvider = "carbonintensity.org.uk"
	electricityMapProvider  = "electricitymap.org"

	labelProvider = "provider"
	labelRegion   = "region"
	namespace     = "grid_intensity"

	electricityMapAPITokenEnvVar = "ELECTRICITY_MAP_API_TOKEN"
	providerEnvVar               = "GRID_INTENSITY_PROVIDER"
	regionEnvVar                 = "GRID_INTENSITY_REGION"
)

var (
	actualDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "carbon", "actual"),
		"Actual carbon intensity for the electricity grid in this region.",
		[]string{
			labelProvider,
			labelRegion,
		},
		nil,
	)
)

type Exporter struct {
	apiClient gridintensity.Provider
	provider  string
	region    string
}

func NewExporter(provider, region string) (*Exporter, error) {
	var apiClient gridintensity.Provider
	var err error

	switch provider {
	case carbonIntensityProvider:
		apiClient, err = carbonintensity.New()
		if err != nil {
			return nil, err
		}
	case electricityMapProvider:
		token := os.Getenv(electricityMapAPITokenEnvVar)
		if token == "" {
			return nil, fmt.Errorf("please set the env variable %q", electricityMapAPITokenEnvVar)
		}

		apiClient, err = electricitymap.New(token)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("provider %q not supported", provider)
	}

	exporter := &Exporter{
		apiClient: apiClient,
		provider:  provider,
		region:    region,
	}

	return exporter, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	actualIntensity, err := e.apiClient.GetCarbonIntensity(ctx, e.region)
	if err != nil {
		log.Printf("failed to get carbon intensity %#v", err)
	}

	ch <- prometheus.MustNewConstMetric(
		actualDesc,
		prometheus.GaugeValue,
		actualIntensity,
		e.provider,
		e.region,
	)
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- actualDesc
}

func main() {
	provider := os.Getenv(providerEnvVar)
	region := os.Getenv(regionEnvVar)

	// Default provider to carbonintensity.co.uk since it doesn't need API key.
	if provider == "" {
		provider = carbonIntensityProvider
	}
	// Default region to UK which is the only region supported.
	if provider == carbonIntensityProvider && region == "" {
		region = "UK"
	}

	exporter, err := NewExporter(provider, region)
	if err != nil {
		log.Fatalf("failed to create exporter %#v", err)
	}

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
