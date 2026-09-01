import "../../../css/ui/ToggleSwitch.css";
import { createElement } from "../createElement.js";

// ---- Types & Interfaces ----

export interface ToggleSwitchOptions {
  checked?: boolean;
  disabled?: boolean;
  id?: string;
  label?: string;
}

export type OnToggleCallback = (checked: boolean) => void;

/**
 * Creates an accessible ToggleSwitch component.
 */
const ToggleSwitch = (
  onToggle: OnToggleCallback,
  {
    checked = false,
    disabled = false,
    id = "",
    label = "",
  }: ToggleSwitchOptions = {}
): HTMLLabelElement => {
  // Fallback unique ID generation if no ID is passed
  const switchId = id || `toggle-${Math.random().toString(36).substring(2, 9)}`;

  const inputAttributes: Record<string, unknown> = {
    type: "checkbox",
    id: switchId,
    checked: Boolean(checked),
    disabled: Boolean(disabled),
    "aria-checked": String(Boolean(checked)),
    class: "sr-only",
    events: {
      change: (e: Event) => {
        const target = e.target as HTMLInputElement;
        target.setAttribute("aria-checked", String(target.checked));
        onToggle(target.checked);
      },
    },
  };

  const input = createElement("input", inputAttributes) as HTMLInputElement;
  const slider = createElement("span", {
    class: "slider",
    "aria-hidden": "true",
  });

  const trackWrapper = createElement(
    "span",
    { class: "toggle-track" },
    [input, slider]
  );

  const labelChildren: (HTMLElement | string)[] = [trackWrapper];

  if (label) {
    const labelText = createElement("span", { class: "toggle-label-text" }, [
      label,
    ]);
    labelChildren.push(labelText);
  } else {
    input.setAttribute("aria-label", "Toggle switch");
  }

  const labelAttributes: Record<string, unknown> = {
    class: `toggle-switch${disabled ? " toggle-disabled" : ""}`,
    htmlFor: switchId,
  };

  const toggleContainer = createElement(
    "label",
    labelAttributes,
    labelChildren
  ) as HTMLLabelElement;

  return toggleContainer;
};

export default ToggleSwitch;