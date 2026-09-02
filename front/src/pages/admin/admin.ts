import "../../../css/inistyles/adminpage.css";
import { displayAdminDash } from "../../services/admin/adminDash.js";

export async function Admin(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  await displayAdminDash(contentContainer, isLoggedIn);
}