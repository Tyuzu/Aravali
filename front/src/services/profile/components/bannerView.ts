import Bannerx from "../../../components/base/Bannerx.js";

interface UserProfile {
  userid?: string | number;
  username?: string;
  banner?: string;
}

export function createBanner(profile: UserProfile = {} as UserProfile, isCreator: boolean = false): HTMLElement {
  return Bannerx({
    isCreator,
    bannerkey: profile.banner || "",
    banneraltkey: `Banner for ${profile.username || "User"}`,
    bannerentitytype: "user",
    stateentitykey: "user",
    bannerentityid: profile.userid ? String(profile.userid) : "",
    previewElementId: "banner-picture-preview"
  });
}
