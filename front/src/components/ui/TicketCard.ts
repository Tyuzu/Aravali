import "../../../css/ui/TicketCard.css";
import { createElement } from "../../components/createElement.js";
import { applyButtonColors } from "../../utils/lumicolor.js";
import Button from "../base/Button.js";

// ---- Types & Interfaces ----

export interface TicketCardProps {
  isl?: boolean;
  seatstart?: string | number;
  seatend?: string | number;
  creator?: boolean;
  name: string;
  price: string | number;
  quantity: number;
  color: string;
  attributes?: Record<string, unknown>;
  onClick?: (name: string, quantity: number) => void;
}

/**
 * Creates an enhanced, accessible TicketCard component node.
 */
const TicketCard = ({
  isl,
  seatstart,
  seatend,
  creator,
  name,
  price,
  quantity,
  color,
  attributes = {},
  onClick,
}: TicketCardProps): HTMLDivElement => {
  // Folded corner accent
  const corner = createElement("div", {
    class: "corner",
    style: { borderTop: `40px solid ${color}` },
  });

  // Left color stripe with perforation effect
  const perforation = createElement("div", { class: "perforation" });
  const stripeContainer = createElement(
    "div",
    {
      class: "stripe-container",
      style: { backgroundColor: color },
    },
    [perforation]
  );

  // Main content elements
  const nameElement = createElement("h2", {
    class: "ticket-title",
    style: { color },
  }, [name]);

  const contentChildren: (HTMLElement | string)[] = [nameElement];

  if (seatstart && seatend) {
    const seats = createElement("p", { class: "ticket-seats" }, [
      `Seats: ${seatstart} – ${seatend}`
    ]);
    contentChildren.push(seats);
  }

  const priceElement = createElement("span", { class: "ticket-price" }, [`Price: ${price}`]);
  const availableElement = createElement("span", { class: "ticket-available" }, [`Available: ${quantity}`]);

  const info = createElement("div", { class: "tickinfo" }, [
    priceElement,
    availableElement,
  ]);

  contentChildren.push(info);

  // Action button (Buy / Sold Out)
  if (!creator && isl) {
    const button = Button({ 
      title: quantity > 0 ? "Buy Ticket" : "Sold Out", 
      id: "", 
      styles: "buttonx primary" 
    }) as HTMLButtonElement;

    if (quantity > 0) {
      if (onClick) {
        button.addEventListener("click", () => onClick(name, quantity));
      }
      applyButtonColors(button, color);
    } else {
      button.disabled = true;
      button.classList.add("sold-out");
      Object.assign(button.style, { backgroundColor: "var(--color-bg-tertiary, #f0f0f0)", color: "var(--color-muted, #616161)" });
    }
    contentChildren.push(button);
  }

  const content = createElement("div", { class: "content" }, contentChildren);

  // Root card element merging dynamic custom attributes
  const card = createElement(
    "div",
    {
      class: "ticket-card",
      role: "region",
      "aria-label": `Ticket for ${name}`,
      ...attributes,
    },
    [corner, stripeContainer, content]
  ) as HTMLDivElement;

  return card;
};

export default TicketCard;