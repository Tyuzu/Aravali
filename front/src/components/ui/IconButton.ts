import { createElement } from "../createElement.js";

export default function createIconButton(
  svg: string,
  href?: string | null,
  onClick?: (e: MouseEvent) => void
): HTMLDivElement {
  const icon = createElement("span", { class: "icon" }, []);
  icon.innerHTML = svg;

  const anchor = createElement("div", { class: "iconic-button" }, [icon]) as HTMLDivElement;
  if (href) {
    (anchor as unknown as Record<string, unknown>)["href"] = href;
  }
  if (onClick) {
    anchor.addEventListener("click", onClick as EventListener);
  }

  return anchor;
}
