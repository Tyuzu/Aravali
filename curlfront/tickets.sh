#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
EVENT_ID="${EVENT_ID:-1}"
TICKET_ID="${TICKET_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./tickets.sh list [eventId]
  ./tickets.sh get [eventId] [ticketId]
  ./tickets.sh create [eventId]
  ./tickets.sh verify [eventId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for ticket writes
  EVENT_ID       default: 1
  TICKET_ID      default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    ID="${2:-$EVENT_ID}"
    curl --silent --show-error "$API_BASE_URL/ticket/event/$ID"
    ;;

  get)
    EVENT="${2:-$EVENT_ID}"
    TICKET="${3:-$TICKET_ID}"
    curl --silent --show-error "$API_BASE_URL/ticket/event/$EVENT/$TICKET"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$EVENT_ID}"
    curl --silent --show-error \
      -X POST "$API_BASE_URL/ticket/event/$ID" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"General Admission","price":75}'
    ;;

  verify)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    ID="${2:-$EVENT_ID}"
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE_URL/ticket/verify/$ID"
    ;;

  *)
    usage
    exit 1
    ;;
esac
