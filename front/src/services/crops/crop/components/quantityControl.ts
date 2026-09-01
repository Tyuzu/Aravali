import { createElement } from "../../../../components/createElement";

export interface QuantityControl {
  element: HTMLElement;
  getQuantity: () => number;
  setDisabled: (val: boolean) => void;
  onChange: (fn: (q: number) => void) => void;
}

export function createQuantityControl(initial = 1, max = 99): QuantityControl {
  let quantity = initial;
  let disabled = false;
  const listeners: Array<(q: number) => void> = [];

  const quantityDisplay = createElement("span", {
    class: "quantity-value",
    "aria-live": "polite",
    "aria-label": "Selected quantity"
  }, [String(quantity)]) as HTMLElement;

  const updateQuantity = () => {
    quantityDisplay.textContent = String(quantity);
    listeners.forEach((fn) => fn(quantity));
  };

  const decrementBtn = createElement("button", { type: "button", "aria-label": "Decrease quantity" }, ["−"]) as HTMLButtonElement;
  const incrementBtn = createElement("button", { type: "button", "aria-label": "Increase quantity" }, ["+"]) as HTMLButtonElement;

  decrementBtn.addEventListener("click", () => {
    if (disabled) return;
    if (quantity > 1) {
      quantity -= 1;
      updateQuantity();
    }
  });

  incrementBtn.addEventListener("click", () => {
    if (disabled) return;
    if (quantity < max) {
      quantity += 1;
      updateQuantity();
    }
  });

  const wrapper = createElement("div", { class: "quantity-control", role: "group", "aria-label": "Quantity" }, [decrementBtn, quantityDisplay, incrementBtn]);

  return {
    element: wrapper as HTMLElement,
    getQuantity: () => quantity,
    setDisabled: (val: boolean) => {
      disabled = Boolean(val);
      decrementBtn.disabled = disabled;
      incrementBtn.disabled = disabled;
    },
    onChange: (fn: (q: number) => void) => {
      listeners.push(fn);
    }
  };
}
