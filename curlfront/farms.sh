#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
FARMS_BASE="${API_BASE_URL%/}/farms"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
FARM_ID="${FARM_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./farms.sh list
  ./farms.sh get [farmId]
  ./farms.sh create

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for farm writes
  FARM_ID        default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$FARMS_BASE"
    ;;

  get)
    ID="${2:-$FARM_ID}"
    curl --silent --show-error "$FARMS_BASE/farm/$ID"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$FARMS_BASE" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo Farm","location":"Remote valley"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
