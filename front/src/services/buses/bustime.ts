import { createElement } from "../../components/createElement.js";
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
  const wrapper = createElement("div", { class: "bus-page" }, [createBusSkeleton()]) as HTMLDivElement;
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const data = await apiFetch<BusTimingResponse>("/bus/time", "GET");
      const routes = data.routes ?? [];

      wrapper.replaceChildren();

      const title = createElement("h2", {}, ["Bus Timetable Schedules"]);
      const info = createElement("p", {}, [
        data.generated ? `Timetable Updated: ${new Date(data.generated).toLocaleString()}` : "Live timetable schedules active"
      ]);

      const controls = createElement("div", { class: "bus-controls" }, [
        createElement("input", {
          type: "text",
          placeholder: "Search routes or stops..."
        }),
        createElement("label", {
          class: "bus-inline-btn",
          style: { display: "inline-flex", alignItems: "center", gap: "0.5rem" }
        }, [
          createElement("input", {
            type: "checkbox",
            "aria-label": "Show saved timetables only"
          }),
          createElement("span", {}, ["Saved only"])
        ])
      ]);

      const searchInput = controls.querySelector("input[type='text']") as HTMLInputElement;
      const favoriteToggle = controls.querySelector("input[type='checkbox']") as HTMLInputElement;
      const list = createElement("div", { class: "bus-list" }) as HTMLDivElement;

      const renderList = (filteredRoutes: BusRoute[]) => {
        list.replaceChildren();

        if (!filteredRoutes.length) {
          list.appendChild(createBusEmptyState("No schedules found", "Try a different route name or clear the saved filter."));
          return;
        }

        filteredRoutes.forEach((route) => {
          const routeId = asRouteId(route);
          const card = createElement("article", { class: "bus-card" }, [
            createElement("div", { class: "bus-card-header" }, [
              createElement("h3", {}, [`${route.name ?? "Route"} ${route.alias ? `(${route.alias})` : ""}`.trim()]),
              buildFavoriteButton(routeId, route.name ?? routeId)
            ]),
            createElement("p", {}, [`Status: ${route.status ?? "Unknown"} • Depot: ${route.depot_name ?? "Unknown depot"} • Next Departure: ${route.next_departure ?? "—"}`]),
            createElement("ol", { class: "bus-stop-list" }, (route.stops ?? []).map((stop) =>
              createElement("li", {}, [`${stop.order ?? "-"}. ${stop.name ?? "Stop"}${stop.arrival ? ` (Arrival: ${stop.arrival})` : ""}`])
            )),
            createElement("a", {
              href: `/bus/time/${encodeURIComponent(routeId)}`,
              style: { display: "inline-block", marginTop: "8px" }
            }, ["View Complete Route Details"])
          ]);
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
