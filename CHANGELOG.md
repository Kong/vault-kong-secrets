# vault-kong-secrets changelog

## 0.2.0

**Release date:** 2026/08/20

### Feature

#### Plugin

* Added a dual-host build: the default binary targets HashiCorp Vault, and `go build -tags openbao` targets OpenBao.
* Added GitHub Release artifacts for both hosts (`vault-kong-secrets_*` and `openbao-kong-secrets_*`).

#### Docs

* Documented that `ttl` and `max_ttl` are Vault or OpenBao lease durations. Kong key-auth credentials are created without a Kong TTL and are deleted when the lease is revoked.
* Documented Kong Gateway OSS versus Enterprise Admin access (`admin_token`, `workspace`).
* Documented tested versions:
  * Kong Gateway OSS 3.9
  * Kong Gateway Enterprise 3.15
  * HashiCorp Vault 1.19 / 1.20 / 1.21 / 2.0
  * OpenBao 2.4 / 2.5 / 2.6

## 0.1.0

**Release date:** 2026/08/19

### Feature

#### Plugin

* Issued Kong key-auth credentials through Vault without storing the key in Vault. The lease revoke callback deletes the credential from Kong.
* Added Kong Gateway Enterprise Admin access via `Kong-Admin-Token` and workspace selection. OSS uses Admin API root and omits both fields.
* Registered the plugin with `vault plugin register`, `ServeMultiplex`, and `RunningVersion` `v0.1.0`.

#### Docs

* Added English and Korean setup docs, including plugin register and checksum verification.

### Dependency

* Switched the Kong Admin client to `github.com/kong/go-kong`.
* Bumped HashiCorp Vault API/SDK and set the Go toolchain to 1.25.7.

### CI

* Added a tag-triggered GitHub Release workflow that publishes cross-platform binaries and `SHA256SUMS`.
