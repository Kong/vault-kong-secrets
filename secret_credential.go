package main

import (
	"context"

	"github.com/hbagdi/go-kong/kong"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const SecretTokenType = "token"

func secretCredential(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The credential key",
			},
			"id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The credential id",
			},
		},

		Renew:  b.secretCredentialRenew,
		Revoke: b.secretCredentialRevoke,
	}
}

func (b *backend) secretCredentialRenew(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func (b *backend) secretCredentialRevoke(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	client, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return nil, userErr
	}

	idRaw, ok := req.Secret.InternalData["id"]
	if !ok {
		return nil, nil // pre 0.5.3 problem
	}
	id := idRaw.(string)

	consumerRaw := req.Secret.InternalData["consumer"]
	consumer := consumerRaw.(string)

	if err := client.KeyAuths.Delete(ctx, kong.String(consumer), kong.String(id)); err != nil {
		return nil, err
	}

	return nil, nil
}
