package main

import (
	"context"
	"fmt"

	"github.com/hbagdi/go-kong/kong"

	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) client(ctx context.Context, s logical.Storage) (*kong.Client,
	error, error) {

	conf, userErr, intErr := b.readConfigAccess(ctx, s)
	if intErr != nil {
		return nil, nil, intErr
	}
	if userErr != nil {
		return nil, userErr, nil
	}
	if conf == nil {
		return nil, nil, fmt.Errorf("no error received but no config found")
	}

	client, err := kong.NewClient(kong.String(conf.BaseURL), nil)
	if err != nil {
		return nil, nil, err
	}

	return client, nil, nil
}
