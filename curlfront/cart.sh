#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
CART_BASE="${API_BASE_URL%/}/cart"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"

usage() {
  cat <<'EOF'
Usage:
  ./cart.sh get
  ./cart.sh add
  ./cart.sh update
  ./cart.sh remove
  ./cart.sh clear

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for cart routes
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  get)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$CART_BASE"
    ;;

  add)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$CART_BASE" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"product_id":1,"quantity":1}'
    ;;

  update)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$CART_BASE/update" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"product_id":1,"quantity":2}'
    ;;

  remove)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X DELETE "$CART_BASE/item" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"product_id":1}'
    ;;

  clear)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X DELETE "$CART_BASE" \
      -H "Authorization: Bearer $TOKEN"
    ;;

  *)
    usage
    exit 1
    ;;
esac
