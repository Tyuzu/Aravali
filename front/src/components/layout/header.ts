import "../../../css/layout/header.css";
import { getState, subscribe } from "../../state/state.js";
import { webSiteName } from "../../config/env.js";
import { navigate } from "../../routes/navigate.js";
import { moonSVG } from "../svgs/featherSVGs";
import { createElement } from "../createElement.js";
import { createDropdownMenu, DropdownMenuItem } from "../ui/Dropdown.js";
// import Imagex from "../base/Imagex.js";
// import { sticky } from "./sticky.js";
import Button from "../base/Button.js";
import { loadTheme, toggleTheme } from "./themeManager.js";
import createIconButton from "../ui/IconButton.js";
import { createProfileSection, getUserAvatarSrc, getCurrentUserState, UserState } from "./ProfileSection.js";

function renderUserSection(): HTMLDivElement {
  const container = createElement("div", { class: "user-area" }) as HTMLDivElement;
  const isLoggedIn = getState("isLoggedIn") ?? Boolean(getState("user")?.id || getState("user")?.userid);

  if (isLoggedIn) {
    container.append(createProfileSection());
  } else {
    const loginBtn = Button({
      title: "Login",
      id: "login-button",
      events: {
        click: () => navigate("/login")
      },
      classes: "login-btn",
      styles: { border: "none", cursor: "pointer" }
    });
    container.append(loginBtn);
  }

  return container;
}

function buildNav(): HTMLDivElement {
  const nav = createElement("div", { class: "header-content" }) as HTMLDivElement;
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

function createHeader(): void {
  const header = document.getElementById("pageheader");
  if (!header || header.hasChildNodes()) {
    return;
  }

  header.className = "main-header";

  const logo = createElement("div", { class: "logo" }, [
    createElement("a", { href: "/home", class: "logo-link" }, [webSiteName])
  ]);

  const sky = createElement("div", { class: "hflexcen" });

  // const renderSkyProfile = () => {
  //   sky.replaceChildren();
  //   const user = getCurrentUserState() as UserState;
  //   sky.append(
  //     sticky({
  //       imglink: Imagex({
  //         src: getUserAvatarSrc(user),
  //         alt: "Profile",
  //         classes: "profile-pic"
  //       })
  //     })
  //   );
  // };

  // renderSkyProfile();

  let navRef = buildNav();
  header.append(logo, sky, navRef);

  // Single top-level state listener updates navigation DOM cleanly
  const handleAuthChange = () => {
    // renderSkyProfile();
    const newNav = buildNav();
    navRef.replaceWith(newNav);
    navRef = newNav;
  };

  subscribe("isLoggedIn", handleAuthChange);
  subscribe("user", handleAuthChange);

  loadTheme();
}

export { createHeader as createheader };