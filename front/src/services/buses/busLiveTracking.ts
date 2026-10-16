import { createElement } from "../../components/createElement.js";
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
  _isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  const routeId = asRouteId(params);
  container.replaceChildren();

  const title = createElement("h2", {}, [`Live tracking: ${routeId}`]);
  const updated = createElement("p", {}, []);
  const list = createElement("div", { class: "bus-list" });
  const wrapper = createElement("div", { class: "bus-page" }, [title, updated, list]) as HTMLDivElement;

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
        list.appendChild(createElement("p", {}, ["No active vehicles currently reporting location on this route."]));
        return;
      }

      vehicles.forEach((vehicle) => {
        const capacityValue = vehicle.capacity ?? 50;
        const card = createElement("article", { class: "bus-card" }, [
          createElement("h3", {}, [`Vehicle #${vehicle.id ?? "Unknown"} (${vehicle.status ?? "Active"})`]),
          createElement("p", {}, [`GPS Coordinates: ${vehicle.lat ?? "—"}, ${vehicle.lng ?? "—"}`]),
          createElement("div", {
            style: {
              marginTop: "8px",
              fontSize: "12px",
              marginBottom: "3px"
            }
          }, [`Estimated Capacity: ${capacityValue}%`]),
          createElement("div", {
            style: {
              background: "#e2e8f0",
              borderRadius: "4px",
              height: "8px",
              width: "100%",
              overflow: "hidden"
            }
          }, [
            createElement("div", {
              style: {
                background: capacityValue > 80 ? "#ef4444" : capacityValue > 50 ? "#f59e0b" : "#10b981",
                width: `${capacityValue}%`,
                height: "100%"
              }
            })
          ])
        ]) as HTMLElement;

        list.appendChild(card);
      });
    } catch (error) {
      updated.textContent = error instanceof Error ? error.message : "Error syncing live coordinates.";
    }
  };

  await fetchData();
  timer = window.setInterval(fetchData, 15000);

  const observer = new MutationObserver(() => {
    if (!document.body.contains(container)) {
      if (timer) clearInterval(timer);
      observer.disconnect();
    }
  });
  observer.observe(document.body, { childList: true, subtree: true });
}