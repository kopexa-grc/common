# logger

Enterprise-ready, environment-aware logging for Go applications.

## Features
- Automatic log format selection (console, JSON, GCP JSON)
- Environment detection (Kubernetes, Docker, local)
- Structured logging with [zerolog](https://github.com/rs/zerolog)
- Color support for console output
- Buffered logging with pause/resume
- Configurable log levels via environment variables
- Version tagging in all log entries
- Thread-safe and high-performance

## Usage

```go
import "github.com/kopexa-grc/common/logger"

func main() {
    // Configure logger with version string
    logger.Configure("v1.2.3")
    // ... your app code ...
}
```

## Environment Integration
- **Kubernetes:** Uses GCP-compatible JSON logging for best integration with cloud logging.
- **Docker:** Uses plain JSON logging for easy log aggregation.
- **Local:** Uses colorized, human-friendly console output.

Detection is automatic based on environment variables and file presence.

## Log Level Control
- Set `DEBUG=true` or `TRACE=true` in the environment to increase verbosity.
- Default log level is `info`.

## Example
```go
logger.Configure("v1.0.0")
logger.Set("debug")
log := logger.LogOutputWriter
log.Write([]byte("hello world\n"))
```

## Testing
All core components are covered by unit tests:
- Buffering
- Color output
- Environment detection

## License
Business Source License 1.1 (BUSL-1.1) 