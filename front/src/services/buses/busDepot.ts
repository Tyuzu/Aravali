import { createElement } from "../../components/createElement.js";
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
  _isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const depotId = asDepotId(params);
  container.replaceChildren();

  try {
    const depot = await apiFetch<BusDepot>(`/bus/depot/${encodeURIComponent(depotId)}`, "GET");

    const wrapper = createElement("div", { class: "bus-page" }, [
      createElement("h2", {}, [depot.name ?? "Bus depot"]),
      createElement("div", { class: "bus-summary" }, [
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(depot.vehicles ?? "—")]),
          createElement("span", {}, ["Vehicles"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(depot.open_hours ?? "—")]),
          createElement("span", {}, ["Open hours"])
        ]),
        createElement("div", { class: "bus-stat" }, [
          createElement("strong", {}, [String(depot.location ?? "—")]),
          createElement("span", {}, ["Location"])
        ])
      ]),
      createElement("ul", { class: "bus-stop-list" }, [
        createElement("li", {}, [
          "Phone: ",
          depot.phone ? createElement("a", { href: `tel:${encodeURIComponent(depot.phone)}` }, [depot.phone]) : "—"
        ]),
        createElement("li", {}, [`Coordinate: ${depot.latitude ?? "—"}, ${depot.longitude ?? "—"}`]),
        createElement("li", {}, [`Description: ${depot.description ?? "—"}`])
      ])
    ]) as HTMLDivElement;

    container.appendChild(wrapper);
  } catch (error) {
    const wrapper = createElement("div", { class: "bus-page" }, [
      createElement("h2", {}, [`Bus depot ${depotId}`]),
      createElement("p", {}, [error instanceof Error ? error.message : "Unable to load depot data."])
    ]) as HTMLDivElement;
    container.appendChild(wrapper);
  }
}