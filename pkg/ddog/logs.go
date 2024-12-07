package ddog

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"
)

// TimeAndDurationToRange returns the query arguments for the given time range
// based on the end time and duration. StartTime is calculated by subtracting duration from endTime.
func TimeAndDurationToRange(endTimeString string, layout string, length time.Duration) (map[string]string, error) {
	endTime := time.Now()
	if endTimeString != "" {
		var err error
		endTime, err = time.Parse(layout, endTimeString)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing time %v; with layout %v", endTimeString, layout)
		}
	}

	startTime := endTime.Add(-length)

	return BuildTimeRange(startTime, endTime), nil
}

// BuildTimeRange returns the query arguments for the given time range.
func BuildTimeRange(start time.Time, end time.Time) map[string]string {
	return map[string]string{
		// Multiple by 1000 because we want it in milliseconds
		"from_ts": fmt.Sprintf("%d", start.Unix()*1000),
		"to_ts":   fmt.Sprintf("%d", end.Unix()*1000),
	}
}

// GetLogsLink returns a link to the Datadog logs matching the given query.
func GetLogsLink(baseUrl string, query map[string]string) string {
	// Create a new url.Values object
	queryParams := url.Values{}

	// Add map values to the url.Values object
	for key, value := range query {
		queryParams.Add(key, value)
	}

	// Encode the values into a query string
	encodedQuery := queryParams.Encode()
	u := fmt.Sprintf("%s/logs?%s", baseUrl, encodedQuery)
	return u
}

func buildLogsQuery(labels map[string]string) string {
	// We want the names to appear in sorted order in the link so the link is deterministic.
	names := make([]string, 0, len(labels))
	for n := range labels {
		names = append(names, n)
	}
	sort.Strings(names)

	labelsQuery := []string{}
	for _, n := range names {
		labelsQuery = append(labelsQuery, fmt.Sprintf(`%s:"%s"`, n, labels[n]))
	}
	return url.QueryEscape(strings.Join(labelsQuery, " "))
}

func parseUrl(rawURL string) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Extract query parameters
	queryParams := parsedURL.Query()

	// Iterate over all query parameters
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}
