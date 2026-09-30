/* =========================================================
   TYPES & INTERFACES
========================================================= */
export type ServiceWorkerIncomingMessage =
  | { type: "SW_ACTIVATED"; version?: string }
  | { type: "SW_VERSION"; version?: string }
  | { type: string; [key: string]: unknown };

export interface SkipWaitingMessage {
  type: "SKIP_WAITING";
}

/* =========================================================
   STATE & MAIN EXPORT
========================================================= */
let serviceWorkerRefreshing = false;

export function setupServiceWorker(
  onUpdateFound?: (registration: ServiceWorkerRegistration) => void
): void {
  if (!("serviceWorker" in navigator)) {
    console.warn("[SW] Service workers are not supported.");
    return;
  }

  let hadController = Boolean(navigator.serviceWorker.controller);

  const register = async () => {
    try {
      const registration = await navigator.serviceWorker.register("/service-worker.js", {
        updateViaCache: "none",
      });
      console.log("[SW] Registered:", registration.scope);

      // Handle a worker already in waiting state (e.g. installed from previous tab)
      if (registration.waiting) {
        handleWaitingWorker(registration, onUpdateFound);
      }

      // Listen for newly installed workers
      registration.addEventListener("updatefound", () => {
        const newWorker = registration.installing;
        if (!newWorker) return;

        newWorker.addEventListener("statechange", () => {
          if (newWorker.state === "installed" && navigator.serviceWorker.controller) {
            console.log("[SW] New version available.");
            handleWaitingWorker(registration, onUpdateFound);
          }
        });
      });

      // Check for SW updates
      try {
        await registration.update();
      } catch (error: unknown) {
        console.warn("[SW] Update check failed:", error);
      }
    } catch (error: unknown) {
      console.error("[SW] Registration failed:", error);
    }
  };

  // Safe load guard for SPAs and bundled code
  if (document.readyState === "complete") {
    register();
  } else {
    window.addEventListener("load", register, { once: true });
  }

  // Handle page reload once new SW takes control
  navigator.serviceWorker.addEventListener("controllerchange", () => {
    if (!hadController) {
      hadController = true;
      return;
    }
    if (serviceWorkerRefreshing) return;
    serviceWorkerRefreshing = true;
    console.log("[SW] Controller changed. Reloading application.");
    window.location.reload();
  });

  navigator.serviceWorker.addEventListener("message", handleServiceWorkerMessage);
}

/* =========================================================
   HELPER FUNCTIONS
========================================================= */
function handleWaitingWorker(
  registration: ServiceWorkerRegistration,
  onUpdateFound?: (registration: ServiceWorkerRegistration) => void
): void {
  if (onUpdateFound) {
    // Notify application UI to prompt user
    onUpdateFound(registration);
  } else {
    // Fallback: Activate immediately if no callback provided
    activateServiceWorker(registration.waiting);
  }
}

export function activateServiceWorker(worker: ServiceWorker | null): void {
  if (!worker) return;
  const message: SkipWaitingMessage = { type: "SKIP_WAITING" };
  worker.postMessage(message);
}

function handleServiceWorkerMessage(event: MessageEvent<ServiceWorkerIncomingMessage>): void {
  const data = event.data;
  if (!data?.type) return;

  switch (data.type) {
    case "SW_ACTIVATED":
      console.log("[SW] Active version:", data.version);
      break;
    case "SW_VERSION":
      console.log("[SW] Version:", data.version);
      break;
  }
}