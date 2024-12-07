# Datadog CLI

CLI For Datadog.

This CLI is intended to make it easy to use Datadog with [foyle.io](https://foyle.io/) and 
[Runme.Dev](https://runme.dev/).

You can use `ddctl` to easily generate links to Datadog logs. Using a CLI to generate links allows you to use Foyle
to generate the links automatically. To use it

```bash
ddctl logs querytourl --base-url=https://acme.datadoghq.com --query="query: service:foyle" --duration=1h --end-time="2024-12-06 15:20 PST" --open=true
```

To avoid having to specify the `--base-url` flag every time, you can set it in your configuration

```bash

ddctl config set baseURL https://acme.datadoghq.com
```
