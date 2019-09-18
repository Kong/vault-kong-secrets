package main

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListConsumers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "consumers/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathConsumersList,
		},
	}
}

func (b *backend) pathConsumersList(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

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
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "username of the Kong Consumer",
			},

			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "TTL for the key-auth credential created for the consumer",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "Max TTL for the key-auth credential created for the consumer",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConsumerRead,
			logical.UpdateOperation: b.pathConsumerUpdate,
			logical.DeleteOperation: b.pathConsumerDelete,
		},
	}
}

func (b *backend) pathConsumerRead(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	username := d.Get("username").(string)

	entry, err := req.Storage.Get(ctx, "consumer/"+username)
	if err != nil {
		return nil, err
	}

	var result consumerConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	// generate the response
	return &logical.Response{
		Data: map[string]interface{}{
			"username": result.Username,
			"ttl":      int(result.TTL.Seconds()),
			"max_ttl":  int(result.MaxTTL.Seconds()),
		},
	}, nil
}

func (b *backend) pathConsumerUpdate(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	username := d.Get("username").(string)

	var ttl time.Duration
	ttlRaw, ok := d.GetOk("ttl")
	if ok {
		ttl = time.Second * time.Duration(ttlRaw.(int))
	}

	var maxTTL time.Duration
	maxTTLRaw, ok := d.GetOk("max_ttl")
	if ok {
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

func (b *backend) pathConsumerDelete(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	username := d.Get("username").(string)

	if err := req.Storage.Delete(ctx, "consumer/"+username); err != nil {
		return nil, err
	}

	return nil, nil
}
