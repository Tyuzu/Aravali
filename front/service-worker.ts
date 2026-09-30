/// <reference lib="webworker" />
declare const self: ServiceWorkerGlobalScope;

const CACHE_VERSION = "v18";
const STATIC_CACHE = `scav-static-${CACHE_VERSION}`;
const STAGE_2_CACHE = `scav-stage2-${CACHE_VERSION}`;
const CURRENT_CACHES = [STATIC_CACHE, STAGE_2_CACHE];

// Stage-1: Core shell required for instant offline booting
const STAGE_1_ASSETS = [
  "/",
  "/index.html",
  "/offline.html",
  "/offWork.html",
  "/manifest.json",
];

// Stage-2: Secondary assets fetched in the background later
const STAGE_2_ASSETS = [
  "/assets/icon-192.png",
  "/assets/icon-512.png",
  "/fonts/main-font.woff2",
];

/* =========================================================
   INSTALL EVENT: Stage 1 Core Caching
========================================================= */
self.addEventListener("install", (event: ExtendableEvent) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(STATIC_CACHE);
      // Resilient caching: won't fail the entire SW if one non-critical file errors
      await Promise.allSettled(
        STAGE_1_ASSETS.map((asset) => cache.add(asset))
      );
      await self.skipWaiting();
    })()
  );
});

/* =========================================================
   ACTIVATE EVENT: Cache Cleanup & Client Takeover
========================================================= */
self.addEventListener("activate", (event: ExtendableEvent) => {
  event.waitUntil(
    (async () => {
      // Clean up legacy cache versions
      const cacheNames = await caches.keys();
      // ✅ FIX: Clean, readable, and properly typed
      await Promise.all(
        cacheNames
          .filter((cacheName) => !CURRENT_CACHES.includes(cacheName))
          .map((cacheName) => {
            console.log(`[SW] Deleting legacy cache: ${cacheName}`);
            return caches.delete(cacheName);
          })
      );
      await self.clients.claim();
    })()
  );
});

/* =========================================================
   FETCH EVENT: Cache-First Strategy with Offline Fallback
========================================================= */
self.addEventListener("fetch", (event: FetchEvent) => {
  if (event.request.method !== "GET") return;

  event.respondWith(
    (async () => {
      // 1. Try cache match across both STATIC and STAGE_2 caches
      const cachedResponse = await caches.match(event.request);
      if (cachedResponse) {
        return cachedResponse;
      }

      // 2. Fall back to network fetch
      try {
        const networkResponse = await fetch(event.request);
        return networkResponse;
      } catch (error) {
        console.warn("[SW] Fetch failed; returning offline fallback.", error);

        // 3. Navigation fallback for HTML pages when offline
        if (event.request.mode === "navigate") {
          const offlinePage = await caches.match("/offline.html");
          if (offlinePage) return offlinePage;
        }

        throw error;
      }
    })()
  );
});

/* =========================================================
   MESSAGE RECEIVER: Stage 2 Trigger & Commands
========================================================= */
self.addEventListener("message", (event: ExtendableMessageEvent) => {
  if (event.data?.type === "WARMUP_STAGE_2") {
    event.waitUntil(executeStage2Boot());
  } else if (event.data?.type === "SKIP_WAITING") {
    self.skipWaiting();
  }
});

/* =========================================================
   STAGE-2 EXECUTION: Background Prefetching
========================================================= */
async function executeStage2Boot(): Promise<void> {
  console.log("[SW Stage-2] Starting background warm-up...");
  const cache = await caches.open(STAGE_2_CACHE);

  for (const asset of STAGE_2_ASSETS) {
    try {
      const response = await fetch(asset, { priority: "low" } as RequestInit);
      if (response.ok) {
        await cache.put(asset, response);
      }
    } catch {
      console.warn(`[SW Stage-2] Pre-fetch skipped for ${asset}`);
    }
  }
  console.log("[SW Stage-2] Warm-up complete.");
}