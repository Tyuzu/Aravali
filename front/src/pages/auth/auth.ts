// import "../../../css/inistyles/authpage.css";
import { renderAuth } from "../../services/auth/auth.js";

export function Auth(
  isLoggedIn: boolean,
  contentContainer: HTMLElement | null
): void {
  renderAuth(isLoggedIn, contentContainer);
}
