//go:build openbao

package hostsdk

import (
	"github.com/openbao/openbao/api/v2"
	"github.com/openbao/openbao/sdk/v2/framework"
	"github.com/openbao/openbao/sdk/v2/logical"
	"github.com/openbao/openbao/sdk/v2/plugin"
)

type PluginAPIClientMeta = api.PluginAPIClientMeta

var VaultPluginTLSProvider = api.VaultPluginTLSProvider

type ServeOpts = plugin.ServeOpts

func ServeMultiplex(opts *ServeOpts) error {
	return plugin.ServeMultiplex(opts)
}

type (
	Backend           = framework.Backend
	Path              = framework.Path
	FieldSchema       = framework.FieldSchema
	DisplayAttributes = framework.DisplayAttributes
	PathOperation     = framework.PathOperation
	Secret            = framework.Secret
	FieldData         = framework.FieldData
	OperationHandler  = framework.OperationHandler
	LogicalBackend    = logical.Backend
	BackendConfig     = logical.BackendConfig
	Storage           = logical.Storage
	Operation         = logical.Operation
	Request           = logical.Request
	Response          = logical.Response
)

const (
	TypeString         = framework.TypeString
	TypeDurationSecond = framework.TypeDurationSecond
	TypeLogical        = logical.TypeLogical
	ReadOperation      = logical.ReadOperation
	UpdateOperation    = logical.UpdateOperation
	DeleteOperation    = logical.DeleteOperation
	ListOperation      = logical.ListOperation
)

var (
	GenericNameRegex = framework.GenericNameRegex
	ErrorResponse    = logical.ErrorResponse
	ListResponse     = logical.ListResponse
	StorageEntryJSON = logical.StorageEntryJSON
)
