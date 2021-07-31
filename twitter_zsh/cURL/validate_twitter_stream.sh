#!/bin/zsh

curl -X GET 'https://api.twitter.com/2/tweets/search/stream/rules' -H "Authorization: Bearer $BEARER_TOKEN"
