export function normalizeNavPath(value: string | null | undefined): string {
  if (typeof value !== "string") return "/";

  let candidate = value.trim();
  if (!candidate) return "/";

  candidate = candidate.startsWith("#/") ? candidate.slice(1) : candidate;
  candidate = candidate.split("#")[0];

  const queryIndex = candidate.indexOf("?");
  if (queryIndex >= 0) {
    candidate = candidate.slice(0, queryIndex);
  }

  if (!candidate.startsWith("/")) {
    candidate = `/${candidate}`;
  }

  if (candidate.length > 1 && candidate.endsWith("/")) {
    candidate = candidate.slice(0, -1);
  }

  return candidate || "/";
}

export function normalizeStoredRedirect(value: string | null): string | null {
  const candidate = normalizeNavPath(value);
  const ignoredPaths = ["/", "/login", "/logout", "/404", "/error/403", "/error/404"];
  return ignoredPaths.includes(candidate) ? null : candidate;
}
