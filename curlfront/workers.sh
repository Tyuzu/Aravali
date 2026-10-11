#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
WORKER_ID="${WORKER_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./workers.sh create-profile
  ./workers.sh update-profile [workerId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for worker profile routes
  WORKER_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  create-profile)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$API_BASE_URL/baitos/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Demo worker","skills":["farming","delivery"]}'
    ;;

  update-profile)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$WORKER_ID}"
    curl --silent --show-error \
      -X PATCH "$API_BASE_URL/baitos/profile/$ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"skills":["farming","delivery","event-support"]}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
