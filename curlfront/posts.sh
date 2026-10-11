#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
POST_ID="${POST_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./posts.sh list
  ./posts.sh get [postId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  POST_ID       default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    curl --silent --show-error "$API_BASE_URL/posts"
    ;;

  get)
    ID="${2:-$POST_ID}"
    curl --silent --show-error "$API_BASE_URL/posts/post/$ID"
    ;;

  *)
    usage
    exit 1
    ;;
esac
