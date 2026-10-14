import "../../../css/inistyles/bustime.css";
import { displayRouteDetails } from "../../services/buses/busRouteDetail.js";

export async function BusRouteDetail(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  await displayRouteDetails(isLoggedIn, contentContainer, params);
}
