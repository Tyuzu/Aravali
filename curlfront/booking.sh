#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
BOOKING_BASE="${API_BASE_URL%/}/bookings"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
BOOKING_ID="${BOOKING_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./booking.sh slots
  ./booking.sh create-slot
  ./booking.sh bookings
  ./booking.sh create-booking

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for booking routes
  BOOKING_ID     default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  slots)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$BOOKING_BASE/slots"
    ;;

  create-slot)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$BOOKING_BASE/slots" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"date":"2026-10-12","start_time":"09:00","end_time":"11:00","capacity":10}'
    ;;

  bookings)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$BOOKING_BASE/bookings"
    ;;

  create-booking)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$BOOKING_BASE/bookings" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"slot_id":1,"guest_count":2}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
