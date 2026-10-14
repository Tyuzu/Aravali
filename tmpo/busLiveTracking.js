import { apiFetch } from "../../api/api.js";

type BusVehicle = {
  id?: string;
  lat?: number;
  lng?: number;
  status?: string;
  capacity?: number;
};

type BusLiveResponse = {
  route_id?: string;
  name?: string;
  vehicles?: BusVehicle[];
  updated_at?: string;
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

export async function displayBusLiveTracking(
  isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const routeId = asRouteId(params);
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Live tracking: ${routeId}</h2><p>Sign in to view live tracking.</p>`;
    container.appendChild(wrapper);
    return;
  }

  try {
    const data = await apiFetch<BusLiveResponse>(`/bus/live/${encodeURIComponent(routeId)}`, "GET");
    const vehicles = data.vehicles ?? [];

    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";

    const title = document.createElement("h2");
    title.textContent = `${data.name ?? "Bus route"} live tracking`;

    const updated = document.createElement("p");
    updated.textContent = data.updated_at ? `Updated ${new Date(data.updated_at).toLocaleString()}` : "No timestamp available";

    const list = document.createElement("div");
    list.className = "bus-list";

    vehicles.forEach((vehicle) => {
      const card = document.createElement("article");
      card.className = "bus-card";

      const heading = document.createElement("h3");
      heading.textContent = vehicle.id ?? "Vehicle";

      const info = document.createElement("p");
      info.textContent = `${vehicle.status ?? "Unknown status"} • Capacity ${vehicle.capacity ?? "—"} • ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`;

      card.append(heading, info);
      list.appendChild(card);
    });

    wrapper.append(title, updated, list);
    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Live tracking: ${routeId}</h2><p>${error instanceof Error ? error.message : "Unable to load live tracking."}</p>`;
    container.appendChild(wrapper);
  }
}