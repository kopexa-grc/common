// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"os"
	"testing"
)

func TestIsKubernetes(t *testing.T) {
	// Backup and restore environment
	origEnv := os.Getenv("KUBERNETES_SERVICE_HOST")
	defer os.Setenv("KUBERNETES_SERVICE_HOST", origEnv)

	// Should be false by default
	os.Unsetenv("KUBERNETES_SERVICE_HOST")

	if isKubernetes() {
		t.Error("expected false when not in k8s")
	}

	// Should be true if env var is set
	os.Setenv("KUBERNETES_SERVICE_HOST", "1.2.3.4")

	if !isKubernetes() {
		t.Error("expected true when KUBERNETES_SERVICE_HOST is set")
	}
}

func TestIsDocker(t *testing.T) {
	// Backup and restore environment
	origEnv := os.Getenv("RUNNING_IN_DOCKER")
	defer os.Setenv("RUNNING_IN_DOCKER", origEnv)

	os.Unsetenv("RUNNING_IN_DOCKER")

	if isDocker() {
		t.Error("expected false when not in docker")
	}

	os.Setenv("RUNNING_IN_DOCKER", "true")

	if !isDocker() {
		t.Error("expected true when RUNNING_IN_DOCKER is set")
	}
}

func TestConfigure_Local(t *testing.T) {
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("RUNNING_IN_DOCKER")

	level, err := Configure("test-version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if level != "info" {
		t.Errorf("expected info, got %s", level)
	}
}

func TestConfigure_DebugEnv(t *testing.T) {
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	level, err := Configure("test-version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if level != "debug" {
		t.Errorf("expected debug, got %s", level)
	}
}
