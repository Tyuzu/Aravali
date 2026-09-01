import profilGen from "../profilegen.js";
import type { UserProfile } from "../../profile/profileGenHelpers.js";

export function renderProfile(profile: UserProfile, isLoggedIn: boolean, onLoadUserData?: (isLoggedIn: boolean, container: HTMLElement, username: string) => void): HTMLElement {
  const profileElement = profilGen(profile, isLoggedIn, onLoadUserData as any);
  return profileElement;
}
