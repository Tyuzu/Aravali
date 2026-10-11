#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
BAITO_BASE="${API_BASE_URL%/}/baitos"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ENTITY_TYPE="${ENTITY_TYPE:-farm}"
ENTITY_ID="${ENTITY_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./baito.sh latest
  ./baito.sh create
  ./baito.sh jobs [entityType] [entityId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for protected baito actions
  ENTITY_TYPE    default: farm
  ENTITY_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  latest)
    curl --silent --show-error "$BAITO_BASE/latest"
    ;;

  jobs)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$ENTITY_ID}"
    curl --silent --show-error "$API_BASE_URL/jobs/$TYPE/$ID"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$BAITO_BASE/baito" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"title":"Demo baito job","description":"Created via curlfront"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
