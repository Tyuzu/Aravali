// buses/busLiveTracking.ts

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

  const wrapper = document.createElement("div");
  wrapper.className = "bus-page";

  const title = document.createElement("h2");
  title.textContent = `Live tracking: ${routeId}`;

  const updated = document.createElement("p");
  const list = document.createElement("div");
  list.className = "bus-list";

  wrapper.append(title, updated, list);
  container.appendChild(wrapper);

  let timer: number | null = null;

  const fetchData = async () => {
    try {
      const data = await apiFetch<BusLiveResponse>(`/bus/live/${encodeURIComponent(routeId)}`, "GET");
      const vehicles = data.vehicles ?? [];

      title.textContent = `${data.name ?? "Bus route"} live tracking`;
      updated.textContent = data.updated_at ? `Last updated: ${new Date(data.updated_at).toLocaleTimeString()}` : "Live data synced";

      list.replaceChildren();

      if (vehicles.length === 0) {
        const emptyMsg = document.createElement("p");
        emptyMsg.textContent = "No active vehicles currently reporting location on this route.";
        list.appendChild(emptyMsg);
        return;
      }

      vehicles.forEach((vehicle) => {
        const card = document.createElement("article");
        card.className = "bus-card";

        const heading = document.createElement("h3");
        heading.textContent = `Vehicle #${vehicle.id ?? "Unknown"} (${vehicle.status ?? "Active"})`;

        const info = document.createElement("p");
        info.textContent = `GPS Coordinates: ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`;

        const capacityBar = document.createElement("div");
        capacityBar.style.marginTop = "8px";
        const capacityValue = vehicle.capacity ?? 50;
        capacityBar.innerHTML = `
          <div style="font-size: 12px; margin-bottom: 3px;">Estimated Capacity: ${capacityValue}%</div>
          <div style="background: #e2e8f0; border-radius: 4px; height: 8px; width: 100%; overflow: hidden;">
            <div style="background: ${capacityValue > 80 ? '#ef4444' : capacityValue > 50 ? '#f59e0b' : '#10b981'}; width: ${capacityValue}%; height: 100%;"></div>
          </div>
        `;

        card.append(heading, info, capacityBar);
        list.appendChild(card);
      });
    } catch (error) {
      updated.textContent = error instanceof Error ? error.message : "Error syncing live coordinates.";
    }
  };

  await fetchData();
  timer = window.setInterval(fetchData, 15000);

  // Clean up interval if container gets replaced
  const observer = new MutationObserver(() => {
    if (!document.body.contains(container)) {
      if (timer) clearInterval(timer);
      observer.disconnect();
    }
  });
  observer.observe(document.body, { childList: true, subtree: true });
}