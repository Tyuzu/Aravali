import { createProfileActions } from "./createProfileActions.js";
import { createProfileInfo } from "./createProfileInfo.js";
import type { UserProfile } from "../profileGenHelpers.js";

export function createProfileDetails(profile: UserProfile, isLoggedIn: boolean): HTMLDivElement {
    const profileDetails = document.createElement("div");
    profileDetails.className = "profile-details";

    const username = document.createElement("strong");
    username.className = "username";
    username.textContent = `@${profile.username}`;

    const name = document.createElement("p");
    name.className = "name";
    name.textContent = profile.name || "";

    const bio = document.createElement("p");
    bio.className = "bio";
    bio.textContent = profile.bio || "";

    const profileActions = createProfileActions(profile, isLoggedIn);
    const profileInfo = createProfileInfo(profile);

    profileDetails.append(profileActions, name, username, bio, profileInfo);
    return profileDetails;
}
