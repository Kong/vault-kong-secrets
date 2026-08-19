package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/kong/go-kong/kong"
)

func pathCredential(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("consumer"),
		Fields: map[string]*framework.FieldSchema{
			"consumer": {
				Type:        framework.TypeString,
				Description: "The consumer for which to generate key-auth credentials",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{Callback: b.pathCredentialRead},
		},
	}
}

func (b *backend) pathCredentialRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	consumer := d.Get("consumer").(string)

	entry, err := req.Storage.Get(ctx, "consumer/"+consumer)
	if err != nil {
		return nil, fmt.Errorf("error retrieving consumer: %w", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("consumer %q not found", consumer)), nil
	}

	var result consumerConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	client, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	keyAuth, err := client.KeyAuths.Create(ctx, kong.String(consumer), &kong.KeyAuth{})
	if err != nil {
		return nil, err
	}
	if keyAuth == nil || keyAuth.Key == nil || keyAuth.ID == nil {
		return nil, fmt.Errorf("kong returned an empty key-auth credential")
	}

	s := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"key": *keyAuth.Key,
		"id":  *keyAuth.ID,
	}, map[string]interface{}{
		"id":       *keyAuth.ID,
		"consumer": consumer,
	})
	s.Secret.TTL = result.TTL
	s.Secret.MaxTTL = result.MaxTTL

	return s, nil
}
