import { createElement } from "../../components/createElement.js";

type BusCoordinate = {
  latitude: number;
  longitude: number;
};

type BusStopLike = {
  id?: string;
  name?: string;
  code?: string;
  order?: number;
  arrival?: string;
  lat?: number;
  lng?: number;
};

export const BUS_FAVORITES_KEY = "bus:favorites";

export function readBusFavorites(): string[] {
  try {
    const raw = localStorage.getItem(BUS_FAVORITES_KEY);
    if (!raw) {
      return [];
    }
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed.filter((item): item is string => typeof item === "string") : [];
  } catch {
    return [];
  }
}

export function toggleBusFavorite(routeId: string): string[] {
  const next = readBusFavorites();
  const exists = next.includes(routeId);
  const updated = exists ? next.filter((id) => id !== routeId) : [...next, routeId];
  localStorage.setItem(BUS_FAVORITES_KEY, JSON.stringify(updated));
  return updated;
}

export function createBusSkeleton(): HTMLElement {
  return createElement("div", { class: "bus-page" }, [
    createElement("div", { class: "bus-skeleton bus-skeleton-title" }),
    createElement("div", { class: "bus-summary" }, [
      createElement("div", { class: "bus-stat bus-skeleton-box" }),
      createElement("div", { class: "bus-stat bus-skeleton-box" }),
      createElement("div", { class: "bus-stat bus-skeleton-box" })
    ]),
    createElement("div", { class: "bus-controls" }, [
      createElement("div", { class: "bus-skeleton bus-skeleton-input" }),
      createElement("div", { class: "bus-skeleton bus-skeleton-input short" })
    ]),
    createElement("div", { class: "bus-list" }, [
      createElement("div", { class: "bus-card bus-skeleton-box" }),
      createElement("div", { class: "bus-card bus-skeleton-box" }),
      createElement("div", { class: "bus-card bus-skeleton-box" })
    ])
  ]) as HTMLElement;
}

export function createBusEmptyState(title: string, message: string): HTMLElement {
  return createElement("div", { class: "bus-empty-state" }, [
    createElement("h3", {}, [title]),
    createElement("p", {}, [message])
  ]) as HTMLElement;
}

export function createBusErrorState(message: string, retryAction?: () => void): HTMLElement {
  const error = createElement("div", { class: "bus-error" }, [message]) as HTMLDivElement;

  if (retryAction) {
    const button = createElement("button", {
      type: "button",
      class: "bus-retry-btn"
    }, ["Retry"]) as HTMLButtonElement;
    button.addEventListener("click", retryAction);
    error.appendChild(button);
  }

  return error;
}

export function buildFavoriteButton(routeId: string, label: string): HTMLButtonElement {
  const isSaved = readBusFavorites().includes(routeId);
  const button = createElement("button", {
    type: "button",
    class: "bus-favorite-btn",
    title: isSaved ? `Remove ${label} from saved routes` : `Save ${label}`,
    "aria-pressed": String(isSaved)
  }, [isSaved ? "★ Saved" : "☆ Save"]) as HTMLButtonElement;

  button.addEventListener("click", () => {
    const favorites = toggleBusFavorite(routeId);
    const saved = favorites.includes(routeId);
    button.textContent = saved ? "★ Saved" : "☆ Save";
    button.title = saved ? `Remove ${label} from saved routes` : `Save ${label}`;
    button.setAttribute("aria-pressed", String(saved));
  });

  return button;
}

export function getCurrentLocation(): Promise<BusCoordinate | null> {
  if (!navigator.geolocation) {
    return Promise.resolve(null);
  }

  return new Promise((resolve) => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        resolve({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude
        });
      },
      () => resolve(null),
      {
        enableHighAccuracy: true,
        timeout: 8000,
        maximumAge: 60000
      }
    );
  });
}

export function sortStopsByDistance(stops: BusStopLike[], location: BusCoordinate | null): BusStopLike[] {
  if (!location || !stops.length) {
    return stops;
  }

  return [...stops].sort((a, b) => {
    const distA = a.lat != null && a.lng != null ? distanceKm(a.lat, a.lng, location.latitude, location.longitude) : Number.POSITIVE_INFINITY;
    const distB = b.lat != null && b.lng != null ? distanceKm(b.lat, b.lng, location.latitude, location.longitude) : Number.POSITIVE_INFINITY;
    return distA - distB;
  });
}

function distanceKm(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const toRad = (value: number): number => (value * Math.PI) / 180;
  const earthRadiusKm = 6371;
  const dLat = toRad(lat2 - lat1);
  const dLon = toRad(lon2 - lon1);

  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) * Math.sin(dLon / 2);

  return earthRadiusKm * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}
