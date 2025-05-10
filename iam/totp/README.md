# TOTP Package

A secure and enterprise-ready Time-based One-Time Password (TOTP) implementation for Go applications.

## Motivation

In today's digital landscape, securing user accounts is more critical than ever. While passwords remain a fundamental security measure, they are increasingly vulnerable to various attacks. Multi-factor authentication (MFA) has become essential for protecting user accounts and sensitive data.

This package provides a robust implementation of TOTP-based MFA, offering:

- **Enhanced Security**: Adds an additional layer of protection beyond passwords
- **Industry Standard**: Implements RFC 6238 (TOTP) for compatibility with authenticator apps
- **Flexible Delivery**: Supports both TOTP and OTP (email/SMS) authentication methods
- **Enterprise Ready**: Includes features like recovery codes and secure secret management

## Features

- **TOTP Generation**: Generate and validate time-based one-time passwords
- **QR Code Support**: Generate QR codes for easy setup with authenticator apps
- **Multiple Delivery Methods**: Support for email and SMS OTP delivery
- **Recovery Codes**: Generate secure recovery codes for account recovery
- **Secure Storage**: Encrypted secret storage with versioning support
- **NATS Integration**: Distributed storage backend using NATS KV
- **High Test Coverage**: Comprehensive test suite with >90% coverage

## User Workflow: From Activation to Recovery

This section describes the complete user journey for TOTP-based 2FA, from activation in the UI to account recovery using recovery codes.

### 1. Activation (User enables TOTP in the UI)
- **User Action:** Navigates to security settings and chooses to enable 2FA (TOTP).
- **Backend:**
  - Generates a new TOTP secret for the user (`otp.TOTPSecret(user)`).
  - Stores the encrypted secret in the user's profile.
  - Generates a QR code URL (`otp.TOTPQRString(user)`) and a set of recovery codes (`otp.GenerateRecoveryCodes()`).
- **UI:**
  - Displays the QR code for scanning with an authenticator app (e.g., Google Authenticator).
  - Shows the recovery codes **only once** and instructs the user to save them securely.

### 2. Verification (User confirms setup)
- **User Action:** Scans the QR code and enters the generated TOTP code in the UI.
- **Backend:**
  - Validates the code (`otp.ValidateTOTP(ctx, user, code)`).
  - If valid, marks TOTP as enabled for the user.
- **UI:**
  - Confirms successful setup or shows an error if the code is invalid.

### 3. Login with TOTP
- **User Action:** Logs in with username and password, then is prompted for a TOTP code.
- **Backend:**
  - Validates the TOTP code (`otp.ValidateTOTP(ctx, user, code)`).
- **UI:**
  - Grants access if the code is valid, otherwise shows an error.

### 4. Recovery (User lost access to TOTP device)
- **User Action:** Clicks "Can't access your authenticator?" and chooses to use a recovery code.
- **UI:**
  - Prompts for a recovery code.
- **Backend:**
  - Checks if the code matches a stored (hashed) recovery code and is unused.
  - If valid, marks the code as used and allows the user to log in.
  - Optionally, prompts the user to re-enroll TOTP after login.
- **UI:**
  - Informs the user if the code was accepted or rejected.
  - After successful recovery, instructs the user to set up TOTP again for continued protection.

### 5. Best Practices
- Always show recovery codes only once and never store them in plaintext.
- Invalidate used recovery codes immediately.
- Encourage users to store recovery codes in a secure place (e.g., password manager).
- After recovery, require the user to re-enable TOTP for continued account security.

## Persistence and Storage

This section explains how TOTP secrets and recovery codes are persisted and stored securely.

### TOTP Secrets
- **Storage:** TOTP secrets are encrypted and stored in the user's profile. The package uses NATS KV as a distributed storage backend.
- **Encryption:** Secrets are encrypted using AES-CTR before being stored, ensuring they are not accessible in plaintext.
- **Versioning:** The package supports versioning of secrets, allowing for secure updates and rotations.

### Recovery Codes
- **Storage:** Recovery codes are generated securely and should be stored by the user in a safe place (e.g., password manager or printed and kept offline).
- **Backend Storage:** In the backend, recovery codes are hashed and stored in a database. After use, they are marked as used or deleted to prevent reuse.

### NATS KV Integration
- **Configuration:** The package is configured to use NATS KV for storing TOTP secrets. This provides a distributed and scalable storage solution.
- **Example Configuration:**
  ```go
  nc, _ := nats.Connect(nats.DefaultURL)
  js, _ := nc.JetStream()
  kv, _ := js.CreateKeyValue(&nats.KeyValueConfig{
      Bucket: "totp",
  })
  store := totp.NewNATSStore(kv)
  ```

## Installation

```bash
go get github.com/kopexa-grc/common/iam/totp
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/kopexa-grc/common/iam/totp"
    "github.com/nats-io/nats.go"
)

func main() {
    // Create NATS connection
    nc, _ := nats.Connect(nats.DefaultURL)
    js, _ := nc.JetStream()
    kv, _ := js.CreateKeyValue(&nats.KeyValueConfig{
        Bucket: "totp",
    })

    // Create store
    store := totp.NewNATSStore(kv)

    // Create OTP manager
    otp := totp.New(store,
        totp.WithIssuer("MyApp"),
        totp.WithSecret(totp.Secret{
            Version: 1,
            Key:     []byte("your-secret-key"),
        }),
    )

    // Create user
    user := &totp.User{
        ID:    "user123",
        Email: sql.NullString{String: "user@example.com", Valid: true},
    }

    // Generate TOTP secret
    secret, _ := otp.TOTPSecret(user)
    user.TOTPSecret = secret

    // Generate QR code URL
    qrURL, _ := otp.TOTPQRString(user)
    fmt.Println("Scan this QR code with your authenticator app:", qrURL)

    // Validate TOTP code
    code := "123456" // Code from authenticator app
    err := otp.ValidateTOTP(context.Background(), user, code)
    if err != nil {
        fmt.Println("Invalid code")
    } else {
        fmt.Println("Valid code")
    }
}
```

## User Workflow Code Examples

### 1. Activation (User enables TOTP in the UI)
```go
// Generate TOTP secret and QR code URL
secret, err := otp.TOTPSecret(user)
if err != nil {
    // Handle error
}
user.TOTPSecret = secret

qrURL, err := otp.TOTPQRString(user)
if err != nil {
    // Handle error
}
fmt.Println("Scan this QR code with your authenticator app:", qrURL)

// Generate recovery codes
codes := otp.GenerateRecoveryCodes()
for _, code := range codes {
    fmt.Println("Your recovery code:", code)
    // Store a hash of the code in your database for later verification
    // Example: storeRecoveryCodeHash(userID, hash(code))
}
```

### 2. Verification (User confirms setup)
```go
// Validate TOTP code
code := "123456" // Code from authenticator app
err := otp.ValidateTOTP(context.Background(), user, code)
if err != nil {
    fmt.Println("Invalid code")
} else {
    fmt.Println("Valid code")
    // Mark TOTP as enabled for the user
}
```

### 3. Login with TOTP
```go
// Validate TOTP code
code := "123456" // Code from authenticator app
err := otp.ValidateTOTP(context.Background(), user, code)
if err != nil {
    fmt.Println("Invalid code")
} else {
    fmt.Println("Valid code")
    // Grant access to the user
}
```

### 4. Recovery (User lost access to TOTP device)
```go
// When a user submits a recovery code for login:
inputCode := "user-input-code"
if isValidRecoveryCode(userID, inputCode) {
    // Mark the code as used or remove it from the database
    // Example: markRecoveryCodeUsed(userID, inputCode)
    fmt.Println("Recovery code accepted. Please reset your TOTP device.")
} else {
    fmt.Println("Invalid or already used recovery code.")
}
```

## Security Considerations

- All secrets are encrypted using AES-CTR
- Random values are generated using crypto/rand
- TOTP secrets are stored in Base32 format
- Recovery codes are generated using a secure random source
- All cryptographic operations are performed in constant time

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the BUSL-1.1 License - see the LICENSE file for details. 