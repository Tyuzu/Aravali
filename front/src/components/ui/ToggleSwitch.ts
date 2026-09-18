import "../../../css/ui/ToggleSwitch.css";
import { createElement } from "../createElement.js";

// ---- Types & Interfaces ----

export type OnToggleCallback = (checked: boolean) => void;

export interface ToggleSwitchOptions {
  onToggle?: OnToggleCallback;
  checked?: boolean;
  disabled?: boolean;
  id?: string;
  label?: string;
  classes?: string;
  styles?: Partial<CSSStyleDeclaration> | Record<string, string>;
  [key: string]: unknown;
}

/**
 * Creates an accessible ToggleSwitch component supporting both options object and positional calls.
 */
const ToggleSwitch = (...args: any[]): HTMLLabelElement => {
  let opts: ToggleSwitchOptions = {};

  if (args.length === 1 && typeof args[0] === "object" && args[0] !== null) {
    opts = args[0] as ToggleSwitchOptions;
  } else {
    // Legacy / positional signature: ToggleSwitch(onToggle, checked?, label?, options?)
    if (typeof args[0] === "function") {
      opts.onToggle = args[0] as OnToggleCallback;
    }
    if (typeof args[1] === "boolean") opts.checked = args[1];
    if (typeof args[2] === "string") opts.label = args[2];
    if (typeof args[3] === "object" && args[3] !== null) {
      Object.assign(opts, args[3]);
    }
  }

  const {
    onToggle = () => {},
    checked = false,
    disabled = false,
    id = "",
    label = "",
    classes = "",
    styles = {},
    ...rest
  } = opts;

  // Generate fallback unique ID if omitted
  const switchId =
    id ||
    `toggle-${typeof crypto !== "undefined" && crypto.randomUUID ? crypto.randomUUID().slice(0, 8) : Math.random().toString(36).substring(2, 9)}`;

  const inputAttributes: Record<string, unknown> = {
    type: "checkbox",
    id: switchId,
    checked: Boolean(checked),
    disabled: Boolean(disabled),
    class: "sr-only",
    events: {
      change: (e: Event) => {
        const target = e.target as HTMLInputElement;
        if (typeof onToggle === "function") {
          onToggle(target.checked);
        }
      },
    },
  };

  const input = createElement("input", inputAttributes) as HTMLInputElement;
  const slider = createElement("span", {
    class: "slider",
    "aria-hidden": "true",
  });

  const trackWrapper = createElement("span", { class: "toggle-track" }, [
    input,
    slider,
  ]);

  const labelChildren: (HTMLElement | string)[] = [trackWrapper];

  if (label) {
    const labelText = createElement("span", { class: "toggle-label-text" }, [
      label,
    ]);
    labelChildren.push(labelText);
  } else {
    input.setAttribute("aria-label", "Toggle switch");
  }

  const containerClasses = [
    "toggle-switch",
    disabled ? "toggle-disabled" : "",
    classes,
  ]
    .filter(Boolean)
    .join(" ");

  return createElement(
    "label",
    {
      class: containerClasses,
      htmlFor: switchId,
      style: styles,
      ...rest,
    },
    labelChildren
  ) as HTMLLabelElement;
};

export { ToggleSwitch };
export default ToggleSwitch;
export { ToggleSwitch as ToggleSwitchComponent };