import "../../../css/layout/navi.css";
import { t } from "../../i18n/i18n.js";
import { navigate } from "../../routes/navigate.js";
import { getCurrentAllowedFeatures } from "../../config/domainFeatures.js";
import { getState } from "../../state/state.js";
import { enableDragDrop, getNavOrder } from "./navigationDrag.js";

export interface NavItemConfig {
  href: string;
  label: string;
  feature?: string;
  roles?: string[];
}

const normalizeRoles = (value: unknown): string[] => {
  if (Array.isArray(value)) {
    return value
      .filter((role): role is string => typeof role === "string" && role.trim().length > 0)
      .map((role) => role.trim().toLowerCase());
  }

  if (typeof value === "string") {
    const trimmed = value.trim();
    return trimmed ? [trimmed.toLowerCase()] : [];
  }

  return [];
};

const getCurrentUserRoles = (): string[] => {
  const state = getState() as Record<string, any> | null;
  const authRoles = normalizeRoles(state?.auth?.roles ?? state?.roles);
  const userRoles = normalizeRoles(
    state?.userProfile?.roles ?? state?.user?.roles ?? state?.user?.role ?? state?.userProfile?.role
  );

  return [...new Set([...authRoles, ...userRoles])];
};

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

/** Filter nav items according to the active domain's feature flags and user roles */
const getPermittedNavItems = (allNavItems: NavItemConfig[]): NavItemConfig[] => {
  const allowed: string[] = getCurrentAllowedFeatures();
  const userRoles = getCurrentUserRoles();

  return allNavItems.filter((item) => {
    if (item.roles && item.roles.length > 0) {
      const hasRoleAccess = item.roles.some((role) => userRoles.includes(role.toLowerCase()));
      if (!hasRoleAccess) return false;
    }

    // Domain-level feature gate
    if (allowed.includes("ALL")) {
      return true;
    }

    // Shared core items without a feature key are always shown
    if (!item.feature) return true;
    return allowed.includes(item.feature);
  });
};

/** Create navigation bar */
const createNav = (): HTMLDivElement => {
  // 1. Master list of navigation items mapped to feature keys
  const allNavItems: NavItemConfig[] = [
    { href: "/dash", label: t("nav.dash", {}, "Dash"), feature: "farms", roles: ["farmer", "admin"] },
    { href: "/farms", label: t("nav.farms", {}, "Farms"), feature: "farms" },
    { href: "/grocery", label: t("nav.grocery", {}, "Grocery"), feature: "farms" },
    { href: "/recipes", label: t("nav.recipes", {}, "Recipes"), feature: "farms" },
    { href: "/products", label: t("nav.products", {}, "Products"), feature: "farms" },
    { href: "/tools", label: t("nav.tools", {}, "Tools"), feature: "farms" },
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
  toggleLabel.innerText = t("nav.more", {}, "More");

  toggleLabelWrapper.appendChild(toggleLabel);
  inner.appendChild(ul);
  inner.appendChild(toggleLabelWrapper);
  nav.appendChild(toggle);
  nav.appendChild(inner);

  highlightActiveNav(window.location.pathname);

  return nav;
};

export { createNav, createNavItem };