# vault-kong-secrets

Vault secrets engine for Kong API Gateway. Manage Kong authentication credentials
via HashiCorp Vault.

## Rationale

This plugin allows for deeper integration of the Kong API Gateway and HashiCorp
Vault. Users can leverage Vault's built-in authorization mechanisms to automatically
generate API authentication keys, leaving the heavy lifting of cryptographically
secured stored, TTL management, and authentication to Vault. Moreover, Vault's
design allows it to managed the lifecycle of API authentication tokens without
storing the token itself; Vault generates the token in Kong and passes it back
to the client without storing the secret internally.

This plugin was inspired by Nicolas Corrarello's talk at HashiConf 2019 entitled
"Building Vault Plug-ins: From No-Go to Dynamic Secrets".

## Setup

The setup guide assumes some familiarity with Vault and Vault's plugin
ecosystem. You must have a Vault server already running, unsealed, and
authenticated.

1. Download and decompress the latest plugin binary from the Releases tab on
GitHub. Alternatively you can compile the plugin from souce.

1. Move the compiled plugin into Vault's configured `plugin_directory`:

    ```sh
    $ mv vault-kong-secrets /etc/vault/plugins/vault-kong-secrets
    ```

1. Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.
If you are downloading the pre-compiled binary, it is highly recommended that
you use the published checksums to verify integrity.

    ```sh
    $ export SHA256=$(shasum -a 256 "/etc/vault/plugins/vault-kong-secrets" | cut -d' ' -f1)

    $ vault write sys/plugins/catalog/secret/kong-secrets \
        sha_256="${SHA256}" \
        command="vault-kong-secrets"
    ```

1. Mount the secrets engine:

    ```sh
    $ vault secrets enable \
        -path="kong" \
        -plugin-name="kong-secrets" \
        plugin
    ```

## Usage & API

### Configure Kong Access

Manage connection parameters to the Kong cluster.

| Method   | Path                         | Produces                 |
| :------- | :--------------------------- | :----------------------- |
| `PUT`    | `/config/access`             | `204 No Content`         |

#### Parameters

* `baseurl` `(string)`- The URL scheme where the Kong cluster can be found.

#### CLI

```
$ vault write kong/config/access baseurl=http://127.0.0.1:8001
```

### Configure Consumer

Manage Kong Consumer username and credential lease TTLs.

| Method   | Path                         | Produces                 |
| :------- | :--------------------------- | :----------------------- |
| `PUT`    | `/consumers/:username`       | `204 No Content`         |

#### Parameters

* `username` `(string)`- The Kong Consumer username on which to generate credentials.

* `ttl` `(duration: "")`- Specifies the TTL for the generated credential. This
is provided as a string duration with a time suffix like "30s" or "1h" or as
seconds. If not provided, the default Vault TTL is used.

* `max_ttl` `(duration: "")`- Specifies the max TTL for the generated credential.
This is provided as a string duration with a time suffix like "30s" or "1h" or
as seconds. If not provided, the default Vault TTL is used.

#### CLI

```
$ vault write kong/consumers/foo max_ttl=600s ttl=20s
```

### Generate key-auth Credential

Generate a Kong key-auth credential for the given consumer.

| Method   | Path                         | Produces                 |
| :------- | :--------------------------- | :----------------------- |
| `GET`    | `/creds/:username`           | `200 (application/json)` |

#### Parameters

* `username` `(string)`- The Kong Consumer username on which to generate credentials.

#### CLI

```
$ vault read kong/creds/foo
```

### License

Copyright 2019 Kong Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
