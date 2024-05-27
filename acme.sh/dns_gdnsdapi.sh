#!/usr/bin/env sh

# GDNSD_API environment variable should be set to the URL of the gdnsd-acme-dns-api service
GDNSD_API=${GDNSD_API:-"http://localhost:8080"}
GDNSD_API_TOKEN=${GDNSD_API_TOKEN:-""}
dns_gdnsd_add() {
  fulldomain=$1
  txtvalue=$2
  $_post_url="${GDNSD_API}/acme-dns-01"
  $_postContentType="application/json"
  $_data="{\"${fulldomain}\": \"${txtvalue}\"}"
  export _H1="Authorization: Bearer ${GDNSD_API_TOKEN}"
  response="$(_post "$_data" "$_post_url" "" "POST" "$_postContentType")"
  exit_code="$?"
  return $exit_code
}

dns_gdnsd_rm() {
  fulldomain=$1
  txtvalue=$2
  fulldomain=$1
  txtvalue=$2
  $_post_url="${GDNSD_API}/acme-dns-01-flush"
  $_postContentType="application/json"
  $_data="{\"submit\": \"true\"}"
  export _H1="Authorization: Bearer ${GDNSD_API_TOKEN}"
  response="$(_post "$_data" "$_post_url" "" "POST" "$_postContentType")"
  exit_code="$?"
  return $exit_code
}

# acme.sh DNS API hook functions
case "$1" in
  add)
    dns_gdnsd_add "$2" "$3"
    ;;
  rm)
    dns_gdnsd_rm "$2" "$3"
    ;;
esac
