// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"crypto/tls"
	"fmt"

	"github.com/kopexa-grc/common/errors"
	"github.com/nats-io/nats.go"
)

var (
	ErrTLSConfigMissing  = errors.New(errors.BadRequest, "tls is enabled but no valid tls configuration was provided")
	ErrTLSCertMissing    = errors.New(errors.BadRequest, "tls auth method is configured but no certificate was loaded")
	ErrInvalidAuthMethod = errors.New(errors.BadRequest, "invalid auth method")
	ErrRootCAPoolNil     = errors.New(errors.BadRequest, "nats: the root ca pool from the given tls.Config is nil")
)

func NewNATSClient(cfg *NATSConfig, opts ...nats.Option) (*nats.Conn, error) {
	tlsConfig, err := cfg.TLS.TLSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure tls")
	}

	// If TLS is enabled, make the connection to NATS secure
	if cfg.TLS.Enabled {
		if tlsConfig == nil {
			return nil, ErrTLSConfigMissing
		}

		opts = append(opts, nats.Secure(tlsConfig))
	}

	switch cfg.Auth.Method {
	case NatsAuthMethodUser:
		opts = append(opts, nats.UserInfo(cfg.Auth.User, cfg.Auth.Password))
	case NatsAuthMethodToken:
		opts = append(opts, nats.Token(cfg.Auth.Token))
	case NatsAuthMethodTLS:
		// if using TLS auth, make sure the client certificate is loaded
		if tlsConfig == nil || len(tlsConfig.Certificates) == 0 {
			return nil, ErrTLSCertMissing
		}
	case "":
		// noop ~ we aren't using any auth method
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidAuthMethod, cfg.Auth.Method)
	}

	return nats.Connect(cfg.ServerString(), opts...)
}

// NatsRootCAs is a NATS helper option to provide the RootCAs pool from a tls.Config struct. If Secure is
// not already set this will set it as well.
func NatsRootCAs(tlsConf *tls.Config) nats.Option {
	return func(o *nats.Options) error {
		if tlsConf.RootCAs == nil {
			return ErrRootCAPoolNil
		}

		if o.TLSConfig == nil {
			o.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				// nolint:gosec
				InsecureSkipVerify: tlsConf.InsecureSkipVerify,
			}
		}

		o.TLSConfig.RootCAs = tlsConf.RootCAs
		o.Secure = true

		return nil
	}
}
