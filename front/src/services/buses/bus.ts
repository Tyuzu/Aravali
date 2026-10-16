import { createElement } from "../../components/createElement.js";
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
  const wrapper = createElement("div", { class: "bus-page" }, [createBusSkeleton()]) as HTMLDivElement;
  container.replaceChildren(wrapper);

  const load = async (): Promise<void> => {
    try {
      const data = await apiFetch<BusOverviewResponse>("/bustime", "GET");
      const routes = data.routes ?? [];

      wrapper.replaceChildren();

      const heading = createElement("h2", {}, ["Bus overview"]);

      const summary = createElement("div", { class: "bus-summary" }, [
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(data.total_routes ?? routes.length)]),
          createElement("span", {}, ["Total routes"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(routes.filter((route) => route.status === "On Time").length)]),
          createElement("span", {}, ["On time"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [formatTimestamp(data.generated_at)]),
          createElement("span", {}, ["Updated"])
        ])
      ]);

      const controls = createElement("div", { class: "bus-controls" }, [
        createElement("input", {
          type: "text",
          placeholder: "Search routes or depots..."
        }),
        createElement("select", {}, [
          createElement("option", { value: "" }, ["All Statuses"]),
          createElement("option", { value: "On Time" }, ["On Time"]),
          createElement("option", { value: "Delayed" }, ["Delayed"])
        ]),
        createElement("label", {
          class: "bus-inline-btn",
          style: { display: "inline-flex", alignItems: "center", gap: "0.5rem" }
        }, [
          createElement("input", {
            type: "checkbox",
            "aria-label": "Show saved routes only"
          }),
          createElement("span", {}, ["Saved only"])
        ])
      ]);

      const searchInput = controls.querySelector("input[type='text']") as HTMLInputElement;
      const statusFilter = controls.querySelector("select") as HTMLSelectElement;
      const favoriteToggle = controls.querySelector("input[type='checkbox']") as HTMLInputElement;
      const list = createElement("div", { class: "bus-list" }) as HTMLDivElement;

      const renderList = (filteredRoutes: BusRouteSummary[]) => {
        list.replaceChildren();

        if (!filteredRoutes.length) {
          list.appendChild(createBusEmptyState("No routes match", "Try a different route name or clear the filters."));
          return;
        }

        filteredRoutes.forEach((route) => {
          const routeId = asRouteId(route);
          const card = createElement("article", { class: "bus-card" }, [
            createElement("div", { class: "bus-card-header" }, [
              createElement("h3", {}, [
                `${route.name ?? "Bus route"} ${route.alias ? `(${route.alias})` : ""}`.trim()
              ]),
              buildFavoriteButton(routeId, route.name ?? routeId)
            ]),
            createElement("p", {}, [
              `${route.status ?? "Unknown status"} • ${route.depot_name ?? "Unknown depot"} • Next departure: ${route.next_departure ?? "—"}`
            ]),
            createElement("div", {
              class: "bus-meta"
            }, [
              createElement("span", {}, [`Headway: ${route.headway_minutes ?? 0} mins`]),
              createElement("span", {}, [`Vehicles: ${route.vehicles?.length ?? 0}`])
            ]),
            createElement("div", { class: "bus-actions" }, [
              createElement("a", { href: `/bus/time/${encodeURIComponent(routeId)}` }, ["Timetable"]),
              createElement("a", { href: `/bus/live/${encodeURIComponent(routeId)}` }, ["Live"])
            ])
          ]);
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