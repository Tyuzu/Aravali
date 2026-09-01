import "../../../css/layout/header.css";
import { getState, subscribe } from "../../state/state.js";
import { webSiteName } from "../../config/env.js";
import { navigate } from "../../routes/navigate.js";
import { logout } from "../../services/auth/authService.js";
import { settingsSVG, moonSVG, profileSVG, shopBagSVG, logoutSVG, cardSVG } from "../svgs/featherSVGs";
import { createElement } from "../createElement.js";
import { resolveImagePath, EntityType, PictureType } from "../../utils/imagePaths.js";
import Imagex from "../base/Imagex.js";
import { sticky } from "./sticky.js";
import Button from "../base/Button.js";
import { loadTheme, toggleTheme } from "./themeManager.js";

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

function getCurrentUserState(): Partial<UserState> {
  const authUser = (getState("user") || {}) as Partial<UserState>;
  const profileUser = (getState("userProfile") || {}) as Partial<UserState>;
  return {
    ...profileUser,
    ...authUser
  };
}

function getUserAvatarSrc(user: Partial<UserState> = {}): string {
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

function createIconButton(
  svg: string,
  href?: string | null,
  onClick?: (e: MouseEvent) => void
): HTMLDivElement {
  const icon = createElement("span", { class: "icon" }, []);
  icon.innerHTML = svg;

  const anchor = createElement("div", { class: "iconic-button" }, [icon]) as HTMLDivElement;
  if (href) {
    // Use bracket notation here 👇
    (anchor as unknown as Record<string, unknown>)['href'] = href;
  }
  if (onClick) {
    anchor.addEventListener("click", onClick as EventListener);
  }

  return anchor;
}

function createDropdownMenu(
  id: string,
  labelText: string,
  items: DropdownMenuItem[]
): HTMLDivElement {
  const toggle = createElement("button", { id, class: "menu-toggle" }, [labelText]);
  const menu = createElement("div", { class: "menu-content", "aria-label": labelText }, []);

  items.forEach(({ href, text }) => {
    const link = createElement("a", { class: "profile-menu-item", href }, [text]);
    link.addEventListener("click", (e) => {
      e.preventDefault();
      navigate(href);
      menu.classList.remove("open");
    });
    menu.append(link);
  });

  toggle.addEventListener("click", (e) => {
    e.stopPropagation();
    menu.classList.toggle("open");
  });

  return createElement("div", { class: "header-content-dropdown" }, [toggle, menu]) as HTMLDivElement;
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

  const links: ProfileMenuItem[] = [
    { href: "/profile", text: username, icon: profileSVG },
    { href: "/my-orders", text: "My Orders", icon: shopBagSVG },
    { href: "/wallet", text: "Wallet", icon: cardSVG },
    { href: "/settings", text: "Settings", icon: settingsSVG }
  ];

  const menu = createElement("div", { class: "profile-menu" }, []);

  links.forEach(({ href, text, icon }) => {
    const label = createElement("span", {}, [text]);
    const iconSpan = createElement("span", {}, []);
    if (icon) {
      iconSpan.innerHTML = icon;
    }

    const link = createElement("a", { class: "profile-menu-item", href }, [iconSpan, label]);
    link.addEventListener("click", (e) => {
      e.preventDefault();
      navigate(href);
      menu.classList.remove("open");
    });

    menu.append(link);
  });

  const logoutBtn = createElement("button", { class: "profile-menu-item logout" }, []);
  logoutBtn.innerHTML = logoutSVG;
  logoutBtn.append(createElement("span", {}, ["Logout"]));
  logoutBtn.addEventListener("click", () => {
    menu.classList.remove("open");
    logout();
  });
  menu.append(logoutBtn);

  toggle.addEventListener("click", (e) => {
    e.stopPropagation();
    menu.classList.toggle("open");
  });

  toggle.addEventListener("keydown", (e: Event) => {
    const keyboardEvent = e as KeyboardEvent;
    if (keyboardEvent.key === "Enter" || keyboardEvent.key === " ") {
      keyboardEvent.preventDefault();
      menu.classList.toggle("open");
    }
  });

  document.addEventListener("click", () => menu.classList.remove("open"));

  return createElement("div", { class: "dropdown" }, [toggle, menu]) as HTMLDivElement;
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
  let navRef: HTMLDivElement = initialNavRef;

  function updateNav(): void {
    if (!navRef || !navRef.parentNode) return;
    const newNav = buildNav();
    navRef.replaceWith(newNav);
    navRef = newNav;
  }

  subscribe("isLoggedIn", updateNav);
  subscribe("user", updateNav);
}

function createHeader(): void {
  const header = document.getElementById("pageheader");
  if (!header || header.hasChildNodes()) {
    return;
  }

  header.className = "main-header";

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
  loadTheme();
}

export { createHeader as createheader };