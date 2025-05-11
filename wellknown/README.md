# Wellknown Package

The `wellknown` package contains constants and configurations that are used throughout the Kopexa platform. These values are standardized and should be used consistently across all services.

## Contents

### Prometheus

- `PrometheusNamespaceKopexa`: The namespace for all Prometheus metrics in the Kopexa platform
  - Usage: All metrics should use this namespace, e.g. `kopexa_http_requests_total`
  - Format: `kopexa_<service>_<metric_name>`

## Best Practices

1. **Adding New Constants:**
   - Constants should have a clear, descriptive name
   - Each constant should be documented with a comment
   - Constants should be grouped in the appropriate file (e.g., `prometheus.go` for Prometheus-related constants)

2. **Usage:**
   - Always use constants from this package instead of defining magic strings
   - When creating new services or components that need these values, import the constants from this package

3. **Changes:**
   - Changes to these constants should be made with caution as they can affect the entire platform
   - When making changes, all affected services must be updated

## Examples

```go
import "github.com/kopexa-grc/common/wellknown"

// Prometheus metric with correct namespace
prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: wellknown.PrometheusNamespaceKopexa,
        Name:      "http_requests_total",
        Help:      "Total number of HTTP requests",
    },
    []string{"method", "path", "status"},
)
``` 