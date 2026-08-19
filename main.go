package main

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
)

// runningVersion is reported to Vault's plugin catalog (leading v required).
// See https://developer.hashicorp.com/vault/docs/plugins/plugin-development
const runningVersion = "v0.1.0"

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.ServeMultiplex(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		hclog.New(&hclog.LoggerOptions{}).Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}

type backend struct {
	*framework.Backend
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help:           "Dynamic Kong key-auth credentials",
		RunningVersion: runningVersion,
		Paths: []*framework.Path{
			pathConfigAccess(&b),
			pathListConsumers(&b),
			pathConsumers(&b),
			pathCredential(&b),
		},
		Secrets: []*framework.Secret{
			secretCredential(&b),
		},
		BackendType: logical.TypeLogical,
	}
	return &b
}
