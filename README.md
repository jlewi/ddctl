# Datadog CLI

CLI For Datadog.

This CLI is intended to make it easy to use Datadog with [foyle.io](https://foyle.io/) and 
[Runme.Dev](https://runme.dev/).

You can use `ddctl` to easily generate links to Datadog logs. Using a CLI to generate links allows you to use Foyle
to generate the links automatically.

## Quickstart

### Install

1. Download the latest release from the [releases page](https://github.com/jlewi/ddctl/releases)

### Parse An Existing Explore Link

The easiest way to understand the query syntax for a dashboard is by opening up an existing dashboard in Datadog. You
can then copy the link for it.

Once you have the URL you can parse it using `ddctl`

```
ddctl links parse --url=${URL}
```

This will output a `ddctl` resource like the one below.

```
apiVersion: datadog.foyle.io/v1alpha1
kind: DatadogLink
baseURL: https://acme.datadoghq.com
query: RequestLoggingMiddleware env:prod service:feserver* @handler_module:*bert* -@http.method:GET -@http.method:HEAD status:error -@handler_module:*laxmod* -@handler:*laxmod*
viz: pattern
groupInto: count
storage: flex_tier
missing: "true"
topN: 10
source: base
groupBy: status
clusteringPatternFieldPath: message
messageDisplay: inline
streamSort: desc
topO: top
groupBySource: base
aggType: count
columns:
    - host
    - service
refreshMode: paused
fromTS: "1736927929003"
toTS: "1736949529003"
fromUser: "true"
```

### Generate an A Link

You can generate a link for an view by specifying a DatadogLink resource that contains the query parameters for your 
query e.g

```bash
cat <<'EOF' > /tmp/query.yaml
apiVersion: grafctl.foyle.io/v1alpha1
apiVersion: datadog.foyle.io/v1alpha1
kind: DatadogLink
baseURL: https://acme.datadoghq.com
query: RequestLoggingMiddleware env:prod service:feserver* @handler_module:*bert* -@http.method:GET -@http.method:HEAD status:error -@handler_module:*laxmod* -@handler:*laxmod*
viz: pattern
groupInto: count
storage: flex_tier
missing: "true"
topN: 10
source: base
groupBy: status
clusteringPatternFieldPath: message
messageDisplay: inline
streamSort: desc
topO: top
groupBySource: base
aggType: count
columns:
    - host
    - service
refreshMode: paused
fromTS: "now-5m"
toTS: "now"
fromUser: "true"
EOF
ddctl links build -f=/tmp/query.yaml --open
```

**Important** Note that EOF is enclosed in single quotes. This prevents escaping and shell interpolation. Without this
shell escaping and interpolation can prevent the query from being encoded correctly.


## Timestamps

You can use Grafana style time expressions e.g. "now-5m" for `FromTS` and `ToTS`. `ddctl`
automatically converts this into the unix epoch timestamps that Datadog expects.

