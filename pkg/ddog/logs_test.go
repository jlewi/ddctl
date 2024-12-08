package ddog

import (
	"testing"
	"time"
)

func Test_GetLogsLink(t *testing.T) {
	type testCase struct {
		Name     string
		LogQuery map[string]string
		BaseURL  string
		Expected string
	}

	// Define the time layout that matches the input format
	const layout = "2006-01-02 15:04 MST"

	// The time string you want to parse
	startString := "2024-12-06 15:20 PST"
	endString := "2024-12-06 15:40 PST"

	// Load the location for PST
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("Error loading location: %v\n", err)
		return
	}

	// Parse the time string using ParseInLocation to include timezone information
	// N.B. When we used time.Parse we weren't getting the same values when running in GHA as when running locally.
	// So I don't know that time.Parse properly deals with timezones that aren't specified as offsets.
	startTime, err := time.ParseInLocation(layout, startString, location)
	if err != nil {
		t.Fatalf("Error parsing start time: %v", err)
		return
	}

	endTime, err := time.ParseInLocation(layout, endString, location)
	if err != nil {
		t.Fatalf("Error parsing end time: %v", err)
		return
	}

	timeArgs := BuildTimeRange(startTime, endTime)

	cases := []testCase{
		{
			Name: "basic",
			LogQuery: map[string]string{
				"query":          "service:foyle @contextId:01JEF30X8B9A8K5M7XGQMAPQ2Y",
				"from_ts":        timeArgs["from_ts"],
				"stream_sort":    "desc",
				"viz":            "stream",
				"to_ts":          timeArgs["to_ts"],
				"agg_m":          "count",
				"agg_m_source":   "base",
				"cols":           "host,service",
				"fromUser":       "true",
				"live":           "false",
				"agg_t":          "count",
				"messageDisplay": "inline",
				"refresh_mode":   "sliding",
				"storage":        "flex_tier",
			},
			BaseURL:  "https://datadoghq.com",
			Expected: "https://datadoghq.com/logs?agg_m=count&agg_m_source=base&agg_t=count&cols=host%2Cservice&fromUser=true&from_ts=1733527200000&live=false&messageDisplay=inline&query=service%3Afoyle+%40contextId%3A01JEF30X8B9A8K5M7XGQMAPQ2Y&refresh_mode=sliding&storage=flex_tier&stream_sort=desc&to_ts=1733528400000&viz=stream",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actual := GetLogsLink(c.BaseURL, c.LogQuery)

			if actual != c.Expected {
				t.Errorf("Got %v;\n Want %v", actual, c.Expected)
			}
		})
	}
}
