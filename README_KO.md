# vault-kong-secrets

Kong API Gateway용 Vault secrets engine. HashiCorp Vault로 Kong 인증 자격 증명을 관리합니다.

## 목적

Kong API Gateway와 HashiCorp Vault를 더 깊게 연동합니다. Vault의 인가 체계로 Kong의 Consumer에 대한 API 인증 키(Key Authentication)를 자동 발급하고, TTL 관리, 인증은 Vault에 맡깁니다. Vault는 키를 내부에 보관하지 않습니다. Vault가 Kong에서 key-auth 키를 생성해 클라이언트에 전달하고, 리스 수명으로 자격 증명의 생명주기를 관리합니다.

이 플러그인은 HashiConf 2019에서 Nicolas Corrarello가 발표한 ["Building Vault Plug-ins: From No-Go to Dynamic Secrets"](https://www.hashicorp.com/resources/building-vault-plug-ins-from-no-go-to-dynamic-secrets)에서 영감을 받았습니다.

## 적용 방식

Vault와 Vault 플러그인 생태계를 어느 정도 알고 있다고 가정합니다. Vault 서버가 이미 기동·unseal되어 있고, 인증된 상태여야 합니다.

1. GitHub Releases에서 최신 플러그인 바이너리를 받아 압축을 풉니다. 소스에서 빌드할 수도 있습니다 (Go 1.25.7+).

2. 컴파일된 플러그인을 Vault의 `plugin_directory`로 옮깁니다.

    ```sh
    $ mv vault-kong-secrets /etc/vault/plugins/vault-kong-secrets
    ```

3. 플러그인의 SHA256을 계산하고 Vault 플러그인 카탈로그에 등록합니다. 미리 컴파일된 바이너리를 쓰는 경우 게시된 체크섬으로 무결성을 확인하는 것을 권장합니다. 플러그인이 보고하는 `RunningVersion`을 확인합니다. (예제에서는 `v0.1.0`입니다.)

    ```sh
    $ export SHA256=$(shasum -a 256 "/etc/vault/plugins/vault-kong-secrets" | cut -d' ' -f1)

    $ vault plugin register \
        -sha256="${SHA256}" \
        -command="vault-kong-secrets" \
        -version=v0.1.0 \
        secret kong-secrets
    ```

4. secrets engine을 마운트합니다.

    ```sh
    $ vault secrets enable -path="kong" kong-secrets
    ```

## Usage & API

### Configure Kong Access

Kong 클러스터 접속 파라미터를 설정합니다.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `PUT` | `/config/access` | `204 No Content` |

#### 파라미터

* `baseurl` `(string)`- Kong Admin API URL.

* `admin_token` `(string)`- Kong Enterprise RBAC용 Kong-Admin-Token. RBAC가
꺼져 있으면 생략합니다.

* `workspace` `(string)`- Kong workspace 이름. 생략하면 default workspace를
사용합니다.

#### CLI 예시

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

Kong에 이미 있는 Consumer의 사용자 이름과 credential 리스 TTL을 설정합니다.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `PUT` | `/consumers/:username` | `204 No Content` |

#### 파라미터

* `username` `(string)`- 자격 증명을 발급할 Kong Consumer 사용자 이름.

* `ttl` `(duration: "")`- 생성된 자격 증명의 TTL. `"30s"`, `"1h"` 같은 접미사
또는 초 단위 문자열입니다. 생략하면 Vault 기본 TTL을 씁니다.

* `max_ttl` `(duration: "")`- 생성된 자격 증명의 최대 TTL. `"30s"`, `"1h"` 같은
접미사 또는 초 단위 문자열입니다. 생략하면 Vault 기본 최대 TTL을 씁니다.

#### CLI

```
$ vault write kong/consumers/foo max_ttl=10m ttl=5m
```

### Generate key-auth Credential

지정한 consumer에 대해 Kong key-auth 자격 증명을 발급합니다.

| Method | Path | Produces |
| :------- | :--------------------------- | :----------------------- |
| `GET` | `/creds/:username` | `200 (application/json)` |

#### 파라미터

* `username` `(string)`- 자격 증명을 발급할 Kong Consumer 사용자 이름.

#### CLI

요청:
```
$ vault read kong/creds/foo
```

결과:
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
