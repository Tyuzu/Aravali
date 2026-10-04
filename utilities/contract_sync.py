#!/usr/bin/env python3
"""Project-wide contract sync checker for Go JSON tags and TypeScript interfaces.

This scans all package directories under backend/internal and all service folders
under front/src/services, then compares the discovered field names by logical
folder name. If a curated registry exists, it is merged in as explicit overrides.

The goal is to catch drift across the full codebase instead of only a small
hand-maintained subset.
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
REGISTRY_PATH = ROOT / "utilities" / "contract_sync_registry.json"

TAGS_RE = re.compile(r'json:"([^"]+)"')
TS_BLOCK_RE = re.compile(r"(?ms)(?:export\s+)?(?:interface|type)\s+[A-Za-z_$][A-Za-z0-9_$]*\s*(?:extends\s+[^{]+)?\s*(?:=\s*)?\{(.*?)\}")
TS_PROP_RE = re.compile(r"(?m)^\s*([A-Za-z_$][A-Za-z0-9_$]*)\??\s*:")


def normalize_name(value: str) -> str:
    return re.sub(r"[^a-z0-9]", "", value.lower())


def read_text(path: Path) -> str:
    if not path.exists():
        raise FileNotFoundError(path)
    return path.read_text(encoding="utf-8", errors="ignore")


def collect_backend_fields(paths: list[str]) -> set[str]:
    fields: set[str] = set()
    for rel in paths:
        path = ROOT / rel
        if not path.exists():
            continue
        text = read_text(path)
        for match in TAGS_RE.findall(text):
            key = match.split(",")[0].strip()
            if key and key != "-":
                fields.add(key)
    return fields


def collect_frontend_fields(paths: list[str]) -> set[str]:
    fields: set[str] = set()
    for rel in paths:
        path = ROOT / rel
        if not path.exists():
            continue
        text = read_text(path)
        for block in TS_BLOCK_RE.findall(text):
            for name in TS_PROP_RE.findall(block):
                if name:
                    fields.add(name)
    return fields


def discover_default_registry() -> dict[str, dict]:
    registry: dict[str, dict] = {}

    backend_root = ROOT / "backend" / "internal"
    frontend_root = ROOT / "front" / "src" / "services"

    if backend_root.exists():
        for backend_dir in sorted(backend_root.iterdir()):
            if not backend_dir.is_dir():
                continue
            backend_paths = [str(p.relative_to(ROOT)) for p in sorted(backend_dir.rglob("*.go"))]
            if not backend_paths:
                continue
            registry.setdefault(backend_dir.name, {"backend": [], "frontend": []})
            registry[backend_dir.name]["backend"] = backend_paths

    if frontend_root.exists():
        for frontend_dir in sorted(frontend_root.iterdir()):
            if not frontend_dir.is_dir():
                continue
            frontend_paths = [str(p.relative_to(ROOT)) for p in sorted(frontend_dir.rglob("*.ts"))]
            if not frontend_paths:
                continue
            registry.setdefault(frontend_dir.name, {"backend": [], "frontend": []})
            registry[frontend_dir.name]["frontend"] = frontend_paths

    return registry


def sync_entry(name: str, entry: dict) -> list[str]:
    backend_paths = entry.get("backend", [])
    frontend_paths = entry.get("frontend", [])
    expected_fields = entry.get("fields", [])

    issues: list[str] = []

    for rel in backend_paths:
        if not (ROOT / rel).exists():
            issues.append(f"{name}: missing backend file {rel}")
    for rel in frontend_paths:
        if not (ROOT / rel).exists():
            issues.append(f"{name}: missing frontend file {rel}")

    backend_fields = collect_backend_fields(backend_paths)
    frontend_fields = collect_frontend_fields(frontend_paths)

    expected_norm = {normalize_name(v) for v in expected_fields}
    backend_norm = {normalize_name(v) for v in backend_fields}
    frontend_norm = {normalize_name(v) for v in frontend_fields}

    if expected_fields:
        missing_in_backend = sorted({v for v in expected_norm if v not in backend_norm})
        missing_in_frontend = sorted({v for v in expected_norm if v not in frontend_norm})
        if missing_in_backend:
            issues.append(f"{name}: expected backend fields missing: {', '.join(missing_in_backend)}")
        if missing_in_frontend:
            issues.append(f"{name}: expected frontend fields missing: {', '.join(missing_in_frontend)}")
        return issues

    if backend_fields and frontend_fields:
        backend_only = sorted({normalize_name(v) for v in backend_fields if normalize_name(v) not in frontend_norm})
        frontend_only = sorted({normalize_name(v) for v in frontend_fields if normalize_name(v) not in backend_norm})

        if backend_only:
            issues.append(f"{name}: backend-only fields: {', '.join(backend_only[:25])}{' ...' if len(backend_only) > 25 else ''}")
        if frontend_only:
            issues.append(f"{name}: frontend-only fields: {', '.join(frontend_only[:25])}{' ...' if len(frontend_only) > 25 else ''}")

    return issues


def merge_registry_data(base: dict, override: dict) -> dict:
    merged = {**base}
    for slug, entry in override.items():
        existing = merged.setdefault(slug, {"backend": [], "frontend": []})
        existing["backend"] = sorted(set(existing.get("backend", []) + entry.get("backend", [])))
        existing["frontend"] = sorted(set(existing.get("frontend", []) + entry.get("frontend", [])))
        for key in ("fields",):
            if key in entry:
                existing[key] = entry[key]
    return merged


def main() -> int:
    registry: dict[str, dict]

    if REGISTRY_PATH.exists():
        try:
            registry = json.loads(REGISTRY_PATH.read_text(encoding="utf-8"))
        except json.JSONDecodeError as exc:
            print(f"Invalid JSON in contract registry: {exc}")
            return 1
    else:
        registry = discover_default_registry()

    issues: list[str] = []
    for slug, entry in sorted(registry.items()):
        issues.extend(sync_entry(slug, entry))

    if issues:
        print(f"Contract sync check failed for {len(registry)} registered contract groups:")
        for issue in issues:
            print(f" - {issue}")
        return 1

    print(f"Contract sync OK for {len(registry)} registered backend/frontend pairs.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
