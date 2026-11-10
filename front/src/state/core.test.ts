import assert from "node:assert/strict";
import test from "node:test";

import { normalizeAuth, normalizeRoles } from "./auth.ts";

test("normalizeAuth preserves authenticated state and unique roles", () => {
  const auth = normalizeAuth({
    accessToken: "token-123",
    roles: ["admin", "admin", "seller"],
    user: { id: 5, username: "demo" }
  });

  assert.equal(auth.isAuthenticated, true);
  assert.deepEqual(auth.roles, ["admin", "seller"]);
  assert.equal(auth.user?.username, "demo");
});

test("normalizeRoles filters empty and duplicate values", () => {
  assert.deepEqual(normalizeRoles(["a", "", null, "a", undefined]), ["a"]);
  assert.deepEqual(normalizeRoles(null), []);
});
