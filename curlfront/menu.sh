#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
PLACE_ID="${PLACE_ID:-1}"
MENU_ID="${MENU_ID:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./menu.sh list [placeId]
  ./menu.sh get [placeId] [menuId]
  ./menu.sh stock [placeId] [menuId]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  PLACE_ID       default: 1
  MENU_ID        default: 1
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  list)
    ID="${2:-$PLACE_ID}"
    curl --silent --show-error "$API_BASE_URL/places/menu/$ID"
    ;;

  get)
    PLACE="${2:-$PLACE_ID}"
    MENU="${3:-$MENU_ID}"
    curl --silent --show-error "$API_BASE_URL/places/menu/$PLACE/$MENU"
    ;;

  stock)
    PLACE="${2:-$PLACE_ID}"
    MENU="${3:-$MENU_ID}"
    curl --silent --show-error "$API_BASE_URL/places/menu/$PLACE/$MENU/stock"
    ;;

  *)
    usage
    exit 1
    ;;
esac
