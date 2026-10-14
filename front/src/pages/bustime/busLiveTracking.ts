import "../../../css/inistyles/bustime.css";
import { displayBusLiveTracking } from "../../services/buses/busLiveTracking.js";

export async function BusLiveTracking(
  isLoggedIn: boolean,
  params: Record<string, string | undefined> | string | undefined,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBusLiveTracking(isLoggedIn, contentContainer, params);
}
