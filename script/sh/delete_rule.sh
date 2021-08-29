#!/bin/zsh

curl -X POST 'https://api.twitter.com/2/tweets/search/stream/rules' -H 'Content-type: application/json' -H "Authorization: Bearer $BEARER_TOKEN" -d "{ "delete": { "ids": ["$1"]}}"
