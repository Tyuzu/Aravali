import { getState } from "../../state/state.js";
import { resolveImagePath, EntityType, PictureType } from "../../utils/imagePaths.js";
import Imagex from "../base/Imagex.js";
import { createElement } from "../createElement.js";
import { navigate } from "../../routes/navigate.js";
import { logout } from "../../services/auth/authService.js";
import { profileSVG, shopBagSVG, cardSVG, settingsSVG, logoutSVG } from "../svgs/featherSVGs";
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
    { href: "/wallet", text: "Wallet", icon: cardSVG },
    { href: "/settings", text: "Settings", icon: settingsSVG }
  ];

  const container = createDropdownMenu("profile-menu", username, links, toggle);

  // adjust classes to match previous structure
  container.className = "dropdown";
  const menu = container.querySelector(".menu-content") as HTMLDivElement;
  if (menu) menu.className = "profile-menu";

  // add logout button to menu
  const logoutBtn = createElement("button", { class: "profile-menu-item logout" }, []);
  logoutBtn.innerHTML = logoutSVG;
  logoutBtn.append(createElement("span", {}, ["Logout"]));
  logoutBtn.addEventListener("click", () => {
    if (menu) menu.classList.remove("open");
    logout();
  });
  if (menu) menu.appendChild(logoutBtn as unknown as Node);

  // keyboard handling for toggle (kept from previous behavior)
  toggle.addEventListener("keydown", (e: Event) => {
    const keyboardEvent = e as KeyboardEvent;
    if (keyboardEvent.key === "Enter" || keyboardEvent.key === " ") {
      keyboardEvent.preventDefault();
      if (menu) menu.classList.toggle("open");
    }
  });

  return container as HTMLDivElement;
}
