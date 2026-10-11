#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"

usage() {
  cat <<'EOF'
Usage:
  ./recipes.sh tags
  ./recipes.sh list

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     optional bearer token for authenticated recipe queries
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  tags)
    curl --silent --show-error "$API_BASE_URL/recipes/tags"
    ;;

  list)
    if [[ -n "$TOKEN" ]]; then
      curl --silent --show-error \
        -H "Authorization: Bearer $TOKEN" \
        "$API_BASE_URL/recipes"
    else
      curl --silent --show-error "$API_BASE_URL/recipes"
    fi
    ;;

  *)
    usage
    exit 1
    ;;
esac
