import "../../../css/inistyles/bustime.css";
import { displayBusRoute } from "../../services/buses/busRoutePage.js";

export async function BusRoute(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBusRoute(isLoggedIn, contentContainer, params);
}
