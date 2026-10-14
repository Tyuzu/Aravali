import "../../../css/inistyles/bustime.css";
import { displayBus } from "../../services/buses/bus.js";

export async function BusHome(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  await displayBus(isLoggedIn, contentContainer);
}

export async function Bus(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  await BusHome(isLoggedIn, contentContainer);
}