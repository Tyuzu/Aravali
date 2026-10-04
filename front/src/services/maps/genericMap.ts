import { createElement } from "../../components/createElement.js";

export interface MapPoint {
  x?: number;
  y?: number;
  lat?: number;
  lon?: number;
  name?: string;
  type?: string;
}

export interface MapOptions {
  entity?: string;
  title?: string;
  mapImage?: string;
  mapWidth?: number;
  mapHeight?: number;
  mapBounds?: { minLat: number; maxLat: number; minLon: number; maxLon: number };
  currentLocation?: { lat: number; lon: number };
  markers?: Array<MapPoint & { lat?: number; lon?: number; name?: string; type?: string }>;
  showLegend?: boolean;
  locations?: Array<{ id?: string; name?: string; x?: number; y?: number; category?: string }>;
  lockedAreas?: Array<{ id?: string; label?: string; x?: number; y?: number; width?: number; height?: number }>;
  floors?: Array<{ level?: number; name?: string; image?: string }>;
}

export function displayDynamicMap(container: HTMLElement, options: MapOptions): void {
  if (!container) return;
  const mapEl = createElement("div", { class: "dynamic-map-placeholder" }, [
    createElement("p", {}, ["Map placeholder"])
  ]) as HTMLElement;
  container.appendChild(mapEl);
}
