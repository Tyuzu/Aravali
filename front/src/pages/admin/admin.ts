import "../../../css/inistyles/adminpage.css";
import { RoleRequestsPage } from "./roleRequests.js";

export async function Admin(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  await RoleRequestsPage(isLoggedIn, contentContainer);
}