#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
PLACE_ID="${PLACE_ID:-1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"

usage() {
  cat <<'EOF'
Usage:
  ./places.sh menu [placeId]
  ./places.sh rooms [placeId]
  ./places.sh services [placeId]
  ./places.sh products [placeId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  PLACE_ID       default: 1
  AUTH_TOKEN     bearer token for protected place writes
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  menu)
    ID="${2:-$PLACE_ID}"
    curl --silent --show-error "$API_BASE_URL/place/$ID/menu"
    ;;

  rooms)
    ID="${2:-$PLACE_ID}"
    curl --silent --show-error "$API_BASE_URL/place/$ID/rooms"
    ;;

  services)
    ID="${2:-$PLACE_ID}"
    curl --silent --show-error "$API_BASE_URL/place/$ID/services"
    ;;

  products)
    ID="${2:-$PLACE_ID}"
    curl --silent --show-error "$API_BASE_URL/place/$ID/products"
    ;;

  *)
    usage
    exit 1
    ;;
esac
