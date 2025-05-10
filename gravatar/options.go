// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gravatar

// Options allows customization of the generated Gravatar URL.
type Options struct {
	Size          int    // Size in pixels (1–2048)
	DefaultImage  string // Fallback image (404, mp, identicon, etc.)
	ForceDefault  bool   // Always use default image
	Rating        string // Content rating (g, pg, r, x)
	FileExtension string // Optional file extension (.jpg, .png, etc.)
}

type Option func(*Options)

func WithSize(size int) Option {
	return func(o *Options) {
		o.Size = size
	}
}

func WithDefaultImage(image string) Option {
	return func(o *Options) {
		o.DefaultImage = image
	}
}

func WithForceDefault(force bool) Option {
	return func(o *Options) {
		o.ForceDefault = force
	}
}

func WithRating(rating string) Option {
	return func(o *Options) {
		o.Rating = rating
	}
}

func WithFileExtension(ext string) Option {
	return func(o *Options) {
		o.FileExtension = ext
	}
}
