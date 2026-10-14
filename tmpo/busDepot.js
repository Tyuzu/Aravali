import { apiFetch } from "../../api/api.js";

type BusDepot = {
  depot_id?: string;
  id?: string;
  name?: string;
  location?: string;
  phone?: string;
  vehicles?: number;
  latitude?: number;
  longitude?: number;
  open_hours?: string;
  description?: string;
};

const asDepotId = (params?: Record<string, string | undefined> | string): string => {
  if (typeof params === "string") {
    return params;
  }

  if (!params) {
    return "unknown";
  }

  return params.depotid ?? params.id ?? "unknown";
};

export async function displayBusDepot(
  isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const depotId = asDepotId(params);
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus depot ${depotId}</h2><p>Sign in to view depot info.</p>`;
    container.appendChild(wrapper);
    return;
  }

  try {
    const depot = await apiFetch<BusDepot>(`/bus/depot/${encodeURIComponent(depotId)}`, "GET");

    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";

    const title = document.createElement("h2");
    title.textContent = depot.name ?? "Bus depot";

    const summary = document.createElement("div");
    summary.className = "bus-summary";

    const stat1 = document.createElement("div");
    stat1.className = "bus-stat";
    stat1.innerHTML = `<strong>${depot.vehicles ?? "—"}</strong><span>Vehicles</span>`;

    const stat2 = document.createElement("div");
    stat2.className = "bus-stat";
    stat2.innerHTML = `<strong>${depot.open_hours ?? "—"}</strong><span>Open hours</span>`;

    const stat3 = document.createElement("div");
    stat3.className = "bus-stat";
    stat3.innerHTML = `<strong>${depot.location ?? "—"}</strong><span>Location</span>`;

    summary.append(stat1, stat2, stat3);

    const details = document.createElement("ul");
    details.className = "bus-stop-list";
    const items = [
      `Phone: ${depot.phone ?? "—"}`,
      `Coordinate: ${depot.latitude ?? "—"}, ${depot.longitude ?? "—"}`,
      `Description: ${depot.description ?? "—"}`
    ];

    items.forEach((line) => {
      const item = document.createElement("li");
      item.textContent = line;
      details.appendChild(item);
    });

    wrapper.append(title, summary, details);
    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus depot ${depotId}</h2><p>${error instanceof Error ? error.message : "Unable to load depot data."}</p>`;
    container.appendChild(wrapper);
  }
}