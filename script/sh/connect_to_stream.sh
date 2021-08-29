#!/bin/zsh

curl -X GET -H "Authorization: Bearer $BEARER_TOKEN" "https://api.twitter.com/2/tweets/search/stream?tweet.fields=created_at&expansions=author_id&user.fields=created_at"
