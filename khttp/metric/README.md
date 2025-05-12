# Metric Package

The `metric` package provides standardized Prometheus metrics for HTTP services in the Kopexa platform. It implements a comprehensive set of metrics for monitoring HTTP request performance, sizes, and patterns.

## Features

- Pre-configured Prometheus metrics for HTTP services
- Standard metrics for request duration, size, and patterns
- Built-in middleware for easy integration
- Global registry for centralized metric collection
- Support for custom metric labels and buckets

## Usage

```go
import "github.com/kopexa-grc/common/khttp/metric"

// Create a new router
router := chi.NewRouter()

// Apply the metrics middleware
router.Use(metric.Middleware)

// Register the metrics endpoint
router.Get("/metrics", metric.Handler.ServeHTTP)

// Add your routes
router.Get("/api/v1/resource", handler)
```

## Available Metrics

The package provides the following metrics:

- `kopexa_http_in_flight_requests`: Number of HTTP requests currently being processed
- `kopexa_http_api_requests_total`: Total number of HTTP requests by method and status code
- `kopexa_http_request_duration_seconds`: Request duration histogram
- `kopexa_http_time_to_first_header_duration_seconds`: Time to first header histogram
- `kopexa_http_response_size_bytes`: Response size histogram
- `kopexa_http_request_size_bytes`: Request size histogram

## Customization

The package uses the Kopexa Prometheus namespace and provides configurable buckets for histograms:

```go
// Duration buckets for latency metrics
DurationBuckets = []float64{0.001, 0.01, 0.015, .25, .5, 1, 1.5, 2.5, 5, 10}

// Size buckets for request/response size metrics
SizeBuckets = []float64{200, 500, 900, 1500, 2000, 10000}
```

## Best Practices

1. **Metric Names:**
   - All metrics use the `kopexa` namespace
   - Metric names are descriptive and follow Prometheus naming conventions
   - Labels are standardized across all metrics

2. **Performance:**
   - The middleware is optimized for minimal overhead
   - Metrics are collected asynchronously
   - Histogram buckets are carefully chosen for typical HTTP workloads

3. **Integration:**
   - Use the global registry for most cases
   - Create a custom registry only if you need to isolate metrics
   - Always expose metrics on a dedicated `/metrics` endpoint 