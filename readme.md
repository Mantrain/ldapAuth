<p align="center">
<img src="imgs/gandalpher.png" alt="Gandalpher" title="Gandalpher" />
<br>
<em><b>"You shall authenticate to the LDAP to pass"</b> - Gandalpher, the gopher</em>
</p>

---

<h1 align="center">
<img alt="GitHub" src="https://img.shields.io/github/license/wiltonsr/ldapAuth?color=blue">
<img alt="GitHub release (latest by date including pre-releases)" src="https://img.shields.io/github/v/release/wiltonsr/ldapAuth?include_prereleases">
<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/wiltonsr/ldapAuth">
<img alt="GitHub issues" src="https://img.shields.io/github/issues/wiltonsr/ldapAuth">
<img alt="GitHub last commit (branch)" src="https://img.shields.io/github/last-commit/wiltonsr/ldapAuth/main">
</h1>

# Traefik ldapAuth Middleware

This project is an in progress effort to create an open source middleware that enables authentication via LDAP in a similar way to [Traefik Enterprise](https://doc.traefik.io/traefik-enterprise/middlewares/ldap/).

## Requirements

- Yaegi >= [v0.14.2](https://github.com/traefik/yaegi/releases/tag/v0.14.2)
- Traefik >= [v2.5.5](https://github.com/traefik/traefik/releases/tag/v2.5.5)
- go-ldap v3 >= [v3.1.4](https://github.com/go-ldap/ldap/releases/tag/v3.1.4)
- gorilla/sessions >= [v1.2.1](https://github.com/gorilla/sessions/releases/tag/v1.2.1)

[Traefik](https://traefik.io) plugins are developed using the compiled [Go language](https://golang.org). Rather than being pre-compiled and linked, however, plugins are executed on the fly by [Yaegi](https://github.com/traefik/yaegi), an embedded Go interpreter. Due to [traefik/yaegi#1275](https://github.com/traefik/yaegi/issues/1275), the `ldap-go` module only works after the listed version.

## Usage

### Add Plugin to Service

```yml
whoami:
  image: "traefik/whoami"
  container_name: "whoami"
  labels:
    - traefik.enable=true
    - traefik.http.routers.whoami.rule=Host(`whoami.localhost`)
    - traefik.http.routers.whoami.entrypoints=web
    # ldapAuth Register Middleware ====================================================
    - traefik.http.routers.whoami.middlewares=ldap_auth
    # ldapAuth Options=================================================================
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.enabled=true
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.logLevel=DEBUG
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.url=ldap://ldap.forumsys.com
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.port=389
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.baseDN=dc=example,dc=com
    - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.attribute=uid
    # =================================================================================
```

### Bind Mode Example

```yml
[...]
labels:
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.url=ldap://ldap.forumsys.com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.port=389
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.baseDN=dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.attribute=uid
```

### Search Mode Anonymous Example

```yml
[...]
labels:
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.url=ldap://ldap.forumsys.com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.port=389
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.baseDN=dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.attribute=uid
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.searchFilter=({{.Attribute}}={{.Username}})
```

### Search Mode Authenticated Example

```yml
[...]
labels:
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.url=ldap://ldap.forumsys.com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.port=389
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.baseDN=dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.attribute=uid
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.bindDN=uid=tesla,dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.bindPassword=password
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.searchFilter=({{.Attribute}}={{.Username}})
```

### Advanced Search Mode Example

```yml
[...]
labels:
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.url=ldap://ldap.forumsys.com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.port=389
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.baseDN=dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.attribute=uid
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.bindDN=uid=tesla,dc=example,dc=com
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.bindPassword=password
  - traefik.http.middlewares.ldap_auth.plugin.ldapAuth.searchFilter=(&(objectClass=person)({{.Attribute}}={{.Username}}))
```

## Operations Mode

### Bind Mode

If no `searchFilter` is specified in its configuration, the middleware runs in the default bind mode, meaning it tries to make a simple bind request to the LDAP server with the credentials provided in the request headers. If the bind succeeds, the middleware forwards the request, otherwise it returns a 401 Unauthorized status code.

### Search Mode

If a `searchFilter` query is specified in the configuration, then the middleware runs in search mode. In this mode, a search query with the given filter is issued to the LDAP server before trying to bind. If `bindDN` and `bindPassword` have also been provided, then the search query will use this crentials. If result of this search returns only `1` record, it tries to issue a bind request with this record, otherwise it aborts a 401 Unauthorized status code.

## Options

##### `enabled`

_Optional, Default: `true`_

Controls whether requests will be checked against LDAP or not before being delivered.

##### `logLevel`

_Optional, Default: `INFO`_

Set `LogLevel` for detailed information about plugin operation.

##### `url`

_Required, Default: `""`_

LDAP server address where queries will be performed.

##### `port`

_Optional, Default: `389`_

LDAP server port where queries will be performed.

##### `cacheTimeout`
_Optional, Default: `300`_

Indicates the number of `seconds` until the cookie session expires. A zero or negative number will expire the cookie immediately. 
##### `cacheCookieName`
_Optional, Default: `ldapAuth_session_token`_

The cookie session name.

##### `cacheKey`
Needs `traefik` >= [`v2.8.5`](https://github.com/traefik/traefik/releases/tag/v2.8.5)

_Optional, Default: `super-secret-key`_

The key used to criptography cookie session informations. You `must` use a strong value here.

##### `useTLS`
_Optional, Default: `false`_

Set to true if LDAP server should use an encrypted TLS connection, either with STARTTLS or LDAPS.

##### `startTLS`
_Optional, Default: `false`_

If set to true, instructs `ldapAuth` to issue a `StartTLS` request when initializing the connection with the LDAP server. This is not used if the `useTLS` option is set to `false`.

##### `certificateAuthority`
_Optional, Default: `""`_

The `certificateAuthority` option should contain one or more PEM-encoded certificates to use to establish a connection with the LDAP server if the connection uses TLS but that the certificate was signed by a custom Certificate Authority.


Example:
```yml
    certificateAuthority: |-
        -----BEGIN CERTIFICATE-----
        MIIB9TCCAWACAQAwgbgxGTAXBgNVBAoMEFF1b1ZhZGlzIExpbWl0ZWQxHDAaBgNV
        BAsME0RvY3VtZW50IERlcGFydG1lbnQxOTA3BgNVBAMMMFdoeSBhcmUgeW91IGRl
        Y29kaW5nIG1lPyAgVGhpcyBpcyBvbmx5IGEgdGVzdCEhITERMA8GA1UEBwwISGFt
        aWx0b24xETAPBgNVBAgMCFBlbWJyb2tlMQswCQYDVQQGEwJCTTEPMA0GCSqGSIb3
        DQEJARYAMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCJ9WRanG/fUvcfKiGl
        EL4aRLjGt537mZ28UU9/3eiJeJznNSOuNLnF+hmabAu7H0LT4K7EdqfF+XUZW/2j
        RKRYcvOUDGF9A7OjW7UfKk1In3+6QDCi7X34RE161jqoaJjrm/T18TOKcgkkhRzE
        apQnIDm0Ea/HVzX/PiSOGuertwIDAQABMAsGCSqGSIb3DQEBBQOBgQBzMJdAV4QP
        Awel8LzGx5uMOshezF/KfP67wJ93UW+N7zXY6AwPgoLj4Kjw+WtU684JL8Dtr9FX
        ozakE+8p06BpxegR4BR3FMHf6p+0jQxUEAkAyb/mVgm66TyghDGC6/YkiKoZptXQ
        98TwDIK/39WEB/V607As+KoYazQG8drorw==
        -----END CERTIFICATE-----
```

##### `insecureSkipVerify`
_Optional, Default: `false`_

When `useTLS` is enabled, the connection to the LDAP server is verified to be secure. This option allows `ldapAuth` to proceed and operate even for server connections otherwise considered insecure.

##### `attribute`

_Optional, Default: `cn`_

The attribute used to bind a user in [`Bind Mode`](#bind-mode). Bind queries use this pattern: `<attribute>=<username>,<baseDN>`, where the username is extracted from the request header. If [`AllowedGroups`](#allowedGroups) option was used in [`Bind Mode`](#bind-mode), the same pattern is added when searching if the user belongs to group.

##### `searchFilter`

_Optional, Default: `""`_

If not empty, the middleware will run in [`Search Mode`](#search-mode), filtering search results with the given query.

Filter queries can use the `{{.Option}}` format, from [text/template](https://pkg.go.dev/text/template#pkg-overview) go package, as placeholders that are replaced by the equivalent value from config. Additionaly, the username provided in the Authorization header of the request can also be used.

For example: `(&(objectClass=inetOrgPerson)(gidNumber=500)({{.Attribute}}={{.Username}}))`.

Will be replaced to: `(&(objectClass=inetOrgPerson)(gidNumber=500)(uid=tesla))`.

Note1: All filters options must be start with Uppercase to be replaced correctly.

Note2: `searchFilter` must **not** escape curly braces when using [labels](examples/conf-from-labels.yml).

Note3: `searchFilter` must escape curly braces when using [yml file](examples/dynamic-conf/ldapAuth-conf.yml).

Note4: `searchFilter` must escape curly braces when using [toml file](examples/dynamic-conf/ldapAuth-conf.toml).

##### `baseDN`

_Required, Default: `""`_

From where the plugin will search for users.

##### `bindDN`

_Optional, Default: `""`_

The domain name to bind to in order to authenticate to the LDAP server when running on [`Search Mode`](#search-mode). Leaving this empty with [`Search Mode`](#search-mode) means binds are anonymous, which is rarely expected behavior. It is not used when running in [`Bind Mode`](#bind-mode).

##### `bindPassword`

_Optional, Default: `""`_

The password corresponding to the `bindDN` specified when running in [`Search Mode`](#search-mode), used in order to authenticate to the LDAP server.

##### `forwardUsername`

_Optional, Default: `true`_

The `forwardUsername` option can be enabled to forward the username in a specific header, defined using the `forwardUsernameHeader` option.

##### `forwardUsernameHeader`

_Optional, Default: `Username`_

Name of the header to put the username in when forwarding it. This is not used if the `forwardUsername` option is set to `false`.

##### `forwardAuthorization`

_Optional, Default: `false`_

The `forwardAuthorization` option determines if the authorization header will be forwarded or stripped from the request after it has been approved by the middleware. `Attention`, enabling this option may expose the password of the LDAP user who is making the request.

##### `forwardExtraLDAPHeaders`

_Optional, Default: `false`_

The `forwardExtraLDAPHeaders` option determines if the LDAP Extra Headers, `Ldap-Extra-Attr-DN` and
`Ldap-Extra-Attr-CN`, will be added or not to request. This is not used if the `forwardUsername` option is set to `false` or if `searchFilter` is empty.

##### `wwwAuthenticateHeader`

_Optional, Default: `true`_

If the LDAP middleware receives a request with a missing or invalid Authorization header and `wwwAuthenticateHeader` is enabled, it will set a `WWW-Authenticate` header in the 401 Unauthorized response. See the [WWW-Authenticate header documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate) for more information.

##### `wwwAuthenticateHeaderRealm`

_Optional, Default: `""`_

The name of the realm to specify in the `WWW-Authenticate` header. This option is ineffective unless the `wwwAuthenticateHeader` option is set to true.

##### `allowedGroups`
Needs `traefik` >= [`v2.8.2`](https://github.com/traefik/traefik/releases/tag/v2.8.2)

_Optional, Default: `[]`_

The list of LDAP group DNs that users must be members of to be granted access. If a user is in any one of the listed groups, then that user is granted access.

If setted to an empty list will allow all users that have an LDAP account to log in, without performing any group membership checks.