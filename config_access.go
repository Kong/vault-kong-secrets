package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type accessConfig struct {
	BaseURL string `json:"baseurl"`
}

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",

		Fields: map[string]*framework.FieldSchema{
			"baseurl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Kong server base URL",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigAccessRead,
			logical.UpdateOperation: b.pathConfigAccessWrite,
		},
	}
}

func (b *backend) readConfigAccess(ctx context.Context, storage logical.Storage) (*accessConfig, error, error) {
	entry, err := storage.Get(ctx, "config/access")
	if err != nil {
		return nil, nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("access credentials for the backend itself haven't been configured; please configure them at the '/config/access' endpoint"), nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, nil, errwrap.Wrapf("error reading consul access configuration: {{err}}", err)
	}

	return conf, nil, nil
}

func (b *backend) pathConfigAccessRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, userErr, intErr := b.readConfigAccess(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}
	if conf == nil {
		return nil, fmt.Errorf("no user error reported but consul access configuration not found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"baseurl": conf.BaseURL,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		BaseURL: data.Get("baseurl").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}
