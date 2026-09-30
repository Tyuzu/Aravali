import { getState } from "../../state/state.js";
import { resolveImagePath, EntityType, PictureType } from "../../utils/imagePaths.js";
import Imagex from "../base/Imagex.js";
import { createElement } from "../createElement.js";
import { logout } from "../../services/auth/authService.js";
import { profileSVG, shopBagSVG, settingsSVG, logoutSVG } from "../svgs/featherSVGs";
import { createDropdownMenu, DropdownMenuItem } from "../ui/Dropdown.js";

export interface UserState {
  id?: string;
  userid?: string;
  username?: string;
  name?: string;
  avatar?: string;
  profilepicture?: string;
  profileImage?: string;
  image?: string;
  picture?: string;
  role?: string;
  [key: string]: unknown;
}

export function getCurrentUserState(): Partial<UserState> {
  const authUser = (getState("user") || {}) as Partial<UserState>;
  const profileUser = (getState("userProfile") || {}) as Partial<UserState>;
  return {
    ...profileUser,
    ...authUser
  };
}

export function getUserAvatarSrc(user: Partial<UserState> = {}): string {
  const mergedUser = {
    ...getCurrentUserState(),
    ...user
  };

  const avatar =
    typeof mergedUser.avatar === "string" ? mergedUser.avatar :
    typeof mergedUser.profilepicture === "string" ? mergedUser.profilepicture :
    typeof mergedUser.profileImage === "string" ? mergedUser.profileImage :
    typeof mergedUser.image === "string" ? mergedUser.image :
    typeof mergedUser.picture === "string" ? mergedUser.picture :
    "";

  if (avatar) {
    if (/^https?:\/\//i.test(avatar) || avatar.startsWith("/")) {
      return avatar;
    }
    return resolveImagePath(EntityType.USER, PictureType.THUMB, avatar);
  }

  const userid = mergedUser.userid || mergedUser.id || "default";
  return resolveImagePath(EntityType.USER, PictureType.THUMB, `${userid}.jpg`);
}

export function createProfileSection(): HTMLDivElement {
  const user = getCurrentUserState() as UserState;
  const username = user.username || user.name || "Profile";

  const img = Imagex({
    src: getUserAvatarSrc(user),
    alt: username,
    classes: "profile-pic"
  });

  const toggle = createElement("div", { class: "profile-toggle", tabIndex: 0 }, [img]);

  const links: DropdownMenuItem[] = [
    { href: "/profile", text: username, icon: profileSVG },
    { href: "/my-orders", text: "My Orders", icon: shopBagSVG },
    { href: "/settings", text: "Settings", icon: settingsSVG }
  ];

  const container = createDropdownMenu("profile-menu", username, links, toggle) as HTMLDivElement;

  container.className = "dropdown";
  const menu = container.querySelector(".menu-content") as HTMLDivElement;
  if (menu) {
    menu.className = "profile-menu";

    // Add logout button
    const logoutBtn = createElement("button", { class: "profile-menu-item logout" }, [
      createElement("span", {}, ["Logout"])
    ]);
    logoutBtn.insertAdjacentHTML("afterbegin", logoutSVG);
    logoutBtn.addEventListener("click", () => {
      menu.classList.remove("open");
      logout();
    });
    menu.appendChild(logoutBtn);
  }

  // Keyboard accessibility handler
  toggle.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      if (menu) menu.classList.toggle("open");
    }
  });

  return container;
}