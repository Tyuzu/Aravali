
import "../../../css/pagestyles/crops.css";
import { displayCrops } from "../../services/crops/displayCrops.js";

export async function Crops(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayCrops(contentContainer);
}
