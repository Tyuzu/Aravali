#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
PROFILE_BASE="${API_BASE_URL%/}/profile"
TOKEN="${AUTH_TOKEN:-${ACCESS_TOKEN:-}}"
USERNAME="${USERNAME:-demo_user}"
TARGET_USERNAME="${TARGET_USERNAME:-$USERNAME}"

usage() {
  cat <<'EOF'
Usage:
  ./profile.sh get
  ./profile.sh edit
  ./profile.sh delete
  ./profile.sh user <username>

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  AUTH_TOKEN     bearer token for authenticated profile routes
  USERNAME       default: demo_user
  TARGET_USERNAME target user for public profile lookup
  PROFILE_NAME   name field for edit
  PROFILE_BIO    bio field for edit
  PROFILE_EMAIL  email field for edit
  PROFILE_PHONE  phone_number field for edit
  PROFILE_PASSWORD password field for edit
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  get)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required for profile GET."
      exit 1
    fi
    curl --silent --show-error \
      -H "Authorization: Bearer $TOKEN" \
      "$PROFILE_BASE/profile"
    ;;

  edit)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required for profile edit."
      exit 1
    fi

    FORM_ARGS=()
    if [[ -n "${PROFILE_NAME:-}" ]]; then
      FORM_ARGS+=( -F "name=${PROFILE_NAME}" )
    fi
    if [[ -n "${PROFILE_BIO:-}" ]]; then
      FORM_ARGS+=( -F "bio=${PROFILE_BIO}" )
    fi
    if [[ -n "${PROFILE_EMAIL:-}" ]]; then
      FORM_ARGS+=( -F "email=${PROFILE_EMAIL}" )
    fi
    if [[ -n "${PROFILE_PHONE:-}" ]]; then
      FORM_ARGS+=( -F "phone_number=${PROFILE_PHONE}" )
    fi
    if [[ -n "${PROFILE_PASSWORD:-}" ]]; then
      FORM_ARGS+=( -F "password=${PROFILE_PASSWORD}" )
    fi
    if [[ ${#FORM_ARGS[@]} -eq 0 ]]; then
      FORM_ARGS+=( -F "bio=Updated bio from curlfront" )
    fi

    curl --silent --show-error \
      -X PUT \
      -H "Authorization: Bearer $TOKEN" \
      "${FORM_ARGS[@]}" \
      "$PROFILE_BASE/edit"
    ;;

  delete)
    if [[ -z "$TOKEN" ]]; then
      echo "AUTH_TOKEN or ACCESS_TOKEN is required for profile delete."
      exit 1
    fi
    curl --silent --show-error \
      -X DELETE \
      -H "Authorization: Bearer $TOKEN" \
      "$PROFILE_BASE/delete"
    ;;

  user)
    USER="${2:-$TARGET_USERNAME}"
    curl --silent --show-error \
      "$API_BASE_URL/user/$USER"
    ;;

  *)
    usage
    exit 1
    ;;
esac
