Terraform Provider HerokuData
==================

This provider is used to manage Heroku Data resources (i.e., Postgres add-ons) via the https://data-api.heroku.com endpoint.
IMPORTANT: this endpoint is not officially documented and is subject to change, so this provider should not be treated as stable.

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.12.x
-   [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

```sh
$ make build
```

or

```sh
$ go build -o terraform-provider-herokudata
```

Using the provider
----------------------

Example to create credential:

```hcl
provider "herokudata" {
  api_key = "<API_KEY>"
}

resource "herokudata_credential" "my_credential" {
  addon_id = "<ADDON_ID>"
  name = "<CREDENTIAL_NAME>"
  permission = "readonly|readwrite"
}
```
