package main

import (
	"context"
	"fmt"

	"github.com/Kong/vault-kong-secrets/hostsdk"
)

type accessConfig struct {
	BaseURL    string `json:"baseurl"`
	AdminToken string `json:"admin_token"`
	Workspace  string `json:"workspace"`
}

func pathConfigAccess(b *backend) *hostsdk.Path {
	return &hostsdk.Path{
		Pattern: "config/access",
		Fields: map[string]*hostsdk.FieldSchema{
			"baseurl": {
				Type:        hostsdk.TypeString,
				Description: "Kong Admin API base URL (for example http://kong-cp:8001)",
				Required:    true,
			},
			"admin_token": {
				Type:        hostsdk.TypeString,
				Description: "Kong-Admin-Token for RBAC (Enterprise). Leave empty if RBAC is off.",
				DisplayAttrs: &hostsdk.DisplayAttributes{
					Sensitive: true,
				},
			},
			"workspace": {
				Type:        hostsdk.TypeString,
				Description: "Kong workspace name. Empty uses the default workspace.",
			},
		},
		Operations: map[hostsdk.Operation]hostsdk.OperationHandler{
			hostsdk.ReadOperation:   &hostsdk.PathOperation{Callback: b.pathConfigAccessRead},
			hostsdk.UpdateOperation: &hostsdk.PathOperation{Callback: b.pathConfigAccessWrite},
		},
	}
}

func (b *backend) readConfigAccess(ctx context.Context, storage hostsdk.Storage) (*accessConfig, error, error) {
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

func (b *backend) pathConfigAccessRead(ctx context.Context, req *hostsdk.Request, data *hostsdk.FieldData) (*hostsdk.Response, error) {
	conf, userErr, intErr := b.readConfigAccess(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return hostsdk.ErrorResponse(userErr.Error()), nil
	}

	return &hostsdk.Response{
		Data: map[string]interface{}{
			"baseurl":   conf.BaseURL,
			"workspace": conf.Workspace,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(ctx context.Context, req *hostsdk.Request, data *hostsdk.FieldData) (*hostsdk.Response, error) {
	entry, err := hostsdk.StorageEntryJSON("config/access", accessConfig{
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
