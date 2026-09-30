import "../../../css/layout/navi.css";
import { t } from "../../i18n/i18n.js";
import { navigate } from "../../routes/navigate.js";
import { getCurrentAllowedFeatures } from "../../config/domainFeatures.js";
import { getState, subscribe } from "../../state/state.js";
import { enableDragDrop, getNavOrder } from "./navigationDrag.js";
import { createElement } from "../createElement.js";

export interface NavItemConfig {
  href: string;
  label: string;
  feature?: string;
  roles?: string[];
}

/** Get normalized current user roles directly from state.js alias */
const getCurrentUserRoles = (): string[] => {
  const roles = getState("roles");
  if (Array.isArray(roles)) {
    return roles.map((role) => String(role).toLowerCase());
  }
  return [];
};

/** Highlight current active link within a target container or entire document */
export const highlightActiveNav = (path: string, container: ParentNode = document): void => {
  container.querySelectorAll<HTMLAnchorElement>(".navigation__link").forEach((link) => {
    link.classList.toggle("active", link.getAttribute("href") === path);
  });
};

/** Handle navigation cleanly */
const handleNavigation = (event: MouseEvent, href: string): void => {
  event.preventDefault();
  if (!href) {
    console.error("🚨 handleNavigation received invalid href!");
    return;
  }
  navigate(href);
};

/** Create individual navigation list item */
const createNavItem = (href: string, label: string): HTMLLIElement => {
  const anchor = createElement("a", {
    href,
    class: "navigation__link",
    events: {
      click: (e: Event) => handleNavigation(e as MouseEvent, href)
    }
  }, [label]);

  return createElement("li", {
    class: "navigation__item",
    draggable: "false"
  }, [anchor]) as HTMLLIElement;
};

/** Filter nav items against domain feature flags & user state roles */
const getPermittedNavItems = (allNavItems: NavItemConfig[]): NavItemConfig[] => {
  const allowedFeatures: string[] = getCurrentAllowedFeatures();
  const userRoles = getCurrentUserRoles();

  return allNavItems.filter((item) => {
    if (item.roles && item.roles.length > 0) {
      const hasRoleAccess = item.roles.some((role) => userRoles.includes(role.toLowerCase()));
      if (!hasRoleAccess) return false;
    }

    if (allowedFeatures.includes("ALL")) {
      return true;
    }

    if (!item.feature) return true;
    return allowedFeatures.includes(item.feature);
  });
};

/** Build nav links fragment based on master configuration and ordering */
const buildNavList = (): HTMLUListElement => {
  const allNavItems: NavItemConfig[] = [
    // { href: "/dash", label: t("nav.dash", {}, "Dash"), feature: "farms", roles: ["farmer", "admin"] },
    // { href: "/farms", label: t("nav.farms", {}, "Farms"), feature: "farms" },
    // { href: "/grocery", label: t("nav.grocery", {}, "Grocery"), feature: "farms" },
    // { href: "/recipes", label: t("nav.recipes", {}, "Recipes"), feature: "farms" },
    // { href: "/products", label: t("nav.products", {}, "Products"), feature: "farms" },
    // { href: "/tools", label: t("nav.tools", {}, "Tools"), feature: "farms" },
    // { href: "/places", label: t("nav.places", {}, "Places"), feature: "places" }
    { href: "/places", label: t("nav.places", {}, "Places"), feature: "places" }
  ];

  const defaultNavItems = getPermittedNavItems(allNavItems);
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

  const ul = createElement("ul", { class: "navigation__list horizontal" }) as HTMLUListElement;
  navItems.forEach(({ href, label }) => ul.appendChild(createNavItem(href, label)));

  return ul;
};

/** Create full navigation bar container with dynamic state re-rendering support */
const createNav = (): HTMLDivElement => {
  const toggle = createElement("input", {
    class: "toggle",
    type: "checkbox",
    id: "more",
    tabindex: "-1"
  }) as HTMLInputElement;

  let ul = buildNavList();
  enableDragDrop(ul, toggle);

  const toggleLabel = createElement("label", {
    class: "navigation__link",
    for: "more"
  }, [t("nav.more", {}, "More")]);

  const toggleLabelWrapper = createElement("div", { class: "navigation__toggle" }, [toggleLabel]);
  const inner = createElement("div", { class: "navigation__inner" }, [ul, toggleLabelWrapper]);

  const nav = createElement("div", { class: "navigation" }, [toggle, inner]) as HTMLDivElement;

  // Reactively re-render menu items when user updates roles/login state
  const refreshNavList = () => {
    const newUl = buildNavList();
    ul.replaceWith(newUl);
    ul = newUl;
    enableDragDrop(ul, toggle);
    highlightActiveNav(window.location.pathname, nav);
  };

  const unsubRoles = subscribe("roles", refreshNavList);
  const unsubAuth = subscribe("isLoggedIn", refreshNavList);

  // Auto-cleanup observer when component is unmounted from document body
  queueMicrotask(() => {
    const observer = new MutationObserver(() => {
      if (!document.body.contains(nav)) {
        unsubRoles?.();
        unsubAuth?.();
        observer.disconnect();
      }
    });

    if (document.body) {
      observer.observe(document.body, { childList: true, subtree: true });
    }
  });

  highlightActiveNav(window.location.pathname, nav);

  return nav;
};

export { createNav, createNavItem };