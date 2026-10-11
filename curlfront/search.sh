#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TAB_ID="${TAB_ID:-all}"

usage() {
  cat <<'EOF'
Usage:
  ./search.sh ac
  ./search.sh query [tabId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  TAB_ID        default: all
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  ac)
    curl --silent --show-error "$API_BASE_URL/ac"
    ;;

  query)
    TAB="${2:-$TAB_ID}"
    curl --silent --show-error "$API_BASE_URL/search/$TAB?q=demo"
    ;;

  *)
    usage
    exit 1
    ;;
esac
