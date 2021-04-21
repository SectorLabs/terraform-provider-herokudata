#!/usr/bin/env bash

credential_name="$1"
addon_id="$2"
database_name="$3"
auth_token="$4"

if [[ -z $auth_token ]]; then
    auth_token="$(heroku auth:token)"
    echo ${auth_token}
fi

curl --request POST \
  --url https://data-api.heroku.com/graphql \
  --header "authorization: Bearer ${auth_token}" \
  --header "content-type: application/json" \
  --data '{
    "query": "mutation SetPostgresPermissions($addonUUID: ID!, $role: String!, $acls: [PostgresACLInput]!) {\n    setPostgresPermissions(\n      addon_uuid: $addonUUID\n      role: $role\n      acls: $acls\n    )\n  }",
    "variables": {
        "addonUUID": "'${addon_id}'",
        "role": "'${credential_name}'",
        "acls": [
            {
                "kind": "database",
                "name": "'${database_name}'",
                "privileges": ["CONNECT", "TEMPORARY"]
            }, {
                "kind": "table",
                "default": true,
                "privileges": ["SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE"]
            }, {
                "kind": "sequence",
                "default": true,
                "privileges": ["SELECT", "USAGE"]
            }, {
                "kind": "schema",
                "name": "public",
                "privileges": ["USAGE"]
            }, {
                "kind": "table",
                "schema": "public",
                "privileges": ["SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE"]
            }, {
                "kind": "sequence",
                "schema": "public",
                "privileges":["SELECT", "USAGE"]
            }
        ]
    }
}' | python -m json.tool
