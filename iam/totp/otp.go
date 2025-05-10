// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package totp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

// otpNATS defines the interface for NATS operations
type otpNATS interface {
	// Get retrieves a value from the store
	Get(ctx context.Context, key string) ([]byte, error)
	// Set stores a value in the store
	Set(ctx context.Context, key string, value []byte) error
	// Delete removes a value from the store
	Delete(ctx context.Context, key string) error
}

// OTP implements the Manager interface
type OTP struct {
	codeLength         int
	issuer             string
	recoveryCodeCount  int
	recoveryCodeLength int
	secrets            []Secret
	db                 otpNATS
}

// New creates a new OTP instance
func New(db otpNATS, opts ...ConfigOption) *OTP {
	otp := &OTP{
		codeLength:         DefaultLength,
		recoveryCodeCount:  DefaultRecoveryCodeCount,
		recoveryCodeLength: DefaultRecoveryCodeLength,
		db:                 db,
	}

	for _, opt := range opts {
		opt(otp)
	}

	return otp
}

// TOTPQRString returns a URL string used for TOTP code generation
func (o *OTP) TOTPQRString(u *User) (string, error) {
	secret, err := o.TOTPSecret(u)
	if err != nil {
		return "", fmt.Errorf("failed to get secret for QR: %w", err)
	}

	decrypted, err := o.TOTPDecryptedSecret(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Generate QR code URL
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      o.issuer,
		AccountName: u.DefaultName(),
		Secret:      []byte(decrypted),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return key.URL(), nil
}

// TOTPDecryptedSecret decrypts a TOTP secret
func (o *OTP) TOTPDecryptedSecret(secret string) (string, error) {
	if len(o.secrets) == 0 {
		return "", ErrNoSecretKey
	}

	latestSecret := o.secrets[len(o.secrets)-1]

	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", ErrCannotDecodeSecret
	}

	block, err := aes.NewCipher([]byte(latestSecret.Key))
	if err != nil {
		return "", ErrFailedToCreateCipherBlock
	}

	if len(decoded) < aes.BlockSize {
		return "", ErrCipherTextTooShort
	}

	iv := decoded[:aes.BlockSize]
	ciphertext := decoded[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

// TOTPSecret creates a TOTP secret for code generation
func (o *OTP) TOTPSecret(_ *User) (string, error) {
	if len(o.secrets) == 0 {
		return "", ErrNoSecretKey
	}

	// Get the latest secret version
	latestSecret := o.secrets[len(o.secrets)-1]

	// Generate a valid base32 secret for TOTP
	base32Secret := make([]byte, Base32SecretLength)
	if _, err := crand.Read(base32Secret); err != nil {
		return "", ErrFailedToGenerateSecret
	}

	secretString := base32.StdEncoding.EncodeToString(base32Secret)

	// Encrypt the base32 secret
	block, err := aes.NewCipher([]byte(latestSecret.Key))
	if err != nil {
		return "", ErrFailedToCreateCipherBlock
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := crand.Read(iv); err != nil {
		return "", ErrFailedToGenerateSecret
	}

	stream := cipher.NewCTR(block, iv)
	plaintext := []byte(secretString)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	combined := append([]byte{}, iv...)
	combined = append(combined, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(combined)

	return encoded, nil
}

// OTPCode creates a random OTP code and hash
func (o *OTP) OTPCode(address string, method DeliveryMethod) (code, hash string, err error) {
	// Generate random code
	code = generateRandomString(o.codeLength, NumericCode)

	// Create hash
	h := sha256.New()
	h.Write([]byte(code))
	hash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Store hash in NATS
	key := fmt.Sprintf("otp:%s:%s", method, address)
	value := Hash{
		Hash:      hash,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(value)
	if err != nil {
		return "", "", ErrCannotHashOTPString
	}

	if err := o.db.Set(context.Background(), key, data); err != nil {
		return "", "", ErrCannotHashOTPString
	}

	return code, hash, nil
}

// ValidateOTP checks if a User's email/SMS delivered OTP code is valid
func (o *OTP) ValidateOTP(code, hash string) error {
	// Create hash of provided code
	h := sha256.New()
	h.Write([]byte(code))
	codeHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Compare hashes
	if codeHash != hash {
		return ErrInvalidCode
	}

	return nil
}

// ValidateTOTP checks if a User's TOTP code is valid
func (o *OTP) ValidateTOTP(_ context.Context, user *User, code string) error {
	secret, err := o.TOTPDecryptedSecret(user.TOTPSecret)
	if err != nil {
		return ErrFailedToValidateCode
	}

	valid := totp.Validate(code, secret)
	if !valid {
		return ErrInvalidCode
	}

	return nil
}

// GenerateRecoveryCodes creates a set of recovery codes for a user
func (o *OTP) GenerateRecoveryCodes() []string {
	codes := make([]string, o.recoveryCodeCount)
	for i := 0; i < o.recoveryCodeCount; i++ {
		codes[i] = generateRandomString(o.recoveryCodeLength, AlphanumericCode)
	}

	return codes
}

// generateRandomString generates a random string of the specified length using the provided charset
func generateRandomString(length int, charset string) string {
	b := make([]byte, length)
	if _, err := crand.Read(b); err != nil {
		panic("crypto/rand failed")
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
