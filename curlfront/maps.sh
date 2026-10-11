#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"

usage() {
  cat <<'EOF'
Usage:
  ./maps.sh map
  ./maps.sh distance [from] [to]

Environment:
  API_BASE_URL   default: http://localhost:4000/api/v1
  MAP_FROM      default: 10.0,20.0
  MAP_TO        default: 11.0,21.0
EOF
}

if [[ ${1:-} == "" ]]; then
  usage
  exit 1
fi

case "$1" in
  map)
    curl --silent --show-error "$API_BASE_URL/gta/map"
    ;;

  distance)
    FROM="${2:-${MAP_FROM:-10.0,20.0}}"
    TO="${3:-${MAP_TO:-11.0,21.0}}"
    curl --silent --show-error "$API_BASE_URL/gta/map/distance?from=$FROM&to=$TO"
    ;;

  *)
    usage
    exit 1
    ;;
esac
