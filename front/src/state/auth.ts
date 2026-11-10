import type { AppState, AuthState, User } from "./types";

export const AUTH_ALIAS_KEYS = new Set<string>([
  "token",
  "user",
  "roles",
  "permissions",
  "username",
  "userid",
  "isLoggedIn"
]);

export function normalizeRoles(roles: any): string[] {
  if (!Array.isArray(roles)) {
    return roles == null ? [] : [roles];
  }

  return [...new Set(
    roles.filter(
      (role) => role !== null && role !== undefined && String(role).length > 0
    )
  )];
}

export function normalizePermissions(permissions: any): string[] {
  if (!Array.isArray(permissions)) {
    return permissions == null ? [] : [permissions];
  }

  return [...new Set(
    permissions.filter(
      (permission) =>
        permission !== null &&
        permission !== undefined &&
        String(permission).length > 0
    )
  )];
}

export function normalizeAuth(
  authValue: Partial<AuthState> = {},
  previousAuth: Partial<AuthState> = {}
): AuthState {
  const source = authValue && typeof authValue === "object" ? authValue : {};
  const previous = previousAuth && typeof previousAuth === "object" ? previousAuth : {};

  const accessToken = Object.prototype.hasOwnProperty.call(source, "accessToken")
    ? source.accessToken || null
    : previous.accessToken || null;
  const user = Object.prototype.hasOwnProperty.call(source, "user")
    ? source.user || null
    : previous.user || null;
  const roles = normalizeRoles(
    Object.prototype.hasOwnProperty.call(source, "roles") ? source.roles : previous.roles
  );
  const permissions = normalizePermissions(
    Object.prototype.hasOwnProperty.call(source, "permissions")
      ? source.permissions
      : previous.permissions
  );
  const hasUserIdentity = Boolean(
    user && typeof user === "object" && (user.id || user.userid)
  );
  const explicitAuthState = Object.prototype.hasOwnProperty.call(
    source,
    "isAuthenticated"
  )
    ? Boolean(source.isAuthenticated)
    : undefined;
  const isAuthenticated =
    explicitAuthState !== undefined
      ? explicitAuthState
      : Boolean(accessToken || hasUserIdentity);

  return {
    isAuthenticated,
    loading: Object.prototype.hasOwnProperty.call(source, "loading")
      ? Boolean(source.loading)
      : Boolean(previous.loading),
    accessToken,
    user,
    roles,
    permissions
  };
}

export function syncDerivedAuthState(
  rawState: AppState,
  nextAuth: Partial<AuthState> = rawState.auth
): void {
  const user = nextAuth?.user;
  const hasUserIdentity = Boolean(
    user && typeof user === "object" && (user.id || user.userid)
  );
  rawState.isLoggedIn = Boolean(
    nextAuth?.isAuthenticated || nextAuth?.accessToken || hasUserIdentity
  );
}

export function getAuthAlias(rawState: AppState, key: string): any {
  switch (key) {
    case "token":
      return rawState.auth?.accessToken || null;
    case "user":
      return rawState.auth?.user || null;
    case "roles":
      return Array.isArray(rawState.auth?.roles) ? rawState.auth.roles : [];
    case "permissions":
      return Array.isArray(rawState.auth?.permissions) ? rawState.auth.permissions : [];
    case "isLoggedIn": {
      const user = rawState.auth?.user;
      const hasUserValidId = Boolean(
        user && typeof user === "object" && (user.id || user.userid)
      );
      return Boolean(
        rawState.isLoggedIn ||
          rawState.auth?.isAuthenticated ||
          rawState.auth?.accessToken ||
          hasUserValidId
      );
    }
    case "username": {
      const user = rawState.auth?.user;
      if (user && typeof user === "object") {
        return user?.username ?? user?.name ?? "";
      }
      return typeof user === "string" ? user : "";
    }
    case "userid": {
      const user = rawState.auth?.user;
      if (user && typeof user === "object") {
        return user?.id ?? user?.userid ?? null;
      }
      return null;
    }
    default:
      return undefined;
  }
}

export function updateAuthUserProperty(
  rawState: AppState,
  property: string,
  value: any
): void {
  const currentUser = rawState.auth?.user;
  if (currentUser && typeof currentUser === "object" && !Array.isArray(currentUser)) {
    rawState.auth.user = {
      ...currentUser,
      [property]: value
    };
  } else {
    rawState.auth.user = {
      [property]: value
    };
  }
}
