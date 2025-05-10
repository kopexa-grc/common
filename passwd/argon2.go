// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package passwd

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
)

// ===========================================================================
// Derived Key Algorithm
// ===========================================================================

// Argon2Config holds the configuration for the Argon2id algorithm.
type Argon2Config struct {
	Time    uint32 // Number of iterations
	Memory  uint32 // Memory usage in KiB
	Threads uint8  // Number of parallel threads
	KeyLen  uint32 // Length of the derived key in bytes
	SaltLen uint32 // Length of the salt in bytes
}

// DefaultArgon2Config returns the recommended configuration for Argon2id.
// These values follow the recommendations from the Argon2 RFC draft.
func DefaultArgon2Config() Argon2Config {
	return Argon2Config{
		Time:    Argon2DefaultTime,
		Memory:  Argon2DefaultMemory,
		Threads: Argon2DefaultThreads,
		KeyLen:  Argon2DefaultKeyLen,
		SaltLen: Argon2DefaultSaltLen,
	}
}

// Argon2 constants for the derived key (dk) algorithm
const (
	dkAlg = "argon2id" // the derived key algorithm
)

// Argon2 variables for the derived key (dk) algorithm
var (
	dkParse = regexp.MustCompile(`^\$(?P<alg>[\w\d]+)\$v=(?P<ver>\d+)\$m=(?P<mem>\d+),t=(?P<time>\d+),p=(?P<procs>\d+)\$(?P<salt>[\+\/\=a-zA-Z0-9]+)\$(?P<key>[\+\/\=a-zA-Z0-9]+)$`)
)

// CreateDerivedKey creates an encoded derived key with a random hash for the password.
// It uses the default Argon2id configuration.
func CreateDerivedKey(password string) (string, error) {
	return CreateDerivedKeyWithConfig(password, DefaultArgon2Config())
}

// CreateDerivedKeyWithConfig creates an encoded derived key with a custom configuration.
func CreateDerivedKeyWithConfig(password string, config Argon2Config) (string, error) {
	if password == "" {
		return "", ErrCannotCreateDK
	}

	if config.Time == 0 || config.Memory == 0 || config.Threads == 0 || config.KeyLen == 0 || config.SaltLen == 0 {
		return "", ErrInvalidArgon2Config
	}

	salt := make([]byte, config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrCouldNotGenerate
	}

	dk := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)
	b64salt := base64.StdEncoding.EncodeToString(salt)
	b64dk := base64.StdEncoding.EncodeToString(dk)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		dkAlg, argon2.Version, config.Memory, config.Time, config.Threads, b64salt, b64dk), nil
}

// VerifyDerivedKey checks that the submitted password matches the derived key.
func VerifyDerivedKey(dk, password string) (bool, error) {
	if dk == "" || password == "" {
		return false, ErrUnableToVerify
	}

	dkb, salt, t, m, p, err := ParseDerivedKey(dk)
	if err != nil {
		return false, err
	}

	vdk := argon2.IDKey([]byte(password), salt, t, m, p, uint32(len(dkb))) // nolint:gosec

	return bytes.Equal(dkb, vdk), nil
}

// ParseDerivedKey returns the parts of the encoded derived key string.
func ParseDerivedKey(encoded string) (dk, salt []byte, time, memory uint32, threads uint8, err error) {
	if !dkParse.MatchString(encoded) {
		return nil, nil, 0, 0, 0, ErrCannotParseDK
	}

	parts := dkParse.FindStringSubmatch(encoded)

	if len(parts) != 8 { //nolint:mnd
		return nil, nil, 0, 0, 0, ErrCannotParseEncodedEK
	}

	// check the algorithm
	if parts[1] != dkAlg {
		return nil, nil, 0, 0, 0, newParseError("dkAlg", parts[1], dkAlg)
	}

	// check the version
	if version, err := strconv.Atoi(parts[2]); err != nil || version != argon2.Version {
		return nil, nil, 0, 0, 0, newParseError("version", parts[2], fmt.Sprintf("%d", argon2.Version))
	}

	var (
		time64    uint64
		memory64  uint64
		threads64 uint64
	)

	if memory64, err = strconv.ParseUint(parts[3], 10, 32); err != nil {
		return nil, nil, 0, 0, 0, newParseError("memory", parts[3], err.Error())
	}

	memory = uint32(memory64) // nolint:gosec

	if time64, err = strconv.ParseUint(parts[4], 10, 32); err != nil {
		return nil, nil, 0, 0, 0, newParseError("time", parts[4], err.Error())
	}

	time = uint32(time64) // nolint:gosec

	if threads64, err = strconv.ParseUint(parts[5], 10, 8); err != nil {
		return nil, nil, 0, 0, 0, newParseError("threads", parts[5], err.Error())
	}

	threads = uint8(threads64) // nolint:gosec

	if salt, err = base64.StdEncoding.DecodeString(parts[6]); err != nil {
		return nil, nil, 0, 0, 0, newParseError("salt", parts[6], err.Error())
	}

	if dk, err = base64.StdEncoding.DecodeString(parts[7]); err != nil {
		return nil, nil, 0, 0, 0, newParseError("dk", parts[7], err.Error())
	}

	return dk, salt, time, memory, threads, nil
}

// IsDerivedKey checks if a string is a valid derived key.
func IsDerivedKey(s string) bool {
	return dkParse.MatchString(s)
}

// GetDerivedKeyConfig returns the configuration used to create a derived key.
func GetDerivedKeyConfig(dk string) (Argon2Config, error) {
	_, _, time, memory, threads, err := ParseDerivedKey(dk)
	if err != nil {
		return Argon2Config{}, err
	}

	return Argon2Config{
		Time:    time,
		Memory:  memory,
		Threads: threads,
		KeyLen:  Argon2DefaultKeyLen,
		SaltLen: Argon2DefaultSaltLen,
	}, nil
}
