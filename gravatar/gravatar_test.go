// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gravatar_test

import (
	"testing"

	"github.com/kopexa-grc/common/gravatar"
	"github.com/stretchr/testify/require"
)

func TestGravatar_DefaultOptions(t *testing.T) {
	email := "julian@kopexa.com"
	url := gravatar.URL(email,
		gravatar.WithSize(80),
		gravatar.WithDefaultImage("robohash"),
		gravatar.WithRating("pg"),
		gravatar.WithForceDefault(false),
	)

	require.Equal(t,
		"https://www.gravatar.com/avatar/c6a9958d84d231fc31124cb3d44ea601?d=robohash&r=pg&s=80",
		url,
	)
}

func TestGravatar_WithExtension(t *testing.T) {
	email := "julian@kopexa.com"
	url := gravatar.URL(email, gravatar.WithFileExtension(".png"))

	require.Equal(t,
		"https://www.gravatar.com/avatar/c6a9958d84d231fc31124cb3d44ea601.png?d=robohash&r=pg&s=80",
		url,
	)
}

func TestHash(t *testing.T) {
	input := "julian@kopexa.com"
	expected := "c6a9958d84d231fc31124cb3d44ea601"
	require.Equal(t, expected, gravatar.Hash(input))
}
