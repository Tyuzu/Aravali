import profilGen from "../profilegen.js";
import type { UserProfile } from "../../profile/profileGenHelpers";
import type { LoadUserDataCallback } from "../../profile/profilegen.js";

export function renderProfile(
  profile: UserProfile,
  isLoggedIn: boolean,
  onLoadUserData?: LoadUserDataCallback
): HTMLElement {
  return profilGen(profile, isLoggedIn, onLoadUserData || null);
}