import { apiFetch } from "../../api/api.js";

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
  isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = "<h2>Bus map</h2><p>Sign in to access live map details.</p>";
    container.appendChild(wrapper);
    return;
  }

  try {
    const data = await apiFetch<BusMapResponse>("/bus/map", "GET");
    const stops = data.stops ?? [];
    const routes = data.routes ?? [];

    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";

    const title = document.createElement("h2");
    title.textContent = "Bus Network Map";

    const map = document.createElement("div");
    map.className = "bus-map";
    map.style.position = "relative";
    map.style.minHeight = "280px";
    map.style.border = "1px solid #d9d9d9";
    map.style.borderRadius = "12px";
    map.style.background = "linear-gradient(135deg, #f7fafc, #edf2f7)";

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

        const marker = document.createElement("div");
        marker.textContent = stop.code ?? stop.name?.slice(0, 1) ?? "•";
        marker.title = `${stop.name ?? "Stop"} (${stop.code ?? "N/A"})`;
        marker.style.position = "absolute";
        marker.style.left = `${Math.min(Math.max(point.x, 4), 96)}%`;
        marker.style.top = `${Math.min(Math.max(point.y, 4), 96)}%`;
        marker.style.transform = "translate(-50%, -50%)";
        marker.style.width = "24px";
        marker.style.height = "24px";
        marker.style.borderRadius = "50%";
        marker.style.background = "#1d4ed8";
        marker.style.color = "#fff";
        marker.style.display = "grid";
        marker.style.placeItems = "center";
        marker.style.fontSize = "10px";
        marker.style.fontWeight = "700";
        marker.style.cursor = "pointer";
        marker.style.boxShadow = "0 2px 8px rgba(29, 78, 216, 0.35)";
        map.appendChild(marker);
      });
    }

    const info = document.createElement("p");
    info.textContent = `Map Center Coordinate: ${data.center_lat ?? "—"}, ${data.center_lng ?? "—"} • Zoom Level: ${data.zoom ?? "—"} • Total Stops: ${stops.length}`;

    const routeList = document.createElement("div");
    routeList.className = "bus-list";

    routes.forEach((route) => {
      const card = document.createElement("article");
      card.className = "bus-card";
      const heading = document.createElement("h3");
      heading.textContent = route.name ?? "Route";
      const details = document.createElement("p");
      details.textContent = `${route.alias ?? ""}${route.status ? ` • Status: ${route.status}` : ""}`.trim();
      card.append(heading, details);
      routeList.appendChild(card);
    });

    const stopList = document.createElement("ul");
    stopList.className = "bus-stop-list";
    stops.forEach((stop) => {
      const item = document.createElement("li");
      item.textContent = `${stop.name ?? "Stop"} (Code: ${stop.code ?? "—"})`;
      stopList.appendChild(item);
    });

    wrapper.append(title, info, map, routeList, stopList);
    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus map</h2><p>${error instanceof Error ? error.message : "Unable to load map data."}</p>`;
    container.appendChild(wrapper);
  }
}