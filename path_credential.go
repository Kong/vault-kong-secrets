package main

import (
	"context"
	"fmt"

	"github.com/hbagdi/go-kong/kong"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredential(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("consumer"),

		Fields: map[string]*framework.FieldSchema{
			"consumer": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The consumer for which to generate key-auth credentials",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredentialRead,
		},
	}
}

func (b *backend) pathCredentialRead(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	consumer := d.Get("consumer").(string)

	entry, err := req.Storage.Get(ctx, "consumer/"+consumer)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving consumer: {{ err }}", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("consumer %q not found", consumer)), nil
	}

	var result consumerConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	// kong client
	client, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	// create the key-auth credential
	keyAuth, err := client.KeyAuths.Create(ctx, kong.String(consumer), &kong.KeyAuth{})
	if err != nil {
		return nil, err
	}

	// user the helper to create the Secret
	s := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"token": keyAuth.Key,
		"id":    keyAuth.ID,
	}, map[string]interface{}{
		"id":       keyAuth.ID,
		"consumer": consumer,
	})
	s.Secret.TTL = result.TTL
	s.Secret.MaxTTL = result.MaxTTL

	return s, nil
}
