package api

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	TraceGVK = schema.FromAPIVersionAndKind(Group+"/"+Version, "DatadogTrace")
)

// DatadogTrace represents a link to the APM traces in Datadog
type DatadogTrace struct {
	APIVersion string   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// BaseURL is the base URL for links generated from this template
	BaseURL string `json:"baseURL,omitempty" yaml:"baseURL,omitempty"`

	TraceID          string `json:"traceID,omitempty" yaml:"traceID,omitempty"`
	SpanID           string `json:"spanID,omitempty" yaml:"spanID,omitempty"`
	GraphType        string `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	PanelTab         string `json:"panelTab,omitempty" yaml:"panelTab,omitempty"`
	ShouldShowLegend bool   `json:"shouldShowLegend,omitempty" yaml:"shouldShowLegend,omitempty"`
	Sort             string `json:"sort,omitempty" yaml:"sort,omitempty"`
	TimeHint         string `json:"timeHint,omitempty" yaml:"timeHint,omitempty"`

	// ExtraParams is a map of extra parameters to include in the link
	ExtraParams map[string]string `json:"extraParams,omitempty" yaml:"extraParams,omitempty"`
}
