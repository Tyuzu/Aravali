from __future__ import annotations

import base64
import hashlib
import hmac
import json
import os
import time
from typing import Any

from fastapi import FastAPI, Header, HTTPException, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, EmailStr

try:
    import bcrypt
except Exception:  # pragma: no cover - optional dependency, fallback to PBKDF2
    bcrypt = None

from routes import route_registry, router


JWT_SECRET = "aravali-signature"
DEFAULT_TOKEN_TTL_SECONDS = 60 * 60


class RegisterPayload(BaseModel):
    username: str
    email: EmailStr
    password: str


class LoginPayload(BaseModel):
    username: str
    password: str


USERS: dict[str, dict[str, Any]] = {}
ACTIVE_TOKENS: dict[str, str] = {}


def _jwt_payload(user: dict[str, Any]) -> dict[str, Any]:
    return {
        "userid": str(user["id"]),
        "userID": str(user["id"]),
        "sub": str(user["id"]),
        "username": user["username"],
        "email": user["email"],
        "roles": user.get("roles", ["user"]),
        "permissions": user.get("permissions", ["read"]),
        "exp": int(time.time()) + DEFAULT_TOKEN_TTL_SECONDS,
    }


def _b64url(value: bytes) -> str:
    return base64.urlsafe_b64encode(value).rstrip(b"=").decode("utf-8")


def _generate_refresh_token() -> str:
    return hashlib.sha256(os.urandom(32)).hexdigest()


def _sign_hs256(data: bytes) -> bytes:
    return hmac.new(JWT_SECRET.encode("utf-8"), data, hashlib.sha256).digest()


def _make_token(user: dict[str, Any]) -> str:
    header = {"alg": "HS256", "typ": "JWT"}
    payload = _jwt_payload(user)
    encoded_header = _b64url(json.dumps(header, separators=(",", ":")).encode("utf-8"))
    encoded_payload = _b64url(json.dumps(payload, separators=(",", ":")).encode("utf-8"))
    signing_input = f"{encoded_header}.{encoded_payload}".encode("utf-8")
    signature = _b64url(_sign_hs256(signing_input))
    return f"{encoded_header}.{encoded_payload}.{signature}"


def _hash_password(password: str) -> str:
    if bcrypt is not None:
        return bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")

    salt = os.urandom(16)
    derived = hashlib.pbkdf2_hmac("sha256", password.encode("utf-8"), salt, 100_000)
    return f"pbkdf2_sha256${salt.hex()}${derived.hex()}"


def _verify_password(password: str, user: dict[str, Any]) -> bool:
    stored_password = user.get("password")
    if isinstance(stored_password, str) and stored_password == password:
        return True

    raw_hash = user.get("password_hash") or user.get("PasswordHash")
    if not isinstance(raw_hash, str) or not raw_hash:
        return False

    if bcrypt is not None:
        try:
            if raw_hash.startswith("$2"):
                return bcrypt.checkpw(password.encode("utf-8"), raw_hash.encode("utf-8"))
        except ValueError:
            return False

    if raw_hash.startswith("pbkdf2_sha256$"):
        try:
            _, salt_hex, digest_hex = raw_hash.split("$", 2)
            salt = bytes.fromhex(salt_hex)
            expected = bytes.fromhex(digest_hex)
            computed = hashlib.pbkdf2_hmac("sha256", password.encode("utf-8"), salt, 100_000)
            return hmac.compare_digest(computed, expected)
        except (ValueError, TypeError):
            return False

    return False


def _user_response(user: dict[str, Any], token: str | None = None, refresh_token: str | None = None) -> dict[str, Any]:
    if token is None:
        token = _make_token(user)
    if refresh_token is None:
        refresh_token = _generate_refresh_token()

    ACTIVE_TOKENS[token] = user["username"]
    ACTIVE_TOKENS[refresh_token] = user["username"]

    return {
        "message": "Success",
        "status": 200,
        "token": token,
        "Token": token,
        "refreshToken": refresh_token,
        "refresh_token": refresh_token,
        "userid": str(user["id"]),
        "UserID": str(user["id"]),
        "username": user["username"],
        "roles": user.get("roles", ["user"]),
        "role": user.get("roles", ["user"]),
        "permissions": user.get("permissions", ["read"]),
        "user": {
            "id": user["id"],
            "userid": str(user["id"]),
            "userId": str(user["id"]),
            "username": user["username"],
            "email": user["email"],
            "roles": user.get("roles", ["user"]),
            "permissions": user.get("permissions", ["read"]),
        },
    }


app = FastAPI(
    title="Aravali",
    version="1.0.0",
    description="FastAPI backend exposing the frontend route catalog used by the Aravali app.",
)

app.add_middleware(
    CORSMiddleware,
    allow_origin_regex=r"https?://(localhost|127\.0\.0\.1):(5173|3000|4173)",
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Expose the same route catalog under the frontend API prefix as well as root.
app.include_router(router, prefix="/api/v1")
app.include_router(router)


@app.get("/health", tags=["system"], summary="Health check")
async def health_check() -> dict[str, Any]:
    return {
        "status": "ok",
        "service": "aravali",
        "routes": len(route_registry),
    }


@app.get("/api/v1/health", tags=["system"], summary="Health check under the frontend API path")
async def api_health_check() -> dict[str, Any]:
    return await health_check()


@app.post("/register")
@app.post("/auth/register")
@app.post("/api/v1/register")
@app.post("/api/v1/auth/register")
async def register_user(payload: RegisterPayload) -> dict[str, Any]:
    username = payload.username.strip()
    email = payload.email.lower().strip()
    password = payload.password.strip()

    if not username or not email or not password:
        raise HTTPException(status_code=400, detail="Username, email, and password are required.")

    if len(password) < 6:
        raise HTTPException(status_code=400, detail="Password must be at least 6 characters long.")

    existing = next((u for u in USERS.values() if u["username"].lower() == username.lower() or u["email"].lower() == email), None)
    if existing:
        raise HTTPException(status_code=409, detail="A user with that username or email already exists.")

    password_hash = _hash_password(password)
    user = {
        "id": str(len(USERS) + 1),
        "username": username,
        "email": email,
        "password": "",
        "password_hash": password_hash,
        "PasswordHash": password_hash,
        "roles": ["user"],
        "permissions": ["read"],
    }
    USERS[username] = user
    token = _make_token(user)
    response = _user_response(user, token)
    response["message"] = "User registered successfully."
    response["status"] = 201
    return response


@app.post("/login")
@app.post("/auth/login")
@app.post("/api/v1/login")
@app.post("/api/v1/auth/login")
async def login_user(payload: LoginPayload, response: Response) -> dict[str, Any]:
    username = payload.username.strip()
    password = payload.password.strip()

    if not username or not password:
        raise HTTPException(status_code=400, detail="Username and password are required.")

    user = USERS.get(username) or next((u for u in USERS.values() if u["username"].lower() == username.lower()), None)
    if not user or not _verify_password(password, user):
        raise HTTPException(status_code=401, detail="Invalid username or password.")

    refresh_token = _generate_refresh_token()
    token = _make_token(user)
    auth_response = _user_response(user, token, refresh_token)
    response.set_cookie(
        key="refresh_token",
        value=refresh_token,
        path="/",
        httponly=True,
        samesite="lax",
        secure=False,
    )
    auth_response["message"] = "Login successful."
    auth_response["status"] = 200
    return auth_response


@app.post("/logout")
@app.post("/auth/logout")
@app.post("/api/v1/logout")
@app.post("/api/v1/auth/logout")
async def logout_user(
    request: Request,
    response: Response,
    authorization: str | None = Header(default=None, alias="Authorization"),
) -> dict[str, Any]:
    token = None
    if authorization and authorization.lower().startswith("bearer "):
        token = authorization.split(" ", 1)[1].strip()
    if token is None:
        token = request.cookies.get("refresh_token") or request.cookies.get("access_token")
    if token in ACTIVE_TOKENS:
        ACTIVE_TOKENS.pop(token, None)

    response.delete_cookie(key="refresh_token", path="/")
    response.delete_cookie(key="access_token", path="/")

    return {
        "message": "Logout successful.",
        "status": 200,
        "success": True,
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
