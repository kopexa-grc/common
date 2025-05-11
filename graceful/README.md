# graceful

`graceful` is a lightweight Go package for structured **graceful shutdowns** of services, servers, and background tasks. It enables reliable termination on system signals and allows fine-grained control over individual shutdown targets.

## Features

- Register multiple shutdown targets implementing `Shutdownable`
- Per-target shutdown timeouts
- Listens to OS signals (`SIGINT`, `SIGTERM`) and triggers shutdown
- Manual shutdown trigger support
- Safe for concurrent use
- Fully tested with support for error cases and timeouts
- Clean integration with [`zerolog`](https://github.com/rs/zerolog)

---

## 📦 Installation

```bash
go get github.com/kopexa-grc/kopexa/pkg/graceful
```

## 🚀 Usage
1. Implement the Shutdownable interface
```go
type MyServer struct {
	httpServer *http.Server
}

func (m *MyServer) Shutdown(ctx context.Context) error {
	return m.httpServer.Shutdown(ctx)
}
```

1. Integrate graceful.Closer in your main.go
```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/your-org/your-repo/graceful"
)

func main() {
	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		}),
	}

	closer := graceful.NewCloser()
	closer.Register("http-server", &MyServer{httpServer: srv}, 5*time.Second)

	triggerShutdown := closer.DetectShutdown()

	go func() {
		log.Info().Msg("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Optional: manually trigger shutdown from another goroutine
	// triggerShutdown()
}
```

## 🧠 Design Notes

- All shutdowns run concurrently and block until complete
- Each shutdown target can define its own timeout
- skipExit allows clean testing without terminating the process
- Thread safety is enforced using sync.Mutex and atomic
