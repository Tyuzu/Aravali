import "../../../css/ui/Dropdown.css";
import { navigate } from "../../routes/navigate.js";

export interface DropdownOption {
  label: string;
  value?: string;
  href?: string;
}

export type DropdownItem = string | DropdownOption;

export interface DropdownProps {
  options?: DropdownItem[];
  onChange?: (value: string) => void;
  defaultText?: string;
}

export interface DropdownMenuItem {
  href?: string;
  text: string;
  icon?: string;
  onClick?: (e: MouseEvent) => void;
}

const Dropdown = ({
  options = [],
  onChange = () => {},
  defaultText = "Select an option",
}: DropdownProps = {}): HTMLDivElement => {
  const container = document.createElement("div");
  container.className = "dropdown";

  const button = document.createElement("button");
  button.className = "dropdown-button";
  button.type = "button";
  button.textContent = defaultText;

  const menu = document.createElement("ul");
  menu.className = "dropdown-menu";

  options.forEach((option) => {
    const label = typeof option === "string" ? option : option.label;
    const value = typeof option === "string" ? option : option.value ?? option.label;

    const menuItem = document.createElement("li");
    menuItem.className = "dropdown-item";

    if (typeof option === "object" && option.href) {
      const anchor = document.createElement("a");
      anchor.href = option.href as string;
      anchor.textContent = label;
      anchor.addEventListener("click", (e: MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        navigate(option.href as string);
        menu.classList.remove("show");
      });
      menuItem.appendChild(anchor);
    } else {
      menuItem.textContent = label;
      menuItem.addEventListener("click", (e: MouseEvent) => {
        e.stopPropagation();
        button.textContent = label;
        onChange(value);
        menu.classList.remove("show");
      });
    }

    menu.appendChild(menuItem);
  });

  button.addEventListener("click", (e: MouseEvent) => {
    e.stopPropagation();
    menu.classList.toggle("show");
  });

  // Close dropdown when clicking outside
  document.addEventListener("click", () => {
    menu.classList.remove("show");
  });

  container.appendChild(button);
  container.appendChild(menu);

  return container;
};

export function createDropdownMenu(
  id: string,
  labelText: string,
  items: DropdownMenuItem[],
  toggleElement?: HTMLElement
): HTMLDivElement {
  const toggle = toggleElement ?? document.createElement("button");
  if (!toggleElement) {
    toggle.id = id;
    toggle.className = "menu-toggle";
    (toggle as HTMLButtonElement).type = "button";
    toggle.textContent = labelText;
  } else {
    // ensure id is set for accessibility/consistency
    if (!toggle.id) toggle.id = id;
  }

  const menu = document.createElement("div");
  menu.className = "menu-content";
  menu.setAttribute("aria-label", labelText);

  items.forEach(({ href, text, icon, onClick }) => {
    const itemEl = href ? document.createElement("a") : document.createElement("button");
    itemEl.className = "profile-menu-item";
    if (href) (itemEl as HTMLAnchorElement).href = href as string;

    if (icon) {
      const iconSpan = document.createElement("span");
      iconSpan.innerHTML = icon;
      const label = document.createElement("span");
      label.textContent = text;
      itemEl.appendChild(iconSpan);
      itemEl.appendChild(label);
    } else {
      itemEl.textContent = text;
    }

    itemEl.addEventListener("click", (e: Event) => {
      e.preventDefault();
      e.stopPropagation();
      const me = e as MouseEvent;
      if (onClick) {
        onClick(me);
      } else if (href) {
        navigate(href as string);
      }
      menu.classList.remove("open");
    });

    menu.appendChild(itemEl);
  });

  toggle.addEventListener("click", (e: MouseEvent) => {
    e.stopPropagation();
    menu.classList.toggle("open");
  });

  document.addEventListener("click", () => menu.classList.remove("open"));

  const container = document.createElement("div");
  container.className = "header-content-dropdown";
  container.append(toggle, menu);
  return container;
}

export default Dropdown;