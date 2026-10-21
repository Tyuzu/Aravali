
// import "../../../css/inistyles/profilexx.css";
import { displayProfile } from "../services/profile/displayMyProfile.js";
import { getState } from "../state/state.js";

export function Home(isLoggedIn: boolean, container: HTMLElement): void {
  container.innerHTML = "";
  const content = document.createElement("div");
  content.className = "profilepage";
  container.appendChild(content);

  // Prefer centralized state snapshot to determine auth, fall back to route param
  const logged = Boolean(getState("isLoggedIn") || isLoggedIn);
  displayProfile(Boolean(logged), content);
}