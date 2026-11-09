import "../../../css/layout/header.css";
import { getState, subscribe } from "../../state/state.js";
import { webSiteName } from "../../config/envConfig.js";
import { navigate } from "../../routes/navigate.js";
import { logout } from "../../services/auth/authService.js";
import { settingsSVG, moonSVG, profileSVG, shopBagSVG, logoutSVG, cardSVG } from "../svgs/featherSVGs";
import { createElement } from "../createElement.js";
import { createDropdownMenu } from "../ui/Dropdown.js";
import { resolveImagePath, EntityType, PictureType } from "../../utils/imagePaths.js";
import Imagex from "../base/Imagex.js";
import { sticky } from "./sticky.js";
import Button from "../base/Button.js";
import { loadTheme, toggleTheme } from "./themeManager.js";
import createIconButton from "../ui/IconButton.js";
import { createProfileSection, getUserAvatarSrc, getCurrentUserState } from "./ProfileSection.js";

export interface DropdownMenuItem {
  href: string;
  text: string;
}

export interface ProfileMenuItem extends DropdownMenuItem {
  icon?: string;
}

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

function renderUserSection(): HTMLDivElement {
  const container = createElement("div", { class: "user-area" }, []) as HTMLDivElement;

  function update(): void {
    container.replaceChildren();
    const isLoggedIn = getState("isLoggedIn") ?? Boolean(getState("user")?.id || getState("user")?.userid);

    if (isLoggedIn) {
      container.append(createProfileSection());
    } else {
      const loginBtn = Button({
        title: "Login", id: "login-button", events: {
          click: () => {
            navigate("/login");
          }
        }, classes: "login-btn", styles: { border: "none", cursor: "pointer" }
      });

      container.append(loginBtn);
    }
  }

  subscribe("isLoggedIn", update);
  subscribe("user", update);

  update();
  return container;
}

function buildNav(): HTMLDivElement {
  const nav = createElement("div", { class: "header-content" }, []) as HTMLDivElement;
  const isLoggedIn = getState("isLoggedIn") ?? Boolean(getState("user")?.id || getState("user")?.userid);

  if (isLoggedIn) {
    const createLinks: DropdownMenuItem[] = [
      { href: "/create-farm", text: "Farm" },
      { href: "/create-recipe", text: "Recipe" }
    ];
    nav.append(createDropdownMenu("create-menu", "Create", createLinks));
  }

  nav.append(
    createIconButton(moonSVG, null, toggleTheme),
    renderUserSection()
  );

  return nav;
}

function enableNavAutoUpdate(initialNavRef: HTMLDivElement): void {
  const navEl = initialNavRef as HTMLDivElement & { __headerNavBound?: boolean };
  if (navEl.__headerNavBound) {
    return;
  }
  navEl.__headerNavBound = true;

  let navRef: HTMLDivElement = navEl;

  function updateNav(): void {
    if (!navRef || !navRef.parentNode) return;
    const newNav = buildNav();
    navRef.replaceWith(newNav);
    navRef = newNav;
  }

  subscribe("isLoggedIn", updateNav);
  subscribe("user", updateNav);
}

function renderHeaderContents(): void {
  const header = document.getElementById("pageheader");
  if (!header) {
    return;
  }

  header.className = "main-header";
  header.replaceChildren();

  const logo = createElement("div", { class: "logo" }, [
    createElement("a", { href: "/home", class: "logo-link" }, [webSiteName])
  ]);

  const user = getCurrentUserState() as UserState;

  const sky = createElement("div", { class: "hflexcen" }, []);
  sky.append(
    sticky({
      imglink: Imagex({
        src: getUserAvatarSrc(user),
        alt: "Profile",
        classes: "profile-pic"
      })
    })
  );

  const nav = buildNav();
  header.append(logo, sky, nav);
  enableNavAutoUpdate(nav);
}

function bindHeaderAuthRefresh(): void {
  const header = document.getElementById("pageheader") as HTMLElement & { __headerAuthBound?: boolean } | null;
  if (!header || header.__headerAuthBound) {
    return;
  }

  header.__headerAuthBound = true;

  const refreshHeader = () => {
    renderHeaderContents();
  };

  subscribe("isLoggedIn", refreshHeader);
  subscribe("user", refreshHeader);
  subscribe("token", refreshHeader);

  window.addEventListener("auth:logout", refreshHeader);
  window.addEventListener("auth:unauthorized", refreshHeader);
}

function createHeader(): void {
  const header = document.getElementById("pageheader");
  if (!header) {
    return;
  }

  bindHeaderAuthRefresh();
  renderHeaderContents();
  loadTheme();
}

export { createHeader as createheader };