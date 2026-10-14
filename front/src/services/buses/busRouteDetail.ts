import { apiFetch } from "../../api/api.js";

type BusStop = {
  id?: string;
  name?: string;
  code?: string;
  order?: number;
  arrival?: string;
  lat?: number;
  lng?: number;
};

type BusVehicle = {
  id?: string;
  status?: string;
  capacity?: number;
  lat?: number;
  lng?: number;
};

type BusRoute = {
  route_id?: string;
  id?: string;
  name?: string;
  alias?: string;
  depot_id?: string;
  depot_name?: string;
  status?: string;
  next_departure?: string;
  headway_minutes?: number;
  distance_minutes?: number;
  stops?: BusStop[];
  vehicles?: BusVehicle[];
};

const asRouteId = (params?: Record<string, string | undefined> | string): string => {
  if (typeof params === "string") {
    return params;
  }

  if (!params) {
    return "unknown";
  }

  return params.busid ?? params.id ?? "unknown";
};

export async function displayRouteDetails(
  isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const routeId = asRouteId(params);
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus route ${routeId}</h2><p>Sign in to view route details.</p>`;
    container.appendChild(wrapper);
    return;
  }

  try {
    const route = await apiFetch<BusRoute>(`/bustime/${encodeURIComponent(routeId)}`, "GET");

    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";

    const title = document.createElement("h2");
    title.textContent = `${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim();

    const metrics = document.createElement("div");
    metrics.className = "bus-summary";

    const stat1 = document.createElement("div");
    stat1.className = "bus-stat";
    stat1.innerHTML = `<strong>${route.status ?? "—"}</strong><span>Status</span>`;

    const stat2 = document.createElement("div");
    stat2.className = "bus-stat";
    stat2.innerHTML = `<strong>${route.next_departure ?? "—"}</strong><span>Next departure</span>`;

    const stat3 = document.createElement("div");
    stat3.className = "bus-stat";
    stat3.innerHTML = `<strong>${route.headway_minutes ?? "—"}</strong><span>Headway mins</span>`;

    metrics.append(stat1, stat2, stat3);

    const details = document.createElement("p");
    details.textContent = `Depot: ${route.depot_name ?? route.depot_id ?? "Unknown"} • Estimated Total Time: ${route.distance_minutes ?? "—"} mins`;

    const stopsHeader = document.createElement("h3");
    stopsHeader.textContent = "Route Stops";
    const stopList = document.createElement("ul");
    stopList.className = "bus-stop-list";
    (route.stops ?? []).forEach((stop) => {
      const item = document.createElement("li");
      item.textContent = `${stop.order ?? "-"}. ${stop.name ?? "Stop"} (Code: ${stop.code ?? "—"})${stop.arrival ? ` • Arrival: ${stop.arrival}` : ""}`;
      stopList.appendChild(item);
    });

    const vehiclesHeader = document.createElement("h3");
    vehiclesHeader.textContent = "Assigned Vehicles";
    const vehicleList = document.createElement("div");
    vehicleList.className = "bus-list";
    (route.vehicles ?? []).forEach((vehicle) => {
      const card = document.createElement("article");
      card.className = "bus-card";
      const heading = document.createElement("h4");
      heading.textContent = `Vehicle #${vehicle.id ?? "Unknown"} (${vehicle.status ?? "Unknown status"})`;
      const info = document.createElement("p");
      info.textContent = `Capacity: ${vehicle.capacity ?? "—"}% • Coordinates: ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`;
      card.append(heading, info);
      vehicleList.appendChild(card);
    });

    wrapper.append(title, metrics, details, stopsHeader, stopList, vehiclesHeader, vehicleList);
    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus route ${routeId}</h2><p>${error instanceof Error ? error.message : "Unable to load route details."}</p>`;
    container.appendChild(wrapper);
  }
}