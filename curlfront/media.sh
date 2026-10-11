#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ENTITY_TYPE="${ENTITY_TYPE:-event}"
ENTITY_ID="${ENTITY_ID:-1}"
MEDIA_ID="${MEDIA_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./media.sh list [entityType] [entityId]
  ./media.sh get [entityType] [entityId] [mediaId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  ENTITY_TYPE    default: event
  ENTITY_ID      default: 1
  MEDIA_ID      default: 1
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
    curl --silent --show-error "$API_BASE_URL/media/$TYPE/$ID"
    ;;

  get)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$ENTITY_ID}"
    MEDIA="${4:-$MEDIA_ID}"
    curl --silent --show-error "$API_BASE_URL/media/$TYPE/$ID/$MEDIA"
    ;;

  *)
    usage
    exit 1
    ;;
esac
