// buses/bustime.ts

import { apiFetch } from "../../api/api.js";
import {
  buildFavoriteButton,
  createBusEmptyState,
  createBusErrorState,
  createBusSkeleton,
  readBusFavorites
} from "./busUtils.js";

type BusStop = {
  id?: string;
  name?: string;
  code?: string;
  order?: number;
  arrival?: string;
};

type BusRoute = {
  route_id?: string;
  id?: string;
  name?: string;
  alias?: string;
  status?: string;
  next_departure?: string;
  depot_name?: string;
  stops?: BusStop[];
};

type BusTimingResponse = {
  routes?: BusRoute[];
  generated?: string;
};

const asRouteId = (route: BusRoute): string => route.route_id ?? route.id ?? "unknown";

export async function displayBusTime(
  _isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  const wrapper = document.createElement("div");
  wrapper.className = "bus-page";
  wrapper.appendChild(createBusSkeleton());
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const data = await apiFetch<BusTimingResponse>("/bus/time", "GET");
      const routes = data.routes ?? [];

      wrapper.replaceChildren();

      const title = document.createElement("h2");
      title.textContent = "Bus Timetable Schedules";

      const info = document.createElement("p");
      info.textContent = data.generated ? `Timetable Updated: ${new Date(data.generated).toLocaleString()}` : "Live timetable schedules active";

      const controls = document.createElement("div");
      controls.className = "bus-controls";

      const searchInput = document.createElement("input");
      searchInput.type = "text";
      searchInput.placeholder = "Search routes or stops...";

      const favoriteOnly = document.createElement("label");
      favoriteOnly.className = "bus-inline-btn";
      favoriteOnly.style.display = "inline-flex";
      favoriteOnly.style.alignItems = "center";
      favoriteOnly.style.gap = "0.5rem";
      const favoriteToggle = document.createElement("input");
      favoriteToggle.type = "checkbox";
      favoriteToggle.setAttribute("aria-label", "Show saved timetables only");
      const favoriteText = document.createElement("span");
      favoriteText.textContent = "Saved only";
      favoriteOnly.append(favoriteToggle, favoriteText);

      controls.append(searchInput, favoriteOnly);

      const list = document.createElement("div");
      list.className = "bus-list";

      const renderList = (filteredRoutes: BusRoute[]) => {
        list.replaceChildren();

        if (!filteredRoutes.length) {
          list.appendChild(createBusEmptyState("No schedules found", "Try a different route name or clear the saved filter."));
          return;
        }

        filteredRoutes.forEach((route) => {
          const routeId = asRouteId(route);
          const card = document.createElement("article");
          card.className = "bus-card";

          const header = document.createElement("div");
          header.className = "bus-card-header";

          const heading = document.createElement("h3");
          heading.textContent = `${route.name ?? "Route"} ${route.alias ? `(${route.alias})` : ""}`.trim();

          const favoriteButton = buildFavoriteButton(routeId, route.name ?? routeId);
          header.append(heading, favoriteButton);

          const summary = document.createElement("p");
          summary.textContent = `Status: ${route.status ?? "Unknown"} • Depot: ${route.depot_name ?? "Unknown depot"} • Next Departure: ${route.next_departure ?? "—"}`;

          const stopList = document.createElement("ol");
          stopList.className = "bus-stop-list";

          (route.stops ?? []).forEach((stop) => {
            const item = document.createElement("li");
            item.textContent = `${stop.order ?? "-"}. ${stop.name ?? "Stop"}${stop.arrival ? ` (Arrival: ${stop.arrival})` : ""}`;
            stopList.appendChild(item);
          });

          const detailsLink = document.createElement("a");
          detailsLink.href = `/bus/time/${encodeURIComponent(routeId)}`;
          detailsLink.textContent = "View Complete Route Details";
          detailsLink.style.display = "inline-block";
          detailsLink.style.marginTop = "8px";

          card.append(header, summary, stopList, detailsLink);
          list.appendChild(card);
        });
      };

      const applyFilters = () => {
        const query = searchInput.value.toLowerCase();
        const favorites = readBusFavorites();
        const filtered = routes.filter((route) => {
          const routeId = asRouteId(route);
          const matchesQuery =
            (route.name ?? "").toLowerCase().includes(query) ||
            (route.alias ?? "").toLowerCase().includes(query) ||
            (route.stops ?? []).some((stop) => (stop.name ?? "").toLowerCase().includes(query)) ||
            routeId.toLowerCase().includes(query);
          const matchesFavorite = !favoriteToggle.checked || favorites.includes(routeId);
          return matchesQuery && matchesFavorite;
        });
        renderList(filtered);
      };

      searchInput.addEventListener("input", applyFilters);
      favoriteToggle.addEventListener("change", applyFilters);

      renderList(routes);
      wrapper.append(title, info, controls, list);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unable to load timetable data.";
      wrapper.replaceChildren(createBusErrorState(message, load));
    }
  };

  await load();
}
