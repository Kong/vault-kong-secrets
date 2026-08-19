package main

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/kong/go-kong/kong"
)

const SecretTokenType = "token"

func secretCredential(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"key": {
				Type:        framework.TypeString,
				Description: "The credential key",
			},
			"id": {
				Type:        framework.TypeString,
				Description: "The credential id",
			},
		},
		Revoke: b.secretCredentialRevoke,
	}
}

func (b *backend) secretCredentialRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return nil, userErr
	}

	idRaw, ok := req.Secret.InternalData["id"]
	if !ok {
		return nil, nil
	}
	id, _ := idRaw.(string)
	consumerRaw, _ := req.Secret.InternalData["consumer"]
	consumer, _ := consumerRaw.(string)
	if id == "" || consumer == "" {
		return nil, nil
	}

	if err := client.KeyAuths.Delete(ctx, kong.String(consumer), kong.String(id)); err != nil {
		return nil, err
	}
	return nil, nil
}
