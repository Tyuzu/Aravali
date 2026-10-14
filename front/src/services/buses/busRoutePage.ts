// buses/busRoutePage.ts

import { displayRouteDetails } from "./busRouteDetail.js";

export async function displayBusRoute(
  isLoggedIn: boolean,
  container: HTMLElement,
  params?: Record<string, string | undefined> | string
): Promise<void> {
  await displayRouteDetails(isLoggedIn, container, params);
}