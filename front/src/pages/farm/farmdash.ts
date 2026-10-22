import "../../../css/pagestyles/farmdash.css";
import "../../../css/pagestyles/orders.css";
import { displayFarmDash } from "../../services/farm/farmDash.js";

export async function Dash(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  await displayFarmDash(contentContainer, isLoggedIn);
}