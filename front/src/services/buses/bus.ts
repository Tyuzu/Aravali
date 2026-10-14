// buses/bus.ts

import { apiFetch } from "../../api/api.js";
import {
  buildFavoriteButton,
  createBusEmptyState,
  createBusErrorState,
  createBusSkeleton,
  readBusFavorites
} from "./busUtils.js";

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
  _isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  const wrapper = document.createElement("div");
  wrapper.className = "bus-page";
  wrapper.appendChild(createBusSkeleton());
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const data = await apiFetch<BusOverviewResponse>("/bustime", "GET");
      const routes = data.routes ?? [];

      wrapper.replaceChildren();

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

      const controls = document.createElement("div");
      controls.className = "bus-controls";

      const searchInput = document.createElement("input");
      searchInput.type = "text";
      searchInput.placeholder = "Search routes or depots...";

      const statusFilter = document.createElement("select");
      statusFilter.innerHTML = `
        <option value="">All Statuses</option>
        <option value="On Time">On Time</option>
        <option value="Delayed">Delayed</option>
      `;

      const favoriteOnly = document.createElement("label");
      favoriteOnly.className = "bus-inline-btn";
      favoriteOnly.style.display = "inline-flex";
      favoriteOnly.style.alignItems = "center";
      favoriteOnly.style.gap = "0.5rem";
      const favoriteToggle = document.createElement("input");
      favoriteToggle.type = "checkbox";
      favoriteToggle.setAttribute("aria-label", "Show saved routes only");
      const favoriteText = document.createElement("span");
      favoriteText.textContent = "Saved only";
      favoriteOnly.append(favoriteToggle, favoriteText);

      controls.append(searchInput, statusFilter, favoriteOnly);

      const list = document.createElement("div");
      list.className = "bus-list";

      const renderList = (filteredRoutes: BusRouteSummary[]) => {
        list.replaceChildren();

        if (!filteredRoutes.length) {
          list.appendChild(createBusEmptyState("No routes match", "Try a different route name or clear the filters."));
          return;
        }

        filteredRoutes.forEach((route) => {
          const routeId = asRouteId(route);
          const card = document.createElement("article");
          card.className = "bus-card";

          const header = document.createElement("div");
          header.className = "bus-card-header";

          const title = document.createElement("h3");
          title.textContent = `${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim();

          const favoriteButton = buildFavoriteButton(routeId, route.name ?? routeId);
          header.append(title, favoriteButton);

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
          card.append(header, meta, details, actions);
          list.appendChild(card);
        });
      };

      const applyFilters = () => {
        const query = searchInput.value.toLowerCase();
        const status = statusFilter.value;
        const favorites = readBusFavorites();

        const filtered = routes.filter((route) => {
          const routeId = asRouteId(route);
          const matchesQuery =
            (route.name ?? "").toLowerCase().includes(query) ||
            (route.alias ?? "").toLowerCase().includes(query) ||
            (route.depot_name ?? "").toLowerCase().includes(query) ||
            routeId.toLowerCase().includes(query);

          const matchesStatus = !status || route.status === status;
          const matchesFavorite = !favoriteToggle.checked || favorites.includes(routeId);

          return matchesQuery && matchesStatus && matchesFavorite;
        });

        renderList(filtered);
      };

      searchInput.addEventListener("input", applyFilters);
      statusFilter.addEventListener("change", applyFilters);
      favoriteToggle.addEventListener("change", applyFilters);

      renderList(routes);
      wrapper.append(heading, summary, controls, list);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unable to load bus overview.";
      wrapper.replaceChildren(createBusErrorState(message, load));
    }
  };

  await load();
}