# vault-kong-secrets

Vault secrets engine for Kong API Gateway. Manage Kong authentication credentials via HashiCorp Vault.

## Rationale

This plugin allows for deeper integration of the Kong API Gateway and HashiCorp Vault. Users can leverage Vault's built-in authorization mechanisms to automatically generate API authentication keys, leaving the heavy lifting of cryptographically secured stored, TTL management, and authentication to Vault. Moreover, Vault's design allows it to managed the lifecycle of API authentication tokens without storing the token itself; Vault generates the token in Kong and passes it back to the client without storing the secret internally.

This plugin was inspired by Nicolas Corrarello's talk at HashiConf 2019 entitled ["Building Vault Plug-ins: From No-Go to Dynamic Secrets"](https://www.hashicorp.com/resources/building-vault-plug-ins-from-no-go-to-dynamic-secrets).

## Setup

The setup guide assumes some familiarity with Vault and Vault's plugin ecosystem. You must have a Vault server already running, unsealed, and authenticated.

1. Download and decompress the latest plugin binary from the Releases tab on GitHub. Alternatively you can compile the plugin from source (Go 1.25+).

2. Move the compiled plugin into Vault's configured `plugin_directory`:

    ```sh
    $ mv vault-kong-secrets /etc/vault/plugins/vault-kong-secrets
    ```

3. Calculate the SHA256 of the plugin and register it in Vault's plugin catalog. If you are downloading the pre-compiled binary, it is highly recommended that you use the published checksums to verify integrity. The plugin reports `RunningVersion` `v0.1.0`.

    ```sh
    $ export SHA256=$(shasum -a 256 "/etc/vault/plugins/vault-kong-secrets" | cut -d' ' -f1)

    $ vault plugin register \
        -sha256="${SHA256}" \
        -command="vault-kong-secrets" \
        -version=v0.1.0 \
        secret kong-secrets
    ```

4. Mount the secrets engine:

    ```sh
    $ vault secrets enable -path="kong" kong-secrets
    ```

## Usage & API

### Configure Kong Access

Manage connection parameters to the Kong cluster.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `PUT` | `/config/access` | `204 No Content` |

#### Parameters

* `baseurl` `(string)`- Kong Admin API URL.

* `admin_token` `(string)`- Kong-Admin-Token for Kong Enterprise RBAC. Omit when
RBAC is disabled.

* `workspace` `(string)`- Kong workspace name. Omit to use the default workspace.

#### CLI

```
$ vault write kong/config/access baseurl=http://127.0.0.1:8001
```

```
$ vault write kong/config/access \
    baseurl=http://127.0.0.1:8001 \
    admin_token=<kong-admin-token> \
    workspace=<workspace>
```

### Configure Consumer

Manage username and credential lease TTLs for an existing Kong Consumer.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `PUT` | `/consumers/:username` | `204 No Content` |

#### Parameters

* `username` `(string)`- The Kong Consumer username on which to generate credentials.

* `ttl` `(duration: "")`- Specifies the TTL for the generated credential. This is provided as a string duration with a time suffix like "30s" or "1h" or as seconds. If not provided, the default Vault TTL is used.

* `max_ttl` `(duration: "")`- Specifies the max TTL for the generated credential. This is provided as a string duration with a time suffix like "30s" or "1h" or as seconds. If not provided, the default Vault max TTL is used.

#### CLI

```
$ vault write kong/consumers/foo max_ttl=600s ttl=5m
```

### Generate key-auth Credential

Generate a Kong key-auth credential for the given consumer.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `GET` | `/creds/:username` | `200 (application/json)` |

#### Parameters

* `username` `(string)`- The Kong Consumer username on which to generate credentials.

#### CLI

Request:
```
$ vault read kong/creds/foo
```

Result:
```
Key                Value
---                -----
lease_id           kong/creds/foo/<lease-id>
lease_duration     5m
lease_renewable    false
id                 <credential-id>
key                <key>
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
