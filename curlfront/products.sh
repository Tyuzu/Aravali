#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
ENTITY_TYPE="${ENTITY_TYPE:-farm}"
ENTITY_ID="${ENTITY_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./products.sh items
  ./products.sh categories
  ./products.sh orders-mine
  ./products.sh product-details [entityType] [entityId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for order-related actions
  ENTITY_TYPE    default: farm
  ENTITY_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  items)
    curl --silent --show-error "$API_BASE_URL/farm/items"
    ;;

  categories)
    curl --silent --show-error "$API_BASE_URL/farm/items/categories"
    ;;

  orders-mine)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE_URL/orders/mine"
    ;;

  product-details)
    TYPE="${2:-$ENTITY_TYPE}"
    ID="${3:-$ENTITY_ID}"
    curl --silent --show-error "$API_BASE_URL/products/$TYPE/$ID"
    ;;

  *)
    usage
    exit 1
    ;;
esac
