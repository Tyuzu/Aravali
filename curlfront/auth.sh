#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
AUTH_BASE="${API_BASE_URL%/}/auth"
COOKIE_JAR="${COOKIE_JAR:-$PWD/.curl-cookie-jar}"
mkdir -p "$(dirname "$COOKIE_JAR")"

USERNAME="${USERNAME:-demo_user}"
EMAIL="${EMAIL:-demo@example.com}"
PASSWORD="${PASSWORD:-StrongPass123!}"
OTP_CODE="${OTP_CODE:-123456}"

usage() {
  cat <<'EOF'
Usage:
  ./auth.sh register
  ./auth.sh login
  ./auth.sh refresh
  ./auth.sh logout
  ./auth.sh logout-all
  ./auth.sh request-otp
  ./auth.sh verify-otp

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  COOKIE_JAR     default: $PWD/.curl-cookie-jar
  USERNAME       default: demo_user
  EMAIL          default: demo@example.com
  PASSWORD       default: StrongPass123!
  AUTH_TOKEN     bearer token for protected routes
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  register)
    curl --silent --show-error \
      -X POST "$AUTH_BASE/register" \
      -H "Content-Type: application/json" \
      -d "{\"username\":\"$USERNAME\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}"
    ;;

  login)
    RESPONSE=$(curl --silent --show-error \
      -X POST "$AUTH_BASE/login" \
      -H "Content-Type: application/json" \
      -c "$COOKIE_JAR" \
      -b "$COOKIE_JAR" \
      -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
    echo "$RESPONSE"
    printf '%s\n' "$RESPONSE" | sed -n 's/.*"token":"\([^"]*\)".*/ACCESS_TOKEN=\1/p'
    ;;

  refresh)
    if [[ -z "${AUTH_TOKEN:-}" ]]; then
      echo "AUTH_TOKEN is not set; refresh uses the refresh cookie jar and the login response token is not required."
    fi
    curl --silent --show-error \
      -X POST "$AUTH_BASE/refresh" \
      -H "Content-Type: application/json" \
      -H "X-Refresh-Intent: 1" \
      -c "$COOKIE_JAR" \
      -b "$COOKIE_JAR"
    ;;

  logout)
    curl --silent --show-error \
      -X POST "$AUTH_BASE/logout" \
      -H "Content-Type: application/json" \
      -H "X-Refresh-Intent: 1" \
      -c "$COOKIE_JAR" \
      -b "$COOKIE_JAR"
    ;;

  logout-all)
    if [[ -z "${AUTH_TOKEN:-}" ]]; then
      echo "AUTH_TOKEN is required for logout-all."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$AUTH_BASE/logout-all" \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      -H "Content-Type: application/json" \
      -c "$COOKIE_JAR" \
      -b "$COOKIE_JAR"
    ;;

  request-otp)
    curl --silent --show-error \
      -X POST "$AUTH_BASE/request-otp" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$EMAIL\"}"
    ;;

  verify-otp)
    curl --silent --show-error \
      -X POST "$AUTH_BASE/verify-otp" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$EMAIL\",\"otp\":\"$OTP_CODE\"}"
    ;;

  *)
    usage
    exit 1
    ;;
esac
