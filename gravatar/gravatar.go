// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gravatar

import (
	"crypto/md5" //nolint:gosec // Acceptable for Gravatar hashes
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
)

// URL returns the Gravatar image URL for the given email and options.
func URL(email string, opts ...Option) string {
	o := &Options{
		Size:         defaultSize,
		DefaultImage: defaultImage,
		Rating:       defaultRating,
	}

	for _, opt := range opts {
		opt(o)
	}

	hash := Hash(email)
	if ext := strings.TrimPrefix(o.FileExtension, "."); ext != "" {
		hash += "." + ext
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "www.gravatar.com",
		Path:     "avatar/" + hash,
		RawQuery: buildQuery(o).Encode(),
	}

	return u.String()
}

// Hash computes the lowercase-trimmed MD5 hash of the email.
func Hash(email string) string {
	normalized := strings.ToLower(strings.TrimSpace(email))
	sum := md5.Sum([]byte(normalized)) //nolint:gosec

	return hex.EncodeToString(sum[:])
}

// buildQuery constructs the query parameters based on Options.
func buildQuery(opts *Options) url.Values {
	q := url.Values{}

	if opts.Size > 0 {
		q.Set("s", strconv.Itoa(opts.Size))
	}

	if opts.DefaultImage != "" {
		q.Set("d", opts.DefaultImage)
	}

	if opts.ForceDefault {
		q.Set("f", "y")
	}

	if opts.Rating != "" {
		q.Set("r", opts.Rating)
	}

	return q
}
