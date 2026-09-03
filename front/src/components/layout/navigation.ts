import "../../../css/layout/navi.css";
import { navigate } from "../../routes/navigate.js";
import { getCurrentAllowedFeatures } from "../../config/domainFeatures.js";
import { enableDragDrop, getNavOrder } from "./navigationDrag.js";

export interface NavItemConfig {
  href: string;
  label: string;
  feature?: string;
}

/** Highlight current active link */
export const highlightActiveNav = (path: string): void => {
  document.querySelectorAll<HTMLAnchorElement>(".navigation__link").forEach((link) => {
    link.classList.toggle("active", link.getAttribute("href") === path);
  });
};

/** Handle navigation */
const handleNavigation = (event: MouseEvent, href: string): void => {
  event.preventDefault();
  if (!href) {
    console.error("🚨 handleNavigation received null href!");
    return;
  }
  navigate(href);
};

/** Create one navigation item */
const createNavItem = (href: string, label: string): HTMLLIElement => {
  const li = document.createElement("li");
  li.className = "navigation__item";

  // Start as non-draggable so normal clicks fire cleanly
  li.setAttribute("draggable", "false");

  const anchor = document.createElement("a");
  anchor.href = href;
  anchor.className = "navigation__link";
  anchor.textContent = label;
  anchor.addEventListener("click", (e: MouseEvent) => handleNavigation(e, href));

  li.appendChild(anchor);
  return li;
};

/** Filter nav items according to the active domain's feature flags */
const getPermittedNavItems = (allNavItems: NavItemConfig[]): NavItemConfig[] => {
  const allowed: string[] = getCurrentAllowedFeatures();

  // If domain allows ALL features, return everything
  if (allowed.includes("ALL")) {
    return allNavItems;
  }

  return allNavItems.filter((item) => {
    // Shared core items without a feature key are always shown
    if (!item.feature) return true;
    return allowed.includes(item.feature);
  });
};

/** Create navigation bar */
const createNav = (): HTMLDivElement => {
  // 1. Master list of navigation items mapped to feature keys
  const allNavItems: NavItemConfig[] = [
    { href: "/dash", label: "Dash", feature: "farms" },
    { href: "/farms", label: "Farms", feature: "farms" },
    { href: "/grocery", label: "Grocery", feature: "farms" },
    { href: "/recipes", label: "Recipes", feature: "farms" },
  ];

  // 2. Filter available items based on domain permissions
  const defaultNavItems = getPermittedNavItems(allNavItems);

  // 3. Apply custom drag-and-drop ordering (stored in localStorage)
  const savedOrder = getNavOrder();
  let navItems = defaultNavItems;

  if (savedOrder) {
    navItems = savedOrder
      .map((href) => defaultNavItems.find((item) => item.href === href))
      .filter((item): item is NavItemConfig => Boolean(item));

    defaultNavItems.forEach((item) => {
      if (!navItems.find((i) => i.href === item.href)) {
        navItems.push(item);
      }
    });
  }

  const nav = document.createElement("div");
  nav.className = "navigation";

  const toggle = document.createElement("input");
  toggle.className = "toggle";
  toggle.type = "checkbox";
  toggle.id = "more";
  toggle.setAttribute("tabindex", "-1");

  const inner = document.createElement("div");
  inner.className = "navigation__inner";

  const ul = document.createElement("ul");
  ul.className = "navigation__list horizontal";

  navItems.forEach(({ href, label }) => ul.appendChild(createNavItem(href, label)));

  enableDragDrop(ul, toggle);

  const toggleLabelWrapper = document.createElement("div");
  toggleLabelWrapper.className = "navigation__toggle";

  const toggleLabel = document.createElement("label");
  toggleLabel.className = "navigation__link";
  toggleLabel.setAttribute("for", "more");
  toggleLabel.innerText = "More";

  toggleLabelWrapper.appendChild(toggleLabel);
  inner.appendChild(ul);
  inner.appendChild(toggleLabelWrapper);
  nav.appendChild(toggle);
  nav.appendChild(inner);

  highlightActiveNav(window.location.pathname);

  return nav;
};

export { createNav, createNavItem };