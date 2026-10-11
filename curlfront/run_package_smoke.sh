#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."
export API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"

cmds=(
  "./curlfront/artists.sh list"
  "./curlfront/baito.sh latest"
  "./curlfront/booking.sh slots"
  "./curlfront/cart.sh get"
  "./curlfront/comments.sh list"
  "./curlfront/crops.sh list"
  "./curlfront/deliveries.sh list"
  "./curlfront/events.sh list"
  "./curlfront/faqs.sh list"
  "./curlfront/farms.sh list"
  "./curlfront/itinerary.sh list"
  "./curlfront/maps.sh map"
  "./curlfront/media.sh list"
  "./curlfront/menu.sh list"
  "./curlfront/merch.sh list"
  "./curlfront/notices.sh list"
  "./curlfront/pay.sh balance"
  "./curlfront/places.sh menu"
  "./curlfront/posts.sh list"
  "./curlfront/products.sh items"
  "./curlfront/recipes.sh list"
  "./curlfront/reviews.sh list"
  "./curlfront/search.sh ac"
  "./curlfront/settings.sh get"
  "./curlfront/tickets.sh list"
  "./curlfront/vendors.sh list"
  "./curlfront/workers.sh create-profile"
)

for cmd in "${cmds[@]}"; do
  echo "=== $cmd ==="
  bash -lc "$cmd"
  echo
  done
