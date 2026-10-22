
import "../../../css/pagestyles/farms.css";
import { displayFarms } from "../../services/farm/FarmsHome.js";

export async function Farms(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayFarms(contentContainer, isLoggedIn);
}
