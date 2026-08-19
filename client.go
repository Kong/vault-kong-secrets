package main

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/kong/go-kong/kong"
)

type adminTokenTransport struct {
	token string
}

func (t *adminTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Kong-Admin-Token", t.token)
	return http.DefaultTransport.RoundTrip(req)
}

func (b *backend) client(ctx context.Context, s logical.Storage) (*kong.Client, error, error) {
	conf, userErr, intErr := b.readConfigAccess(ctx, s)
	if intErr != nil {
		return nil, nil, intErr
	}
	if userErr != nil {
		return nil, userErr, nil
	}

	var httpClient *http.Client
	if conf.AdminToken != "" {
		httpClient = &http.Client{
			Timeout:   kong.DefaultTimeout,
			Transport: &adminTokenTransport{token: conf.AdminToken},
		}
	}

	client, err := kong.NewClient(kong.String(conf.BaseURL), httpClient)
	if err != nil {
		return nil, nil, err
	}
	if conf.Workspace != "" {
		client.SetWorkspace(conf.Workspace)
	}

	return client, nil, nil
}
