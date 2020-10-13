#!/usr/bin/env bash

addon_id="$1"
auth_token="$2"

if [[ -z $auth_token ]]; then
    auth_token="$(heroku auth:token)"
    echo ${auth_token}
fi

curl --request POST \
  --url https://data-api.heroku.com/graphql \
  --header "authorization: Bearer ${auth_token}" \
  --header "content-type: application/json" \
  --data '{
    "query":"query FetchPostgresDetails($addonUUID: ID!) { postgres (addon_uuid: $addonUUID) { credentials_list {name uuid state}}}","variables":{"addonUUID":"'${addon_id}'"}
}
' | python -m json.tool


