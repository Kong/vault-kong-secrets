package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Kong/vault-kong-secrets/hostsdk"
)

func pathListConsumers(b *backend) *hostsdk.Path {
	return &hostsdk.Path{
		Pattern: "consumers/?$",
		Operations: map[hostsdk.Operation]hostsdk.OperationHandler{
			hostsdk.ListOperation: &hostsdk.PathOperation{Callback: b.pathConsumersList},
		},
	}
}

func (b *backend) pathConsumersList(ctx context.Context, req *hostsdk.Request, d *hostsdk.FieldData) (*hostsdk.Response, error) {
	entries, err := req.Storage.List(ctx, "consumer/")
	if err != nil {
		return nil, err
	}
	return hostsdk.ListResponse(entries), nil
}

type consumerConfig struct {
	Username string        `json:"username"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
}

func pathConsumers(b *backend) *hostsdk.Path {
	return &hostsdk.Path{
		Pattern: "consumers/" + hostsdk.GenericNameRegex("username"),
		Fields: map[string]*hostsdk.FieldSchema{
			"username": {
				Type:        hostsdk.TypeString,
				Description: "username of the Kong Consumer",
			},
			"ttl": {
				Type:        hostsdk.TypeDurationSecond,
				Description: "TTL for the key-auth credential created for the consumer",
			},
			"max_ttl": {
				Type:        hostsdk.TypeDurationSecond,
				Description: "Max TTL for the key-auth credential created for the consumer",
			},
		},
		Operations: map[hostsdk.Operation]hostsdk.OperationHandler{
			hostsdk.ReadOperation:   &hostsdk.PathOperation{Callback: b.pathConsumerRead},
			hostsdk.UpdateOperation: &hostsdk.PathOperation{Callback: b.pathConsumerUpdate},
			hostsdk.DeleteOperation: &hostsdk.PathOperation{Callback: b.pathConsumerDelete},
		},
	}
}

func (b *backend) pathConsumerRead(ctx context.Context, req *hostsdk.Request, d *hostsdk.FieldData) (*hostsdk.Response, error) {
	username := d.Get("username").(string)

	entry, err := req.Storage.Get(ctx, "consumer/"+username)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return hostsdk.ErrorResponse(fmt.Sprintf("consumer %q not found", username)), nil
	}

	var result consumerConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &hostsdk.Response{
		Data: map[string]interface{}{
			"username": result.Username,
			"ttl":      int(result.TTL.Seconds()),
			"max_ttl":  int(result.MaxTTL.Seconds()),
		},
	}, nil
}

func (b *backend) pathConsumerUpdate(ctx context.Context, req *hostsdk.Request, d *hostsdk.FieldData) (*hostsdk.Response, error) {
	username := d.Get("username").(string)

	var ttl time.Duration
	if ttlRaw, ok := d.GetOk("ttl"); ok {
		ttl = time.Second * time.Duration(ttlRaw.(int))
	}

	var maxTTL time.Duration
	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		maxTTL = time.Second * time.Duration(maxTTLRaw.(int))
	}

	entry, err := hostsdk.StorageEntryJSON("consumer/"+username, consumerConfig{
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

func (b *backend) pathConsumerDelete(ctx context.Context, req *hostsdk.Request, d *hostsdk.FieldData) (*hostsdk.Response, error) {
	username := d.Get("username").(string)
	if err := req.Storage.Delete(ctx, "consumer/"+username); err != nil {
		return nil, err
	}
	return nil, nil
}
