# Sessions Package

The `sessions` package provides a type-safe session management system for web applications. It supports generic types for session values and includes features like:

- Secure session ID generation using UUID v4
- Configurable session storage
- Cookie-based session management
- Thread-safe operations

## Features

- **Type Safety**: Uses Go Generics for type-safe session values
- **Security**: 
  - AES-GCM encryption for session data
  - Secure cookie configuration (HttpOnly, Secure, SameSite)
  - UUID v4 for session IDs
- **Flexibility**: 
  - Configurable session storage
  - Customizable cookie settings
  - Extensible store implementations

## Installation

```bash
go get github.com/kopexa-grc/common/iam/sessions
```

## Usage

### Basic Usage

```go
package main

import (
    "net/http"
    
    "github.com/kopexa-grc/common/iam/sessions"
    "github.com/kopexa-grc/common/iam/sessions/cookie"
)

func main() {
    // Create store with options
    store := cookie.NewStore[string](
        cookie.WithSigningKey("your-32-byte-signing-key-here"),
        cookie.WithEncryptionKey("your-32-byte-encryption-key"),
        cookie.WithMaxAge(3600), // 1 hour
        cookie.WithSecure(true),
        cookie.WithHTTPOnly(true),
        cookie.WithSameSite(sessions.CookieSameSiteLax),
    )

    // Für Entwicklungs- und Testumgebungen:
    devStore := cookie.NewStore[string](
        cookie.WithSigningKey("your-32-byte-signing-key-here"),
        cookie.WithEncryptionKey("your-32-byte-encryption-key"),
        cookie.WithDevMode(true),  // Aktiviert Entwicklungsmodus
        cookie.WithSecure(false),  // Erlaubt HTTP
        cookie.WithHTTPOnly(false), // Erlaubt JavaScript-Zugriff
        cookie.WithDomain("localhost"), // Lokale Entwicklung
    )

    // Create session
    session := sessions.NewSession(store, "user_session")

    // Set values
    session.Set("user_id", "123")
    session.Set("role", "admin")

    // Save session
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if err := session.Save(w); err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    })

    // Load session
    http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
        loaded, err := store.Load(r, "user_session")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        userID := loaded.Get("user_id")
        role := loaded.Get("role")
        // ...
    })
}
```

### Custom Types

```go
type UserData struct {
    ID       string
    Username string
    Role     string
}

// Create store with custom type
store := cookie.NewStore[UserData](
    cookie.WithSigningKey("your-32-byte-signing-key-here"),
    cookie.WithEncryptionKey("your-32-byte-encryption-key"),
)

// Create session with custom type
session := sessions.NewSession(store, "user_session")

// Set structured data
session.Set("user", UserData{
    ID:       "123",
    Username: "john.doe",
    Role:     "admin",
})

// Get data
user, ok := session.GetOk("user")
if ok {
    fmt.Printf("User: %s, Role: %s\n", user.Username, user.Role)
}
```

## Security Notes

1. **Keys**: 
   - Use strong, random keys for `SigningKey` and `EncryptionKey`
   - Minimum length: 32 bytes
   - Keep keys secret

2. **Cookie Settings**:
   - `Secure: true` for HTTPS
   - `HttpOnly: true` against XSS
   - `SameSite: "Lax"` or `"Strict"` against CSRF
   - **Domain for Subdomains**: If you want sessions to be shared across subdomains (e.g. `auth.kopexa.com`, `app.kopexa.com`, `console.kopexa.com`), set the cookie domain explicitly to `.kopexa.com` when creating the store. This ensures all subdomains can access the session cookie securely.

3. **Development Mode**:
   - Use `WithDevMode(true)` for local development
   - Allows insecure settings (HTTP, JavaScript access)
   - Never use in production!
   - Useful for local testing and development

4. **Session Timeout**:
   - Set an appropriate `MaxAge`
   - Implement server-side session validation

5. **Configuration Validation**:
   - The store will refuse to start with insecure or invalid settings (e.g. short keys, missing Secure/HTTPOnly, invalid SameSite, MaxAge <= 0).
   - This ensures you cannot accidentally deploy an insecure session configuration.

## Best Practices

1. **Session Size**:
   - Keep session data minimal
   - Store only necessary information

2. **Error Handling**:
   - Always check for `ErrSessionExpired`
   - Implement appropriate error handling

3. **Security**:
   - Regular session rotation
   - CSRF protection implementation
   - Secure cookie configuration

4. **Subdomain Support**:
   - Always use a leading dot in the domain (e.g. `.kopexa.com`)
   - Ensure all subdomains use HTTPS
   - Consider using `SameSite=Lax` for better subdomain compatibility
   - Test session sharing across all subdomains

## Store Implementations

### Cookie Store

The cookie store implementation uses HTTP cookies to store session data. It provides:

- Secure cookie configuration
- AES-GCM encryption
- Configurable session duration
- Domain support for subdomains

### NATS Store

The NATS store implementation uses NATS JetStream for distributed session storage. It provides:

- Distributed session storage
- Session metadata tracking (IP, User-Agent)
- Active session monitoring
- Automatic session expiration
- High availability and scalability

#### Usage

```go
package main

import (
    "net/http"
    
    "github.com/kopexa-grc/common/iam/sessions"
    "github.com/kopexa-grc/common/iam/sessions/nats"
)

func main() {
    // Create NATS store
    store := nats.NewStore[string](
        nats.WithServerURL("nats://localhost:4222"),
        nats.WithBucketName("sessions"),
        nats.WithMaxAge(3600), // 1 hour
    )

    // Create session
    session := sessions.NewSession(store, "user_session")

    // Set values
    session.Set("user_id", "123")
    session.Set("role", "admin")

    // Save session
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if err := session.Save(w); err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    })

    // Load session
    http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
        loaded, err := store.Load(r, "user_session")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        userID := loaded.Get("user_id")
        role := loaded.Get("role")
        // ...
    })

    // Get active sessions
    activeSessions, err := store.GetActiveSessions()
    if err != nil {
        // Handle error
    }

    // Process active sessions
    for _, session := range activeSessions {
        fmt.Printf("Session ID: %s\n", session.Session.ID)
        fmt.Printf("IP: %s\n", session.IP)
        fmt.Printf("User-Agent: %s\n", session.UserAgent)
        fmt.Printf("Last Seen: %s\n", session.LastSeen)
    }
}
```

#### Security Features

1. **Session Metadata**:
   - IP address tracking
   - User-Agent tracking
   - Last seen timestamp
   - Automatic session expiration

2. **Session Monitoring**:
   - Active session listing
   - Session metadata inspection
   - Suspicious activity detection

3. **Best Practices**:
   - Use TLS for NATS connections
   - Configure appropriate session timeouts
   - Monitor session activity
   - Implement session rotation

## License

BUSL-1.1 