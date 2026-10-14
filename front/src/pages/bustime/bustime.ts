import "../../../css/inistyles/bustime.css";
import { displayBusTime } from "../../services/buses/bustime.js";

export async function BusTime(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBusTime(isLoggedIn, contentContainer);
}