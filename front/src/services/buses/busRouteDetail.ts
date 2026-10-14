// buses/busRouteDetail.ts

import { apiFetch } from "../../api/api.js";
import {
  buildFavoriteButton,
  createBusEmptyState,
  createBusErrorState,
  createBusSkeleton,
  getCurrentLocation,
  sortStopsByDistance
} from "./busUtils.js";

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
  _isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const routeId = asRouteId(params);
  const wrapper = document.createElement("div");
  wrapper.className = "bus-page";
  wrapper.appendChild(createBusSkeleton());
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const route = await apiFetch<BusRoute>(`/bustime/${encodeURIComponent(routeId)}`, "GET");
      const nearbyStops = sortStopsByDistance(route.stops ?? [], await getCurrentLocation());

      wrapper.replaceChildren();

      const header = document.createElement("div");
      header.className = "bus-detail-header";

      const title = document.createElement("h2");
      title.textContent = `${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim();
      header.appendChild(title);
      header.appendChild(buildFavoriteButton(routeId, route.name ?? routeId));

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

      const nearbyNote = document.createElement("p");
      nearbyNote.className = "bus-inline-note";
      if (nearbyStops.length) {
        const top = nearbyStops.slice(0, 3).map((stop) => stop.name ?? stop.code ?? "Stop").join(", ");
        nearbyNote.textContent = `Stops near you: ${top}`;
      } else {
        nearbyNote.textContent = "Location access is off. Turn on geolocation to see stops nearest to you.";
      }

      const stopsHeader = document.createElement("h3");
      stopsHeader.textContent = "Route Stops";
      const stopList = document.createElement("ul");
      stopList.className = "bus-stop-list";
      if (!route.stops?.length) {
        stopList.appendChild(createBusEmptyState("No stops listed", "This route is currently missing stop data."));
      } else {
        nearbyStops.forEach((stop) => {
          const item = document.createElement("li");
          item.textContent = `${stop.order ?? "-"}. ${stop.name ?? "Stop"} (Code: ${stop.code ?? "—"})${stop.arrival ? ` • Arrival: ${stop.arrival}` : ""}`;
          stopList.appendChild(item);
        });
      }

      const vehiclesHeader = document.createElement("h3");
      vehiclesHeader.textContent = "Assigned Vehicles";
      const vehicleList = document.createElement("div");
      vehicleList.className = "bus-list";
      if (!route.vehicles?.length) {
        vehicleList.appendChild(createBusEmptyState("No active vehicles", "Vehicle tracking is unavailable for this route right now."));
      } else {
        route.vehicles.forEach((vehicle) => {
          const card = document.createElement("article");
          card.className = "bus-card";
          const heading = document.createElement("h4");
          heading.textContent = `Vehicle #${vehicle.id ?? "Unknown"} (${vehicle.status ?? "Unknown status"})`;
          const info = document.createElement("p");
          info.textContent = `Capacity: ${vehicle.capacity ?? "—"}% • Coordinates: ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`;
          card.append(heading, info);
          vehicleList.appendChild(card);
        });
      }

      wrapper.append(header, metrics, details, nearbyNote, stopsHeader, stopList, vehiclesHeader, vehicleList);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unable to load route details.";
      wrapper.replaceChildren(createBusErrorState(message, load));
    }
  };

  await load();
}