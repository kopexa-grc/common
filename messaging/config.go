// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import "strings"

type NatsAuthMethod string

const (
	NatsAuthMethodUser  NatsAuthMethod = "user"
	NatsAuthMethodToken NatsAuthMethod = "token"
	NatsAuthMethodTLS   NatsAuthMethod = "tls"
)

type NatsAuth struct {
	Method   NatsAuthMethod `mapstructure:"method"`
	User     string         `mapstructure:"user"`
	Password string         `mapstructure:"password"`
	Token    string         `mapstructure:"token"`
}

type NATSConfig struct {
	Servers []string  `mapstructure:"servers"`
	Auth    NatsAuth  `mapstructure:"auth"`
	TLS     TLSConfig `mapstructure:"tls"`
}

// ServerString will build the proper string for nats connect
func (c *NATSConfig) ServerString() string {
	return strings.Join(c.Servers, ",")
}
