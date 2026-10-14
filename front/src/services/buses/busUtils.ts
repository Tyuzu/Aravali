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
  const shell = document.createElement("div");
  shell.className = "bus-page";

  const heading = document.createElement("div");
  heading.className = "bus-skeleton bus-skeleton-title";
  shell.appendChild(heading);

  const summary = document.createElement("div");
  summary.className = "bus-summary";
  [0, 1, 2].forEach(() => {
    const stat = document.createElement("div");
    stat.className = "bus-stat bus-skeleton-box";
    shell.appendChild(stat);
  });

  const controls = document.createElement("div");
  controls.className = "bus-controls";
  const search = document.createElement("div");
  search.className = "bus-skeleton bus-skeleton-input";
  const filter = document.createElement("div");
  filter.className = "bus-skeleton bus-skeleton-input short";
  controls.append(search, filter);
  shell.appendChild(controls);

  const list = document.createElement("div");
  list.className = "bus-list";
  [0, 1, 2].forEach(() => {
    const card = document.createElement("div");
    card.className = "bus-card bus-skeleton-box";
    list.appendChild(card);
  });
  shell.appendChild(list);

  return shell;
}

export function createBusEmptyState(title: string, message: string): HTMLElement {
  const block = document.createElement("div");
  block.className = "bus-empty-state";

  const h3 = document.createElement("h3");
  h3.textContent = title;

  const p = document.createElement("p");
  p.textContent = message;

  block.append(h3, p);
  return block;
}

export function createBusErrorState(message: string, retryAction?: () => void): HTMLElement {
  const error = document.createElement("div");
  error.className = "bus-error";
  error.textContent = message;

  if (retryAction) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "bus-retry-btn";
    button.textContent = "Retry";
    button.addEventListener("click", retryAction);
    error.appendChild(button);
  }

  return error;
}

export function buildFavoriteButton(routeId: string, label: string): HTMLButtonElement {
  const button = document.createElement("button");
  button.type = "button";
  button.className = "bus-favorite-btn";
  button.textContent = readBusFavorites().includes(routeId) ? "★ Saved" : "☆ Save";
  button.title = readBusFavorites().includes(routeId) ? `Remove ${label} from saved routes` : `Save ${label}`;
  button.setAttribute("aria-pressed", String(readBusFavorites().includes(routeId)));
  button.addEventListener("click", () => {
    const favorites = toggleBusFavorite(routeId);
    button.textContent = favorites.includes(routeId) ? "★ Saved" : "☆ Save";
    button.title = favorites.includes(routeId) ? `Remove ${label} from saved routes` : `Save ${label}`;
    button.setAttribute("aria-pressed", String(favorites.includes(routeId)));
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
