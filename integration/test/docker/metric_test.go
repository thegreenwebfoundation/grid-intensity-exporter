// +build dockerrequired

package docker

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_GridIntensityMetric(t *testing.T) {
	metricsURL := "http://localhost:8000/metrics"

	resp, err := http.Get(metricsURL)
	if err != nil {
		t.Fatalf("could not retrieve %s: %v", metricsURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d: got %d", http.StatusOK, resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	metrics := string(respBytes)
	expectedMetricText := "grid_intensity_carbon_actual{provider=\"carbonintensity.org.uk\",region=\"UK\"}"

	if !strings.Contains(metrics, expectedMetricText) {
		t.Fatalf("expected metric text %q not found", expectedMetricText)
	}
}
