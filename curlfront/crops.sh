#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
FARM_ID="${FARM_ID:-1}"
CROP_ID="${CROP_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./crops.sh list
  ./crops.sh add [farmId]
  ./crops.sh edit [farmId] [cropId]
  ./crops.sh delete [farmId] [cropId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for crop management
  FARM_ID        default: 1
  CROP_ID       default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$API_BASE_URL/crops"
    ;;

  add)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$FARM_ID}"
    curl --silent --show-error \
      -X POST "$API_BASE_URL/farms/farm/$ID/crops" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo crop","variety":"Hybrid","price":420}'
    ;;

  edit)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    FARM="${2:-$FARM_ID}"
    CROP="${3:-$CROP_ID}"
    curl --silent --show-error \
      -X PUT "$API_BASE_URL/farms/farm/$FARM/crops/$CROP" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Updated crop","price":460}'
    ;;

  delete)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    FARM="${2:-$FARM_ID}"
    CROP="${3:-$CROP_ID}"
    curl --silent --show-error \
      -X DELETE "$API_BASE_URL/farms/farm/$FARM/crops/$CROP" \
      -H "Authorization: Bearer $TOKEN"
    ;;

  *)
    usage
    exit 1
    ;;
esac
