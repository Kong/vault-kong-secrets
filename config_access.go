package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type accessConfig struct {
	BaseURL    string `json:"baseurl"`
	AdminToken string `json:"admin_token"`
	Workspace  string `json:"workspace"`
}

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"baseurl": {
				Type:        framework.TypeString,
				Description: "Kong Admin API base URL (for example http://kong-cp:8001)",
				Required:    true,
			},
			"admin_token": {
				Type:        framework.TypeString,
				Description: "Kong-Admin-Token for RBAC (Enterprise). Leave empty if RBAC is off.",
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
			"workspace": {
				Type:        framework.TypeString,
				Description: "Kong workspace name. Empty uses the default workspace.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation:   &framework.PathOperation{Callback: b.pathConfigAccessRead},
			logical.UpdateOperation: &framework.PathOperation{Callback: b.pathConfigAccessWrite},
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
		return nil, nil, fmt.Errorf("error reading kong access configuration: %w", err)
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

	return &logical.Response{
		Data: map[string]interface{}{
			"baseurl":   conf.BaseURL,
			"workspace": conf.Workspace,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		BaseURL:    data.Get("baseurl").(string),
		AdminToken: data.Get("admin_token").(string),
		Workspace:  data.Get("workspace").(string),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}
