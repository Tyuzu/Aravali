import test from "node:test";
import assert from "node:assert/strict";

import { normalizeNavPath, normalizeStoredRedirect } from "./navPath.js";

test("normalizeNavPath strips query strings and hashes for active nav matching", () => {
  assert.equal(normalizeNavPath("/farms?tab=all#section"), "/farms");
  assert.equal(normalizeNavPath("#/farms?tab=all"), "/farms");
  assert.equal(normalizeNavPath("/products/"), "/products");
});

test("normalizeStoredRedirect ignores query strings from a stored login redirect", () => {
  assert.equal(normalizeStoredRedirect("/farms?tab=all#details"), "/farms");
  assert.equal(normalizeStoredRedirect("/login"), null);
  assert.equal(normalizeStoredRedirect("/"), null);
});
