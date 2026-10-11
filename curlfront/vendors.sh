#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
VENDOR_ID="${VENDOR_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./vendors.sh list
  ./vendors.sh me
  ./vendors.sh get [vendorId]
  ./vendors.sh create

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for vendor routes
  VENDOR_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$API_BASE_URL/vendors"
    ;;

  me)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE_URL/vendors/me"
    ;;

  get)
    ID="${2:-$VENDOR_ID}"
    curl --silent --show-error "$API_BASE_URL/vendors/vendor/$ID"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$API_BASE_URL/vendors" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo Vendor","category":"food"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
