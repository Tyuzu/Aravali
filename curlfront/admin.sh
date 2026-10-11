#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
ADMIN_BASE="${API_BASE_URL%/}/admin"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"

usage() {
  cat <<'EOF'
Usage:
  ./admin.sh my-role-requests
  ./admin.sh role-requests
  ./admin.sh moderator-apply

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for protected admin routes
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  my-role-requests)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$ADMIN_BASE/role/requests/me"
    ;;

  role-requests)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$ADMIN_BASE/role/requests"
    ;;

  moderator-apply)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required."
      exit 1
    fi
    curl --silent --show-error \
      -X POST "$ADMIN_BASE/moderator/apply" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"message":"I would like to moderate this community."}'
    ;;

  *)
    usage
    exit 1
    ;;
esac
