package ddog

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/ddctl/api"
	"net/url"
	"os"
	"path/filepath"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
	"strconv"
	"testing"
)

func TestBuildURL(t *testing.T) {
	type testCase struct {
		Name        string
		InputFile   string
		ExpectedURL string
	}

	cases := []testCase{
		{
			Name:        "basic",
			InputFile:   "basic.yaml",
			ExpectedURL: "https://acme.datadoghq.com/logs?query=RequestLoggingMiddleware%20env%3Aprod%20service%3Afeserver%2A%20%40handler_module%3A%2Abert%2A%20-%40http.method%3AGET%20-%40http.method%3AHEAD%20status%3Aerror%20-%40handler_module%3A%2Alaxmod%2A%20-%40handler%3A%2Alaxmod%2A&agg_m=count&agg_m_source=base&agg_q=status&agg_q_source=base&agg_t=count&clustering_pattern_field_path=message&cols=host%2Cservice&fromUser=true&messageDisplay=inline&refresh_mode=paused&storage=flex_tier&stream_sort=desc&top_n=10&top_o=top&viz=pattern&x_missing=true&from_ts=1736927929003&to_ts=1736949529003&live=false",
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory")
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tFile := filepath.Join(cwd, "test_data", c.InputFile)
			data, err := os.ReadFile(tFile)
			if err != nil {
				t.Fatalf("Failed to read file %v: %v", tFile, err)
			}

			link := &api.DatadogLink{}
			if err := yaml.Unmarshal(data, link); err != nil {
				t.Fatalf("Failed to unmarshal link data: %v", err)
			}

			resultURL, err := BuildURL(link)
			if err != nil {
				t.Fatalf("Error calling LinkToURL: %v", err)
			}

			// We do the comparison in the URL space because it looks
			// like encode parameters alphabetizes them
			uActual, err := url.Parse(resultURL)
			if err != nil {
				t.Fatalf("Failed to parse URL %v: %v", resultURL, err)
			}
			uExpected, err := url.Parse(c.ExpectedURL)
			if err != nil {
				t.Fatalf("Failed to parse URL %v: %v", c.ExpectedURL, err)
			}

			if uActual.Scheme != uExpected.Scheme {
				t.Fatalf("Scheme does not match; got %v; want %v", uActual.Scheme, uExpected.Scheme)
			}

			if uActual.Host != uExpected.Host {
				t.Fatalf("Host does not match; got %v; want %v", uActual.Host, uExpected.Host)
			}

			if uActual.Path != uExpected.Path {
				t.Fatalf("Path does not match; got %v; want %v", uActual.Path, uExpected.Path)
			}

			if d := cmp.Diff(uExpected.Query(), uActual.Query()); d != "" {
				t.Fatalf("URL query does not match; diff\n%v", d)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	type testCase struct {
		Name         string
		Input        string
		ExpectedFile string
	}

	testUrl := "https://acme.datadoghq.com/logs?query=RequestLoggingMiddleware%20env%3Aprod%20service%3Afeserver%2A%20%40handler_module%3A%2Abert%2A%20-%40http.method%3AGET%20-%40http.method%3AHEAD%20status%3Aerror%20-%40handler_module%3A%2Alaxmod%2A%20-%40handler%3A%2Alaxmod%2A&agg_m=count&agg_m_source=base&agg_q=status&agg_q_source=base&agg_t=count&clustering_pattern_field_path=message&cols=host%2Cservice&fromUser=true&messageDisplay=inline&refresh_mode=paused&storage=flex_tier&stream_sort=desc&top_n=10&top_o=top&viz=pattern&x_missing=true&from_ts=1736927929003&to_ts=1736949529003&live=false"
	cases := []testCase{
		{
			Name:         "basic",
			Input:        testUrl,
			ExpectedFile: "basic.yaml",
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory")
	}
	updateEnvValue := os.Getenv("UPDATE_TEST_DATA")
	if updateEnvValue == "" {
		updateEnvValue = "false"
	}
	updateTestData, err := strconv.ParseBool(updateEnvValue)
	if err != nil {
		t.Fatalf("Failed to parse UPDATE_TEST_DATA")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tFile := filepath.Join(cwd, "test_data", c.ExpectedFile)
			link, err := URLToLink(c.Input)
			if err != nil {
				t.Fatalf("Error calling ParseURL: %v", err)
			}

			if updateTestData {
				b, err := yaml.Marshal(link)
				if err != nil {
					t.Fatalf("Failed to marshal link: %v", err)
				}

				if err := os.WriteFile(tFile, b, 0644); err != nil {
					t.Fatalf("Failed to write file %v: %v", tFile, err)
				}
			}

			expectedData, err := os.ReadFile(tFile)
			if err != nil {
				t.Fatalf("Failed to read file %v: %v", tFile, err)
			}
			expected := &api.DatadogLink{}
			if err := yaml.Unmarshal(expectedData, expected); err != nil {
				t.Fatalf("Failed to unmarshal expected data: %v", err)
			}
			if expected.ExtraParams == nil {
				expected.ExtraParams = map[string]string{}
			}

			if d := cmp.Diff(expected, link); d != "" {
				t.Errorf("Link does not match; diff\n%v", d)
			}

		})
	}
}
