# Password Package

This package provides password validation and hashing functionality using Argon2id, the winner of the 2015 Password Hashing Competition.

## What is a Derived Key?

A derived key is a cryptographically secure value generated from a password using a key derivation function (KDF) such as Argon2id. Instead of storing the password itself, only the derived key (together with its parameters and salt) is stored. This process makes it computationally infeasible for attackers to recover the original password, even if the derived key is compromised. The derived key is unique for each password and salt combination, and is suitable for secure password storage and cryptographic operations (e.g., as a key for encryption).

## Features

- Password strength evaluation
- Argon2id password hashing
- Common password detection
- Leetspeak detection
- Personal information detection

## Usage

### Password Strength Evaluation

```go
feedback := passwd.Evaluate("your-password")
if feedback.Level == passwd.Rejected {
    // Handle rejected password
    for _, msg := range feedback.Messages {
        fmt.Println(msg)
    }
}
```

### Password Hashing with Argon2id

```go
// Create a derived key from a password
dk, err := passwd.CreateDerivedKey("your-password")
if err != nil {
    // Handle error
}

// Verify a password against a derived key
valid, err := passwd.VerifyDerivedKey(dk, "your-password")
if err != nil {
    // Handle error
}
if !valid {
    // Handle invalid password
}
```

## Security

The package uses Argon2id with the following parameters:
- Memory: 64MB (configurable)
- Time: 1 iteration (configurable)
- Parallelism: 2 threads (configurable)
- Salt length: 16 bytes
- Key length: 32 bytes (AES-256)

These parameters follow the recommendations from the Argon2 RFC draft and can be adjusted based on your security requirements.

## Best Practices

1. Always use `CreateDerivedKey` to hash passwords before storage
2. Use `VerifyDerivedKey` to verify passwords
3. Use `Evaluate` to check password strength before hashing
4. Consider using `EvaluateWithContext` to check for personal information

## References

- [Argon2 RFC Draft](https://datatracker.ietf.org/doc/draft-irtf-cfrg-argon2/)
- [Password Hashing Competition](https://password-hashing.net/)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html) 