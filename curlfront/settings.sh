#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"

usage() {
  cat <<'EOF'
Usage:
  ./settings.sh get
  ./settings.sh schema
  ./settings.sh update

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for settings routes
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  get)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE_URL/settings"
    ;;

  schema)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE_URL/settings/schema"
    ;;

  update)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X PATCH "$API_BASE_URL/settings" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"theme":"dark"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
