const APP_SHELL = [
  "/",
  "/index.html",
  "/manifest.json",
  "/offline.html",
  "/offWork.html",
];

const CACHE_NAME = "scav-pwa-v1";

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => cache.addAll(APP_SHELL))
      .then(() => self.skipWaiting())
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) => Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key))))
      .then(() => self.clients.claim())
  );
});

self.addEventListener("fetch", (event) => {
  if (event.request.method !== "GET") return;

  const requestUrl = new URL(event.request.url);
  if (requestUrl.origin !== self.location.origin) return;

  if (event.request.mode === "navigate") {
    event.respondWith(
      fetch(event.request)
        .then((response) => {
          const responseClone = response.clone();
          caches.open(CACHE_NAME).then((cache) => cache.put(event.request, responseClone));
          return response;
        })
        .catch(() => caches.match("/offline.html") || caches.match("/"))
    );
    return;
  }

  event.respondWith(
    caches.match(event.request).then((cachedResponse) => {
      if (cachedResponse) return cachedResponse;

      return fetch(event.request)
        .then((response) => {
          const responseClone = response.clone();
          caches.open(CACHE_NAME).then((cache) => {
            if (requestUrl.pathname.startsWith("/assets/") || requestUrl.pathname.endsWith(".png") || requestUrl.pathname.endsWith(".jpg") || requestUrl.pathname.endsWith(".jpeg") || requestUrl.pathname.endsWith(".webp") || requestUrl.pathname.endsWith(".svg") || requestUrl.pathname.endsWith(".ico")) {
              cache.put(event.request, responseClone);
            }
          });
          return response;
        })
        .catch(() => caches.match("/offline.html") || caches.match("/"));
    })
  );
});
