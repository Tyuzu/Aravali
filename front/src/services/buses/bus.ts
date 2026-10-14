// buses/bus.ts

import { apiFetch } from "../../api/api.js";

type BusRouteSummary = {
  route_id?: string;
  id?: string;
  name?: string;
  alias?: string;
  depot_name?: string;
  status?: string;
  next_departure?: string;
  headway_minutes?: number;
  vehicles?: Array<{ id?: string; status?: string }>;
};

type BusOverviewResponse = {
  routes?: BusRouteSummary[];
  total_routes?: number;
  generated_at?: string;
};

const asRouteId = (route: BusRouteSummary): string =>
  route.route_id ?? route.id ?? "unknown";

const formatTimestamp = (value?: string): string => {
  if (!value) {
    return "—";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
};

export async function displayBus(
  isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = "<h2>Bus</h2><p>Sign in to access the full bus overview.</p>";
    container.appendChild(wrapper);
    return;
  }

  const wrapper = document.createElement("div");
  wrapper.className = "bus-page";

  try {
    const data = await apiFetch<BusOverviewResponse>("/bustime", "GET");
    const routes = data.routes ?? [];

    const heading = document.createElement("h2");
    heading.textContent = "Bus overview";

    const summary = document.createElement("div");
    summary.className = "bus-summary";

    const total = document.createElement("div");
    total.className = "bus-stat";
    total.innerHTML = `<strong>${data.total_routes ?? routes.length}</strong><span>Total routes</span>`;

    const updated = document.createElement("div");
    updated.className = "bus-stat";
    updated.innerHTML = `<strong>${routes.filter((route) => route.status === "On Time").length}</strong><span>On time</span>`;

    const generated = document.createElement("div");
    generated.className = "bus-stat";
    generated.innerHTML = `<strong>${formatTimestamp(data.generated_at)}</strong><span>Updated</span>`;

    summary.append(total, updated, generated);

    // Search and Filter Controls
    const controls = document.createElement("div");
    controls.className = "bus-controls";
    controls.style.display = "flex";
    controls.style.gap = "10px";
    controls.style.marginBottom = "15px";

    const searchInput = document.createElement("input");
    searchInput.type = "text";
    searchInput.placeholder = "Search routes or depots...";
    searchInput.style.flex = "1";
    searchInput.style.padding = "8px 12px";

    const statusFilter = document.createElement("select");
    statusFilter.innerHTML = `
      <option value="">All Statuses</option>
      <option value="On Time">On Time</option>
      <option value="Delayed">Delayed</option>
    `;
    statusFilter.style.padding = "8px";

    controls.append(searchInput, statusFilter);

    const list = document.createElement("div");
    list.className = "bus-list";

    const renderList = (filteredRoutes: BusRouteSummary[]) => {
      list.replaceChildren();
      filteredRoutes.forEach((route) => {
        const routeId = asRouteId(route);
        const card = document.createElement("article");
        card.className = "bus-card";

        const title = document.createElement("h3");
        title.textContent = `${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim();

        const meta = document.createElement("p");
        meta.textContent = `${route.status ?? "Unknown status"} • ${route.depot_name ?? "Unknown depot"} • Next departure: ${route.next_departure ?? "—"}`;

        const details = document.createElement("div");
        details.className = "bus-meta";
        details.innerHTML = `
          <span>Headway: ${route.headway_minutes ?? 0} mins</span>
          <span>Vehicles: ${route.vehicles?.length ?? 0}</span>
        `;

        const actions = document.createElement("div");
        actions.className = "bus-actions";

        const timetableLink = document.createElement("a");
        timetableLink.href = `/bus/time/${encodeURIComponent(routeId)}`;
        timetableLink.textContent = "Timetable";

        const liveLink = document.createElement("a");
        liveLink.href = `/bus/live/${encodeURIComponent(routeId)}`;
        liveLink.textContent = "Live";

        actions.append(timetableLink, liveLink);
        card.append(title, meta, details, actions);
        list.appendChild(card);
      });
    };

    const applyFilters = () => {
      const query = searchInput.value.toLowerCase();
      const status = statusFilter.value;

      const filtered = routes.filter((route) => {
        const matchesQuery = 
          (route.name ?? "").toLowerCase().includes(query) ||
          (route.alias ?? "").toLowerCase().includes(query) ||
          (route.depot_name ?? "").toLowerCase().includes(query);
        
        const matchesStatus = !status || route.status === status;

        return matchesQuery && matchesStatus;
      });

      renderList(filtered);
    };

    searchInput.addEventListener("input", applyFilters);
    statusFilter.addEventListener("change", applyFilters);

    renderList(routes);
    wrapper.append(heading, summary, controls, list);
    container.appendChild(wrapper);
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unable to load bus overview.";
    const alert = document.createElement("div");
    alert.className = "bus-error";
    alert.textContent = message;
    wrapper.appendChild(alert);
    container.appendChild(wrapper);
  }
}