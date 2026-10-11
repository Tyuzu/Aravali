#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ENTITY_TYPE="${ENTITY_TYPE:-event}"
ENTITY_ID="${ENTITY_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./notices.sh list [entityType] [entityId]
  ./notices.sh create [entityType] [entityId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for notice creation
  ENTITY_TYPE    default: event
  ENTITY_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$ENTITY_ID}"
    curl --silent --show-error "$API_BASE_URL/notices/$TYPE/$ID"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$ENTITY_ID}"
    curl --silent --show-error \
      -X POST "$API_BASE_URL/notices/$TYPE/$ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"message":"New update from curlfront."}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
