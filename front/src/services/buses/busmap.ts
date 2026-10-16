import { createElement } from "../../components/createElement.js";
import { apiFetch } from "../../api/api.js";
import {
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
  lat?: number;
  lng?: number;
  order?: number;
};

type BusRoute = {
  route_id?: string;
  id?: string;
  name?: string;
  alias?: string;
  status?: string;
  stops?: BusStop[];
};

type BusMapResponse = {
  center_lat?: number;
  center_lng?: number;
  zoom?: number;
  stops?: BusStop[];
  routes?: BusRoute[];
};

const pointFor = (lat: number, lng: number, minLat: number, maxLat: number, minLng: number, maxLng: number) => ({
  x: maxLng === minLng ? 50 : ((lng - minLng) / (maxLng - minLng)) * 100,
  y: maxLat === minLat ? 50 : ((maxLat - lat) / (maxLat - minLat)) * 100
});

export async function displayBusMap(
  _isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  const wrapper = createElement("div", { class: "bus-page" }, [createBusSkeleton()]) as HTMLDivElement;
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const data = await apiFetch<BusMapResponse>("/bus/map", "GET");
      const stops = data.stops ?? [];
      const routes = data.routes ?? [];
      const location = await getCurrentLocation();
      const nearbyStops = sortStopsByDistance(stops, location);

      wrapper.replaceChildren();

      const title = createElement("h2", {}, ["Bus Network Map"]);
      const map = createElement("div", {
        class: "bus-map",
        style: {
          position: "relative",
          minHeight: "280px",
          border: "1px solid #d9d9d9",
          borderRadius: "12px",
          background: "linear-gradient(135deg, #f7fafc, #edf2f7)"
        }
      }) as HTMLDivElement;

      if (stops.length > 0) {
        const lats = stops.map((stop) => stop.lat ?? 0);
        const lngs = stops.map((stop) => stop.lng ?? 0);
        const minLat = Math.min(...lats);
        const maxLat = Math.max(...lats);
        const minLng = Math.min(...lngs);
        const maxLng = Math.max(...lngs);

        stops.forEach((stop) => {
          const point = pointFor(
            stop.lat ?? 0,
            stop.lng ?? 0,
            minLat,
            maxLat,
            minLng,
            maxLng
          );

          const marker = createElement("div", {
            title: `${stop.name ?? "Stop"} (${stop.code ?? "N/A"})`,
            style: {
              position: "absolute",
              left: `${Math.min(Math.max(point.x, 4), 96)}%`,
              top: `${Math.min(Math.max(point.y, 4), 96)}%`,
              transform: "translate(-50%, -50%)",
              width: "24px",
              height: "24px",
              borderRadius: "50%",
              background: "#1d4ed8",
              color: "#fff",
              display: "grid",
              placeItems: "center",
              fontSize: "10px",
              fontWeight: "700",
              cursor: "pointer",
              boxShadow: "0 2px 8px rgba(29, 78, 216, 0.35)"
            }
          }, [stop.code ?? stop.name?.slice(0, 1) ?? "•"]) as HTMLDivElement;
          map.appendChild(marker);
        });
      }

      const info = createElement("p", {}, [`Map Center Coordinate: ${data.center_lat ?? "—"}, ${data.center_lng ?? "—"} • Zoom Level: ${data.zoom ?? "—"} • Total Stops: ${stops.length}`]);

      const locationNote = createElement("p", { class: "bus-inline-note" }, [
        location
          ? nearbyStops.slice(0, 5).map((stop) => stop.name ?? stop.code ?? "Stop").join(", ")
            ? `Closest stops near you: ${nearbyStops.slice(0, 5).map((stop) => stop.name ?? stop.code ?? "Stop").join(", ")}`
            : "Location found. Nearby stops are ready to review."
          : "Enable location permission to sort stops closest to you."
      ]);

      const routeList = createElement("div", { class: "bus-list" }, routes.map((route) =>
        createElement("article", { class: "bus-card" }, [
          createElement("h3", {}, [route.name ?? "Route"]),
          createElement("p", {}, [`${route.alias ?? ""}${route.status ? ` • Status: ${route.status}` : ""}`.trim()])
        ])
      ));

      const stopList = createElement("ul", { class: "bus-stop-list" });
      const displayStops = nearbyStops.length ? nearbyStops : stops;
      if (!displayStops.length) {
        stopList.appendChild(createBusEmptyState("No stops available", "There are no stop markers for the network right now."));
      } else {
        displayStops.forEach((stop) => {
          const distanceText = location && stop.lat != null && stop.lng != null
            ? ` • ${Math.round((Math.abs(stop.lat - location.latitude) + Math.abs(stop.lng - location.longitude)) * 111)} km away`
            : "";
          stopList.appendChild(createElement("li", {}, [`${stop.name ?? "Stop"} (Code: ${stop.code ?? "—"})${distanceText}`]));
        });
      }

      wrapper.append(title, info, locationNote, map, routeList, stopList);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unable to load map data.";
      wrapper.replaceChildren(createBusErrorState(message, load));
    }
  };

  await load();
}
