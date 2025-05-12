// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// noopLogger is a logger that does nothing
type noopLogger struct{}

func (l *noopLogger) Printf(_ string, _ ...interface{}) {}

// PostgresContainer represents a PostgreSQL test container
type PostgresContainer struct {
	container *postgres.PostgresContainer
	config    *PostgresConfig
}

// PostgresConfig holds the configuration for the PostgreSQL container
type PostgresConfig struct {
	// Image is the PostgreSQL image to use
	Image string
	// Database is the name of the database to create
	Database string
	// Username is the username to use for authentication
	Username string
	// Password is the password to use for authentication
	Password string
}

// DefaultPostgresConfig returns a default PostgreSQL configuration
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Image:    "postgres:16-alpine",
		Database: "test",
		Username: "test",
		Password: "test",
	}
}

// PostgresOption is a function that modifies a PostgresConfig
type PostgresOption func(*PostgresConfig)

// WithImage sets the image name
func WithImage(image string) PostgresOption {
	return func(c *PostgresConfig) {
		c.Image = image
	}
}

// WithDatabase sets the database name
func WithDatabase(database string) PostgresOption {
	return func(c *PostgresConfig) {
		c.Database = database
	}
}

// WithUsername sets the username
func WithUsername(username string) PostgresOption {
	return func(c *PostgresConfig) {
		c.Username = username
	}
}

// WithPassword sets the password
func WithPassword(password string) PostgresOption {
	return func(c *PostgresConfig) {
		c.Password = password
	}
}

// NewPostgresContainer creates a new PostgreSQL test container
func NewPostgresContainer(ctx context.Context, opts ...PostgresOption) (*PostgresContainer, error) {
	config := DefaultPostgresConfig()
	for _, opt := range opts {
		opt(config)
	}

	container, err := postgres.Run(ctx,
		config.Image,
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.Username),
		postgres.WithPassword(config.Password),
		testcontainers.WithLogger(&noopLogger{}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return &PostgresContainer{
		container: container,
		config:    config,
	}, nil
}

// GetDSN returns the Data Source Name for the PostgreSQL container
func (c *PostgresContainer) GetDSN(ctx context.Context) (string, error) {
	host, err := c.container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := c.container.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("failed to get container port: %w", err)
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), c.config.Username, c.config.Password, c.config.Database), nil
}

// Cleanup stops and removes the PostgreSQL container
func (c *PostgresContainer) Cleanup(t *testing.T) {
	testcontainers.CleanupContainer(t, c.container)
}
