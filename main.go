package main

import (
	"context"
	"os"

	"github.com/Kong/vault-kong-secrets/hostsdk"
	"github.com/hashicorp/go-hclog"
)

// runningVersion is reported to Vault's plugin catalog (leading v required).
// Release builds override this from the git tag via -ldflags.
// See https://developer.hashicorp.com/vault/docs/plugins/plugin-development
var runningVersion = "v0.2.0"

func main() {
	apiClientMeta := &hostsdk.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := hostsdk.VaultPluginTLSProvider(tlsConfig)

	err := hostsdk.ServeMultiplex(&hostsdk.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		hclog.New(&hclog.LoggerOptions{}).Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}

type backend struct {
	*hostsdk.Backend
}

func Factory(ctx context.Context, conf *hostsdk.BackendConfig) (hostsdk.LogicalBackend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &hostsdk.Backend{
		Help:           "Dynamic Kong key-auth credentials",
		RunningVersion: runningVersion,
		Paths: []*hostsdk.Path{
			pathConfigAccess(&b),
			pathListConsumers(&b),
			pathConsumers(&b),
			pathCredential(&b),
		},
		Secrets: []*hostsdk.Secret{
			secretCredential(&b),
		},
		BackendType: hostsdk.TypeLogical,
	}
	return &b
}
