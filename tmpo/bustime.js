import { apiFetch } from "../../api/api.js";

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
  isLoggedIn: boolean,
  container: HTMLElement
): Promise<void> {
  container.replaceChildren();

  if (!isLoggedIn) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = "<h2>Bus Timetable</h2><p>Sign in to view bus timetable data.</p>";
    container.appendChild(wrapper);
    return;
  }

  try {
    const data = await apiFetch<BusTimingResponse>("/bus/time", "GET");
    const routes = data.routes ?? [];

    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";

    const title = document.createElement("h2");
    title.textContent = "Bus timetable";

    const info = document.createElement("p");
    info.textContent = data.generated ? `Updated: ${new Date(data.generated).toLocaleString()}` : "Live timetable data";

    const list = document.createElement("div");
    list.className = "bus-list";

    routes.forEach((route) => {
      const routeId = asRouteId(route);
      const card = document.createElement("article");
      card.className = "bus-card";

      const heading = document.createElement("h3");
      heading.textContent = `${route.name ?? "Route"} ${route.alias ? `(${route.alias})` : ""}`.trim();

      const summary = document.createElement("p");
      summary.textContent = `${route.status ?? "Unknown"} • ${route.depot_name ?? "Unknown depot"} • Next departure: ${route.next_departure ?? "—"}`;

      const stopList = document.createElement("ol");
      stopList.className = "bus-stop-list";

      (route.stops ?? []).forEach((stop) => {
        const item = document.createElement("li");
        item.textContent = `${stop.order ?? "-"}. ${stop.name ?? "Stop"}${stop.arrival ? ` • ${stop.arrival}` : ""}`;
        stopList.appendChild(item);
      });

      const detailsLink = document.createElement("a");
      detailsLink.href = `/bus/time/${encodeURIComponent(routeId)}`;
      detailsLink.textContent = "Route details";

      card.append(heading, summary, stopList, detailsLink);
      list.appendChild(card);
    });

    wrapper.append(title, info, list);
    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = document.createElement("div");
    wrapper.className = "bus-page";
    wrapper.innerHTML = `<h2>Bus timetable</h2><p>${error instanceof Error ? error.message : "Unable to load timetable data."}</p>`;
    container.appendChild(wrapper);
  }
}