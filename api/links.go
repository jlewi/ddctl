package api

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	LinkGVK = schema.FromAPIVersionAndKind(Group+"/"+Version, "DatadogLink")
)

// DatadogLink represents a link to a Datadog dashboard
type DatadogLink struct {
	APIVersion string   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// BaseURL is the base URL for links generated from this template
	BaseURL string `json:"baseURL,omitempty" yaml:"baseURL,omitempty"`
	// Query is the query to be used in the link
	Query string `json:"query,omitempty" yaml:"query,omitempty"`
	// VisualizeAs is the visualization to use for the link
	// This is the viz query key
	VisualizeAs string `json:"viz,omitempty" yaml:"viz,omitempty"`

	// GroupInto is the groupInto clause
	// This is the agg_m query key
	GroupInto string `json:"groupInto,omitempty" yaml:"groupInto,omitempty"`

	// Storage is the storage tier to query
	Storage string `json:"storage,omitempty" yaml:"storage,omitempty"`

	// Missing specifies the behavior for fields that maybe missing
	// This is the x_missing query key
	Missing string `json:"missing,omitempty" yaml:"missing,omitempty"`

	// TopN is the topN clause
	TopN int `json:"topN,omitempty" yaml:"topN,omitempty"`

	// Source is the value of the agg_m_source field
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// GroupBy is the value that we GroupBy
	// it is the value of the agg_q query key
	GroupBy string `json:"groupBy,omitempty" yaml:"groupBy,omitempty"`

	// ClusteringPatternFieldPath is the value of the clustering_pattern_field_path query key
	// It is how we cluster the data
	ClusteringPatternFieldPath string `json:"clusteringPatternFieldPath,omitempty" yaml:"clusteringPatternFieldPath,omitempty"`

	// MessageDisplay is the value of the messageDisplay query key
	MessageDisplay string `json:"messageDisplay,omitempty" yaml:"messageDisplay,omitempty"`

	// StreamSort is the value of the stream_sort query key
	StreamSort string `json:"streamSort,omitempty" yaml:"streamSort,omitempty"`

	// Live is the value of the live query key
	Live bool `json:"live,omitempty" yaml:"live,omitempty"`

	// TopO specifies the ordering of the top fields
	// This is the top_o query key
	// Descending means sort in descending order
	TopO string `json:"topO,omitempty" yaml:"topO,omitempty"`

	// GroupBySource is the value of the agg_q_source query key
	GroupBySource string `json:"groupBySource,omitempty" yaml:"groupBySource,omitempty"`

	// AggType is the aggregation type (e.g. count, avg, sum)
	// This is the agg_t query key
	AggType string `json:"aggType,omitempty" yaml:"aggType,omitempty"`

	// Columns is the columns to display
	// This is the cols query key
	Columns []string `json:"columns,omitempty" yaml:"columns,omitempty"`

	// RefreshMode is the value of the refresh_mode query key
	RefreshMode string `json:"refreshMode,omitempty" yaml:"refreshMode,omitempty"`

	// FromTS is the value of the from_ts query key
	// TODO(jeremy): Should we support relative times as we do in grafctl? e.g. with now
	FromTS string `json:"fromTS,omitempty" yaml:"fromTS,omitempty"`
	ToTS   string `json:"toTS,omitempty" yaml:"toTS,omitempty"`

	// Fromuser is the value of the fromUser field. According to chatGPT this is for
	// tracking purposes.
	FromUser string `json:"fromUser,omitempty" yaml:"fromUser,omitempty"`

	// ExtraParams is a map of extra parameters to include in the link
	ExtraParams map[string]string `json:"extraParams,omitempty" yaml:"extraParams,omitempty"`
}
