#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
ITINERARY_BASE="${API_BASE_URL%/}/itineraries"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ITINERARY_ID="${ITINERARY_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./itinerary.sh list
  ./itinerary.sh get [id]
  ./itinerary.sh search [query]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for protected itinerary actions, if needed
  ITINERARY_ID  default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$ITINERARY_BASE"
    ;;

  get)
    ID="${2:-$ITINERARY_ID}"
    curl --silent --show-error "$ITINERARY_BASE/all/$ID"
    ;;

  search)
    QUERY="${2:-demo}"
    curl --silent --show-error "$ITINERARY_BASE/search?q=$QUERY"
    ;;

  *)
    usage
    exit 1
    ;;
esac
