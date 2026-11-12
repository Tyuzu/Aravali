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

type BeforeInstallPromptEvent = Event & {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed"; platform: string }>;
};

/* =========================================================
   STATE & MAIN EXPORT
========================================================= */
let serviceWorkerRefreshing = false;
let deferredInstallPrompt: BeforeInstallPromptEvent | null = null;

export function setupServiceWorker(): void {
  if (!("serviceWorker" in navigator)) {
    console.warn("[SW] Service workers are not supported.");
    return;
  }

  let hadController = Boolean(navigator.serviceWorker.controller);

  window.addEventListener("beforeinstallprompt", (event: Event) => {
    event.preventDefault();
    deferredInstallPrompt = event as BeforeInstallPromptEvent;
    window.dispatchEvent(new CustomEvent("pwa:install-ready"));
  });

  window.addEventListener("appinstalled", () => {
    deferredInstallPrompt = null;
    window.dispatchEvent(new CustomEvent("pwa:installed"));
  });

  window.addEventListener("load", async () => {
    try {
      const registration = await navigator.serviceWorker.register("/service-worker.js", {
        updateViaCache: "none"
      });
      console.log("[SW] Registered:", registration.scope);

      try {
        await registration.update();
      } catch (error: unknown) {
        console.warn("[SW] Update check failed:", error);
      }

      if (registration.waiting) {
        requestServiceWorkerActivation(registration.waiting);
      }

      registration.addEventListener("updatefound", () => {
        const newWorker = registration.installing;
        if (!newWorker) return;

        newWorker.addEventListener("statechange", () => {
          if (newWorker.state === "installed" && navigator.serviceWorker.controller) {
            console.log("[SW] New version available.");
            requestServiceWorkerActivation(newWorker);
          }
        });
      });
    } catch (error: unknown) {
      console.error("[SW] Registration failed:", error);
    }
  });

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
export async function triggerInstallPrompt(): Promise<boolean> {
  if (!deferredInstallPrompt) return false;

  try {
    await deferredInstallPrompt.prompt();
    const choice = await deferredInstallPrompt.userChoice;
    deferredInstallPrompt = null;
    return choice.outcome === "accepted";
  } catch {
    deferredInstallPrompt = null;
    return false;
  }
}

function requestServiceWorkerActivation(worker: ServiceWorker | null): void {
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