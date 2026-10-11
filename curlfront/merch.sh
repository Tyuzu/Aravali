#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ENTITY_TYPE="${ENTITY_TYPE:-event}"
EVENT_ID="${EVENT_ID:-1}"
MERCH_ID="${MERCH_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./merch.sh list [entityType] [eventId]
  ./merch.sh get [entityType] [eventId] [merchId]
  ./merch.sh create [entityType] [eventId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for merch writes
  ENTITY_TYPE    default: event
  EVENT_ID       default: 1
  MERCH_ID       default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$EVENT_ID}"
    curl --silent --show-error "$API_BASE_URL/merch/$TYPE/$ID"
    ;;

  get)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$EVENT_ID}"
    MERCH="${4:-$MERCH_ID}"
    curl --silent --show-error "$API_BASE_URL/merch/$TYPE/$ID/$MERCH"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$EVENT_ID}"
    curl --silent --show-error \
      -X POST "$API_BASE_URL/merch/$TYPE/$ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo merch","price":250}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
