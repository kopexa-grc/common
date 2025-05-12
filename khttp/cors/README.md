# CORS Package

The `cors` package provides a standardized CORS (Cross-Origin Resource Sharing) configuration for Kopexa APIs. It implements a secure and flexible CORS middleware that can be easily integrated into any HTTP server.

## Features

- Pre-configured CORS settings for Kopexa APIs
- Support for common HTTP methods (GET, POST, PUT, PATCH, DELETE, etc.)
- Standard security headers
- Configurable additional headers
- Built-in middleware for easy integration

## Usage

```go
import "github.com/kopexa-grc/common/khttp/cors"

// Create a new router
router := chi.NewRouter()

// Apply the CORS middleware
router.Use(cors.Middleware)

// Add your routes
router.Get("/api/v1/resource", handler)
```

## Customization

The package provides helper functions to extend the default configuration:

```go
// Allow an additional header
cors.AllowExtraHeader("X-Custom-Header")

// Expose an additional header
cors.ExposeExtraHeader("X-Custom-Response-Header")
```

## Default Configuration

The default configuration includes:

- All common HTTP methods
- Standard security headers
- Common REST response headers
- Credentials support
- 6-hour cache duration

## Best Practices

1. **Extending Headers:**
   - Use `AllowExtraHeader` for service-specific headers
   - Use `ExposeExtraHeader` for custom response headers
   - Consider contributing common headers back to this package

2. **Security:**
   - The default configuration is secure for most use cases
   - Review the allowed headers and methods for your specific needs
   - Consider implementing a custom `AllowOriginFunc` for production

3. **Performance:**
   - The middleware is lightweight and has minimal performance impact
   - The 6-hour cache duration reduces preflight requests 