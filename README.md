# vgs-api-client

VGS API Client

## Build

``make build``

## Test

``make test``

## Supported commands

### Organizations

**get organizations**
```console
vgs-api-client get organizations

# output
ACkrF9yRThxLkvABGZVK93jW	required_idp      	INACTIVE  	2020-11-06T17:04:49.615Z
AC3g11N5RmChXA1uMGsFfRyj	Vault mgmt testing	INCOMPLETE	2019-08-01T11:40:33.109Z
ACd8St43Br39HWZj4xDyyzoY	Watchdog          	INACTIVE  	2019-08-15T21:13:06.649Z
```

**describe organization <org_id>**
```console
vgs-api-client describe organization AC3g11N5RmChXA1uMGsFfRyj

# output
{
  "id": "AC3g11N5RmChXA1uMGsFfRyj",
  "name": "Vault mgmt testing",
  "state": "INCOMPLETE",
  "created_at": "2019-08-01T11:40:33.109Z",
  "updated_at": "2019-08-02T14:24:10.739Z",
  "users": null,
  "environments": [
    {
      "id": "ENVxA3SUxU4y2JpcGcCTZziZJ",
      "name": "Live",
      "identifier": "LIVE",
      "region": "US"
    },
    {
      "id": "ENVwPNhoFaaNDNACsvviESU7H",
      "name": "Sandbox",
      "identifier": "SANDBOX",
      "region": "US"
    }
  ]
}
```

### Vaults

**get vaults <org_id>**

```console
vgs-api-client get vaults AC3g11N5RmChXA1uMGsFfRyj

# output
tntol44iuix	Watchdog Provision Test 2021-02-03-11	SANDBOX	2021-02-03T11:04:19.740Z
tnt2a71fm9g	Watchdog Provision Test 2021-01-31-18	SANDBOX	2021-01-31T18:03:19.658Z
tntiyvzmaqu	Watchdog Provision Test 2021-01-28-20	SANDBOX	2021-01-28T20:02:19.730Z
...
```

**describe vault <vault_id>**

```console
vgs-api-client describe vaults tntog8icclh

# output
{
  "id": "tntog8icclh",
  "name": "Watchdog Provision Test 2021-02-14-17",
  "environment": "SANDBOX",
  "state": "PROVISIONED",
  "created_at": "2021-02-14T17:04:19.728Z",
  "updated_at": "2021-02-14T17:04:19.728Z"
}
```

**create vault <org_id>**
```console
vgs-api-client create vault AC3g11N5RmChXA1uMGsFfRyj --vault "zinovii 01"

# output
{
  "id": "tntkboebyfe",
  "name": "zinovii 01",
  "environment": "SANDBOX",
  "state": "PROVISIONED",
  "created_at": "2021-02-14T22:10:09.652Z",
  "updated_at": "2021-02-14T22:10:09.652Z"
}
```

**delete vault <vault_id>**

```console
vgs-api-client delete vault tntvwyrkd5j
```