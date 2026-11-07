from __future__ import annotations

from typing import Any

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from routes import route_registry, router

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


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
