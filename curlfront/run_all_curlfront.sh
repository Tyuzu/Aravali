#!/usr/bin/env bash
set +e

cd "$(dirname "$0")/.."
export API_BASE_URL="${API_BASE_URL:-http://localhost:4000/api/v1}"
export COOKIE_JAR="${COOKIE_JAR:-$PWD/.curl-cookie-jar}"
mkdir -p "$(dirname "$COOKIE_JAR")"

./curlfront/auth.sh register || true
AUTH_TOKEN=$(./curlfront/auth.sh login | sed -n 's/.*"token":"\([^"]*\)".*/\1/p' || true)
export AUTH_TOKEN

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
  "./curlfront/profile.sh get"
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
  bash -lc "$cmd" || true
  echo
  sleep 1
done
