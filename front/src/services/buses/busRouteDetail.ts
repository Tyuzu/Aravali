import { createElement } from "../../components/createElement.js";
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
  const wrapper = createElement("div", { class: "bus-page" }, [createBusSkeleton()]) as HTMLDivElement;
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const route = await apiFetch<BusRoute>(`/bustime/${encodeURIComponent(routeId)}`, "GET");
      const nearbyStops = sortStopsByDistance(route.stops ?? [], await getCurrentLocation());

      wrapper.replaceChildren();

      const header = createElement("div", { class: "bus-detail-header" }, [
        createElement("h2", {}, [`${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim()]),
        buildFavoriteButton(routeId, route.name ?? routeId)
      ]);

      const metrics = createElement("div", { class: "bus-summary" }, [
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [route.status ?? "—"]),
          createElement("span", {}, ["Status"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [route.next_departure ?? "—"]),
          createElement("span", {}, ["Next departure"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(route.headway_minutes ?? "—")]),
          createElement("span", {}, ["Headway mins"])
        ])
      ]);

      const details = createElement("p", {}, [`Depot: ${route.depot_name ?? route.depot_id ?? "Unknown"} • Estimated Total Time: ${route.distance_minutes ?? "—"} mins`]);

      const nearbyNote = createElement("p", { class: "bus-inline-note" }, [
        nearbyStops.length
          ? `Stops near you: ${nearbyStops.slice(0, 3).map((stop) => stop.name ?? stop.code ?? "Stop").join(", ")}`
          : "Location access is off. Turn on geolocation to see stops nearest to you."
      ]);

      const stopsHeader = createElement("h3", {}, ["Route Stops"]);
      const stopList = createElement("ul", { class: "bus-stop-list" });
      if (!route.stops?.length) {
        stopList.appendChild(createBusEmptyState("No stops listed", "This route is currently missing stop data."));
      } else {
        nearbyStops.forEach((stop) => {
          stopList.appendChild(createElement("li", {}, [`${stop.order ?? "-"}. ${stop.name ?? "Stop"} (Code: ${stop.code ?? "—"})${stop.arrival ? ` • Arrival: ${stop.arrival}` : ""}`]));
        });
      }

      const vehiclesHeader = createElement("h3", {}, ["Assigned Vehicles"]);
      const vehicleList = createElement("div", { class: "bus-list" });
      if (!route.vehicles?.length) {
        vehicleList.appendChild(createBusEmptyState("No active vehicles", "Vehicle tracking is unavailable for this route right now."));
      } else {
        route.vehicles.forEach((vehicle) => {
          vehicleList.appendChild(createElement("article", { class: "bus-card" }, [
            createElement("h4", {}, [`Vehicle #${vehicle.id ?? "Unknown"} (${vehicle.status ?? "Unknown status"})`]),
            createElement("p", {}, [`Capacity: ${vehicle.capacity ?? "—"}% • Coordinates: ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`])
          ]));
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