import "../../../css/inistyles/bustime.css";
import { displayBusDepot } from "../../services/buses/busDepot.js";

export async function BusDepot(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBusDepot(isLoggedIn, contentContainer, params);
}
