#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
EVENTS_BASE="${API_BASE_URL%/}/events"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
EVENT_ID="${EVENT_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./events.sh list
  ./events.sh count
  ./events.sh create
  ./events.sh get [eventId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for event creation and updates
  EVENT_ID       default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$EVENTS_BASE/events"
    ;;

  count)
    curl --silent --show-error "$EVENTS_BASE/events/count"
    ;;

  get)
    ID="${2:-$EVENT_ID}"
    curl --silent --show-error "$EVENTS_BASE/event/$ID"
    ;;

  create)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$EVENTS_BASE/event" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"title":"Demo Event","description":"Created from curlfront","date":"2026-10-15"}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
