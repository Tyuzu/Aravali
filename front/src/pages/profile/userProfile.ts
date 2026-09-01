
import "../../../css/inistyles/profilexx.css";
import "../../../css/inistyles/udata.css";
import { displayProfile } from "../../services/profile/displayMyProfile.js";
import { displayUserProfile } from "../../services/profile/otherUserProfileService.js";
import { getState } from "../../state/state.js";

export async function MyProfile(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  const content = document.createElement("div");
  content.className = "profilepage";
  contentContainer.appendChild(content);

  // Prefer centralized state snapshot to determine auth, fall back to route param
  const logged = Boolean(getState("isLoggedIn") || isLoggedIn);
  await displayProfile(Boolean(logged), content);
}

export async function UserProfile(
  isLoggedIn: boolean,
  username: string,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  const content = document.createElement("div");
  content.className = "profilepage";
  contentContainer.appendChild(content);
  displayUserProfile(isLoggedIn, content, username);
}