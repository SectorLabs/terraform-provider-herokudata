#!/usr/bin/env bash

credential_name="$1"
addon_id="$2"
auth_token="$3"

if [[ -z $auth_token ]]; then
    auth_token="$(heroku auth:token)"
    echo ${auth_token}
fi


curl --request POST \
  --url https://data-api.heroku.com/graphql \
  --header "authorization: Bearer ${auth_token}" \
  --header "content-type: application/json" \
  --data '{
    "query":"mutation CreateCredential($addonUUID: ID!, $name: String!) {    createCredential(\n      addon_uuid: $addonUUID\n      name: $name\n    )\n  }","variables":{"addonUUID":"'${addon_id}'","name":"'${credential_name}'"}}
' | python -m json.tool
