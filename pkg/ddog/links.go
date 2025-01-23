package ddog

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/ddctl/api"
	"github.com/jlewi/grafctl/pkg/grafana"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var (
	timeParser *grafana.RelativeTimeParser
)

func init() {
	timeParser = grafana.NewRelativeTimeParser()
}

func addString(values url.Values, name string, value string) {
	if value == "" {
		return
	}

	values.Add(name, value)
}

// relativeToAbsoluteTime handles time expressions with "now" in them.
func relativeToAbsoluteTime(timeVal string) (string, error) {
	if !strings.Contains(timeVal, "now") {
		return timeVal, nil
	}

	newTime, err := timeParser.ParseGrafanaRelativeTime(timeVal)
	if err != nil {
		return "", errors.Wrapf(err, "Error parsing relative time %v", timeVal)
	}

	// Datadog is unix epoch in milliseconds
	newTimestr := fmt.Sprintf("%d", newTime.Unix()*1000)
	return newTimestr, nil
}

func BuildURL(link *api.DatadogLink) (string, error) {
	// Create a new url.Values object
	queryParams := url.Values{}

	from_ts, err := relativeToAbsoluteTime(link.FromTS)
	if err != nil {
		return "", errors.Wrapf(err, "Error converting from_ts relative to absolute time for %v", link.FromTS)
	}
	to_ts, err := relativeToAbsoluteTime(link.ToTS)
	if err != nil {
		return "", errors.Wrapf(err, "Error converting to_ts relative to absolute time for %v", link.ToTS)
	}
	addString(queryParams, "query", link.Query)
	addString(queryParams, "viz", link.VisualizeAs)
	addString(queryParams, "agg_m", link.GroupInto)
	addString(queryParams, "storage", link.Storage)
	addString(queryParams, "x_missing", link.Missing)
	addString(queryParams, "agg_m_source", link.Source)
	addString(queryParams, "agg_q", link.GroupBy)
	addString(queryParams, "clustering_pattern_field_path", link.ClusteringPatternFieldPath)
	addString(queryParams, "stream_sort", link.StreamSort)
	addString(queryParams, "agg_q_source", link.GroupBySource)
	addString(queryParams, "agg_t", link.AggType)
	addString(queryParams, "refresh_mode", link.RefreshMode)
	addString(queryParams, "from_ts", from_ts)
	addString(queryParams, "to_ts", to_ts)
	addString(queryParams, "fromUser", link.FromUser)
	addString(queryParams, "top_n", strconv.Itoa(link.TopN))
	addString(queryParams, "top_o", link.TopO)
	addString(queryParams, "live", strconv.FormatBool(link.Live))
	addString(queryParams, "cols", strings.Join(link.Columns, ","))
	addString(queryParams, "messageDisplay", link.MessageDisplay)
	// Encode the values into a query string
	encodedQuery := queryParams.Encode()
	u := fmt.Sprintf("%s/logs?%s", link.BaseURL, encodedQuery)
	return u, nil
}

// ParseURL parses the input URL and returns
// baseUrl - The base URL
// a map of query parameters
// The panes object.
func ParseURL(inputURL string) (string, map[string][]string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to parse URL: %v", inputURL)
	}

	values := parsedURL.Query()

	// Get only the scheme and host
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	// Clean the path to keep only the base path without the last segment
	cleanedPath := path.Clean("/")
	baseURL = baseURL + cleanedPath

	baseURL = strings.TrimSuffix(baseURL, "/")
	return baseURL, values, nil
}

// queryValueHandler is a function that takes a url.Values and returns a function
// that will parse the value from the url.Values
type queryValHandler func(values []string)

func bindToString(field *string) queryValHandler {
	return func(values []string) {
		if len(values) > 0 {
			*field = values[0]
		}
	}
}

func bindToInt(field *int) queryValHandler {
	log := zapr.NewLogger(zap.L())
	return func(values []string) {
		if len(values) > 0 {
			number, err := strconv.Atoi(values[0])
			if err != nil {
				log.Error(err, "Failed to parse an integer", "value", values[0])
			}
			*field = number
		}
	}
}

func bindToBool(field *bool) queryValHandler {
	log := zapr.NewLogger(zap.L())
	return func(values []string) {
		if len(values) > 0 {
			number, err := strconv.ParseBool(values[0])
			if err != nil {
				log.Error(err, "Failed to parse a boolean", "value", values[0])
			}
			*field = number
		}
	}
}

func bindToStringSlice(field *[]string) queryValHandler {
	return func(values []string) {
		if len(values) > 0 {
			items := strings.Split(values[0], ",")
			*field = append(*field, items...)
		}
	}
}

// URLToLink converts a URL to a DatadogLink
func URLToLink(logUrl string) (*api.DatadogLink, error) {
	baseUrl, queryParams, err := ParseURL(logUrl)
	if err != nil {
		return nil, err
	}

	link := &api.DatadogLink{
		APIVersion:  api.LinkGVK.GroupVersion().String(),
		Kind:        api.LinkGVK.Kind,
		BaseURL:     baseUrl,
		ExtraParams: map[string]string{},
	}

	queryParamMap := map[string]queryValHandler{
		"query":                         bindToString(&link.Query),
		"viz":                           bindToString(&link.VisualizeAs),
		"agg_m":                         bindToString(&link.GroupInto),
		"storage":                       bindToString(&link.Storage),
		"x_missing":                     bindToString(&link.Missing),
		"agg_m_source":                  bindToString(&link.Source),
		"agg_q":                         bindToString(&link.GroupBy),
		"clustering_pattern_field_path": bindToString(&link.ClusteringPatternFieldPath),
		"message_display":               bindToString(&link.MessageDisplay),
		"stream_sort":                   bindToString(&link.StreamSort),
		"agg_q_source":                  bindToString(&link.GroupBySource),
		"agg_t":                         bindToString(&link.AggType),
		"refresh_mode":                  bindToString(&link.RefreshMode),
		"from_ts":                       bindToString(&link.FromTS),
		"to_ts":                         bindToString(&link.ToTS),
		"fromUser":                      bindToString(&link.FromUser),
		"top_n":                         bindToInt(&link.TopN),
		"top_o":                         bindToString(&link.TopO),
		"live":                          bindToBool(&link.Live),
		"cols":                          bindToStringSlice(&link.Columns),
		"messageDisplay":                bindToString(&link.MessageDisplay),
	}

	for key, value := range queryParams {
		if targetFunc, found := queryParamMap[key]; found {
			targetFunc(value)
		} else {
			link.ExtraParams[key] = value[0]
		}
	}

	return link, nil
}
