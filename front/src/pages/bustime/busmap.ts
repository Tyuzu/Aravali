import "../../../css/inistyles/bustime.css";
import { displayBusMap } from "../../services/buses/busmap.js";

export async function BusMap(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBusMap(isLoggedIn, contentContainer);
}
