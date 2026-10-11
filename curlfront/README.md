# Curlfront

This folder contains small, runnable curl scripts for the main API surfaces that are commonly used during local development and smoke testing.

## Quick start

Set the backend URL if it is not running on the default local port:

```bash
export API_BASE_URL="http://localhost:4000/api/v1"
```

For auth flows, the server also stores a refresh token in a cookie jar. The scripts use a local cookie file by default:

```bash
export COOKIE_JAR="$PWD/.curl-cookie-jar"
```

## Auth endpoints

- `auth.sh register`
- `auth.sh login`
- `auth.sh refresh`
- `auth.sh logout`
- `auth.sh logout-all`
- `auth.sh request-otp`
- `auth.sh verify-otp`

## Profile endpoints

- `profile.sh get`
- `profile.sh edit`
- `profile.sh delete`
- `profile.sh user <username>`

## Additional package endpoints

- `admin.sh my-role-requests`
- `artists.sh list`
- `baito.sh latest`
- `booking.sh slots`
- `cart.sh get`
- `comments.sh list`
- `crops.sh list`
- `deliveries.sh list`
- `events.sh list`
- `faqs.sh list`
- `farms.sh list`
- `itinerary.sh list`
- `maps.sh map`
- `media.sh list`
- `menu.sh list`
- `merch.sh list`
- `notices.sh list`
- `pay.sh balance`
- `places.sh menu`
- `posts.sh list`
- `products.sh items`
- `recipes.sh list`
- `reviews.sh list`
- `search.sh ac`
- `settings.sh get`
- `tickets.sh list`
- `vendors.sh list`
- `workers.sh create-profile`

## Examples

```bash
cd curlfront
./auth.sh register
./auth.sh login
./profile.sh get
```

The login script prints the access token so you can reuse it:

```bash
export AUTH_TOKEN=$(./auth.sh login | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
```

The profile edit flow accepts multipart form fields such as `name`, `email`, `bio`, `phone_number`, and `password`.
