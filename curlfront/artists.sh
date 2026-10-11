#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
ARTISTS_BASE="${API_BASE_URL%/}/artists"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ARTIST_ID="${ARTIST_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./artists.sh list
  ./artists.sh get [id]
  ./artists.sh create
  ./artists.sh events [artistId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for authenticated artist writes
  ARTIST_ID     default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$ARTISTS_BASE"
    ;;

  get)
    ID="${2:-$ARTIST_ID}"
    curl --silent --show-error "$ARTISTS_BASE/$ID"
    ;;

  events)
    ID="${2:-$ARTIST_ID}"
    curl --silent --show-error "$ARTISTS_BASE/$ID/events"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$ARTISTS_BASE" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo Artist","bio":"Created from curlfront."}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
