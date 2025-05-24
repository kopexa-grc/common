// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import "github.com/openfga/go-sdk/credentials"

// Config represents the configuration for the OpenFGA client.
// It contains all necessary settings to connect to and interact with the OpenFGA service.
//
// This struct is used to configure:
// - Connection settings (host URL, store details)
// - Authentication credentials
// - Authorization model settings
// - Error handling behavior
type Config struct {
	// Enabled determines whether authorization checks with OpenFGA are active.
	// When disabled, all authorization checks will be bypassed.
	Enabled bool `json:"enabled" koanf:"enabled" jsonschema:"description=enables authorization checks with openFGA" default:"true"`

	// StoreName is the name of the OpenFGA store to use.
	// This is used to identify the authorization store in the OpenFGA service.
	StoreName string `json:"storeName" koanf:"storeName" jsonschema:"description=name of openFGA store" default:"kopexa"`

	// HostURL is the complete URL (including scheme) of the OpenFGA API.
	// This is required and must be a valid URL.
	HostURL string `json:"hostUrl" koanf:"hostUrl" jsonschema:"description=host url with scheme of the openFGA API,required" default:"https://authz.theopenlane.io"`

	// StoreID is the unique identifier of the OpenFGA store.
	// This is required when working with multiple stores.
	StoreID string `json:"storeId" koanf:"storeId" jsonschema:"description=id of openFGA store"`

	// CreateNewModel forces the creation of a new authorization model,
	// even if one already exists in the store.
	CreateNewModel bool `json:"createNewModel" koanf:"createNewModel" jsonschema:"description=force create a new model, even if one already exists" default:"false"`

	// IgnoreDuplicateKeyError determines whether to ignore errors when:
	// - A key already exists during creation
	// - A delete request is made for a non-existent key
	// This is useful for idempotent operations.
	IgnoreDuplicateKeyError bool `json:"ignoreDuplicateKeyError" koanf:"ignoreDuplicateKeyError" jsonschema:"description=ignore duplicate key error" default:"true"`

	// Credentials contains the authentication information for the OpenFGA service.
	// This is required for all API calls to the service.
	Credentials Credentials `json:"credentials" koanf:"credentials" jsonschema:"description=credentials for the openFGA client"`

	// ModelID specifies an existing authorization model ID to use.
	// If not set, the latest model will be used.
	ModelID string `json:"modelId" koanf:"modelId" jsonschema:"description=id of openFGA model"`

	// ModelFile is the path to the FGA model file.
	// This file contains the authorization model definition.
	ModelFile string `json:"modelFile" koanf:"modelFile" jsonschema:"description=path to the fga model file" default:"fga/model/model.fga"`

	// ModelData contains the raw FGA model definition.
	// This can be used instead of ModelFile to provide the model directly.
	ModelData []byte `json:"modelData" koanf:"modelData" jsonschema:"description=data of the fga model"`
}

// Credentials represents the authentication credentials for the OpenFGA service.
// It supports both API token and client credentials authentication methods.
//
// This struct is used to configure:
// - API token authentication
// - OAuth2 client credentials authentication
// - Authentication-related settings (audience, issuer, scopes)
type Credentials struct {
	// APIToken is the API token used to authenticate with the OpenFGA service.
	// This token is required for all API calls to the service.
	// It is used when Method is set to CredentialsMethodApiToken.
	APIToken string `json:"api_token" koanf:"api_token" jsonschema:"description=The API token for the OpenFGA service"`

	// ClientID is the client ID used for OAuth2 client credentials authentication.
	// This is required when using client credentials authentication.
	// It is used when Method is set to CredentialsMethodClientCredentials.
	ClientID string `json:"clientId" koanf:"clientId" jsonschema:"description=client id for the openFGA client, required if using client credentials authentication"`

	// ClientSecret is the client secret used for OAuth2 client credentials authentication.
	// This is required when using client credentials authentication.
	// It is used when Method is set to CredentialsMethodClientCredentials.
	ClientSecret string `json:"clientSecret" koanf:"clientSecret" jsonschema:"description=client secret for the openFGA client, required if using client credentials authentication"`

	// Audience is the OAuth2 audience value for client credentials authentication.
	// This is used to specify the intended recipient of the token.
	// It is used when Method is set to CredentialsMethodClientCredentials.
	Audience string `json:"audience" koanf:"audience" jsonschema:"description=audience for the openFGA client"`

	// Issuer is the OAuth2 token issuer for client credentials authentication.
	// This is used to specify the entity that issued the token.
	// It is used when Method is set to CredentialsMethodClientCredentials.
	Issuer string `json:"issuer" koanf:"issuer" jsonschema:"description=issuer for the openFGA client"`

	// Scopes is the OAuth2 scopes for client credentials authentication.
	// This is used to specify the permissions requested by the client.
	// It is used when Method is set to CredentialsMethodClientCredentials.
	Scopes string `json:"scopes" koanf:"scopes" jsonschema:"description=scopes for the openFGA client"`
}

// Option is a function that configures a Client.
// Options are used to customize the behavior of the FGA client.
// They follow the functional options pattern for flexible configuration.
type Option func(*Client)

// WithStoreID sets the store ID for the FGA client.
// This is required when working with multiple stores.
// The store ID is used to identify which authorization store to use.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithStoreID("store-123"),
//	)
func WithStoreID(storeID string) Option {
	return func(c *Client) {
		c.config.StoreId = storeID
	}
}

// WithIgnoreDuplicateKeyError configures whether duplicate key errors should be ignored.
// When set to true, attempts to write duplicate tuples will be silently ignored.
// This is useful in scenarios where idempotency is desired.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithIgnoreDuplicateKeyError(true),
//	)
func WithIgnoreDuplicateKeyError(ignore bool) Option {
	return func(c *Client) {
		c.IgnoreDuplicateKeyError = ignore
	}
}

// WithToken configures the FGA client with an API token for authentication.
// The token is used to authenticate all requests to the OpenFGA service.
// This option is required for production use of the client.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithToken("your-api-token"),
//	)
func WithToken(token string) Option {
	return func(c *Client) {
		c.config.Credentials = &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: token,
			},
		}
	}
}

// WithAPITokenCredentials sets the credentials for the client with an API token.
// This is an alternative to WithToken that provides more explicit naming.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithAPITokenCredentials("your-api-token"),
//	)
func WithAPITokenCredentials(token string) Option {
	return func(c *Client) {
		c.config.Credentials = &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: token,
			},
		}
	}
}

// WithClientCredentials sets the client credentials for OAuth2 authentication.
// This configures the client to use OAuth2 client credentials flow.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithClientCredentials(
//	        "client-id",
//	        "client-secret",
//	        "https://api.openfga.example",
//	        "https://auth.openfga.example",
//	        "read write",
//	    ),
//	)
func WithClientCredentials(clientID, clientSecret, aud, issuer, scopes string) Option {
	return func(c *Client) {
		c.config.Credentials = &credentials.Credentials{
			Method: credentials.CredentialsMethodClientCredentials,
			Config: &credentials.Config{
				ClientCredentialsClientId:       clientID,
				ClientCredentialsClientSecret:   clientSecret,
				ClientCredentialsApiAudience:    aud,
				ClientCredentialsApiTokenIssuer: issuer,
				ClientCredentialsScopes:         scopes,
			},
		}
	}
}

// WithAuthorizationModelID sets the authorization model ID to use.
// This allows specifying a particular version of the authorization model.
// If not set, the latest model will be used.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithAuthorizationModelID("model-123"),
//	)
func WithAuthorizationModelID(authModelID string) Option {
	return func(c *Client) {
		c.config.AuthorizationModelId = authModelID
	}
}
