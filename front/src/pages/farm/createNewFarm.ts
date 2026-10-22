import { createFarm } from "../../services/farm/createFarm.js";

export async function CreateFarm(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  createFarm(isLoggedIn, contentContainer);
}
