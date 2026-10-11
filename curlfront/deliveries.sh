#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
DELIVERIES_BASE="${API_BASE_URL%/}/deliveries"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
DELIVERY_ID="${DELIVERY_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./deliveries.sh list
  ./deliveries.sh get [deliveryId]
  ./deliveries.sh status [deliveryId]
  ./deliveries.sh create

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for delivery routes
  DELIVERY_ID    default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$DELIVERIES_BASE"
    ;;

  get)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$DELIVERY_ID}"
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$DELIVERIES_BASE/$ID"
    ;;

  status)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$DELIVERY_ID}"
    curl --silent --show-error \
      -X PATCH "$DELIVERIES_BASE/$ID/status" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"status":"picked_up"}'
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$DELIVERIES_BASE" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"order_id":1,"destination":"Demo destination"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
