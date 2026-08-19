package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListConsumers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "consumers/?$",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{Callback: b.pathConsumersList},
		},
	}
}

func (b *backend) pathConsumersList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "consumer/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(entries), nil
}

type consumerConfig struct {
	Username string        `json:"username"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
}

func pathConsumers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "consumers/" + framework.GenericNameRegex("username"),
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "username of the Kong Consumer",
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "TTL for the key-auth credential created for the consumer",
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Max TTL for the key-auth credential created for the consumer",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation:   &framework.PathOperation{Callback: b.pathConsumerRead},
			logical.UpdateOperation: &framework.PathOperation{Callback: b.pathConsumerUpdate},
			logical.DeleteOperation: &framework.PathOperation{Callback: b.pathConsumerDelete},
		},
	}
}

func (b *backend) pathConsumerRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)

	entry, err := req.Storage.Get(ctx, "consumer/"+username)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("consumer %q not found", username)), nil
	}

	var result consumerConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"username": result.Username,
			"ttl":      int(result.TTL.Seconds()),
			"max_ttl":  int(result.MaxTTL.Seconds()),
		},
	}, nil
}

func (b *backend) pathConsumerUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)

	var ttl time.Duration
	if ttlRaw, ok := d.GetOk("ttl"); ok {
		ttl = time.Second * time.Duration(ttlRaw.(int))
	}

	var maxTTL time.Duration
	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		maxTTL = time.Second * time.Duration(maxTTLRaw.(int))
	}

	entry, err := logical.StorageEntryJSON("consumer/"+username, consumerConfig{
		Username: username,
		TTL:      ttl,
		MaxTTL:   maxTTL,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConsumerDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if err := req.Storage.Delete(ctx, "consumer/"+username); err != nil {
		return nil, err
	}
	return nil, nil
}
