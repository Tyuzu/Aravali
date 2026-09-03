import TicketCard from "../../components/ui/TicketCard.js";
import { Button } from "../../components/base/Button.js";
import Modal from "../../components/ui/Modal.js";
import Notify from "../../components/ui/Notify.js";
import { createElement } from "../../components/createElement.js";
import { createFormGroup } from "../../components/form/createFormGroupEnhanced.js";
import { apiFetch } from "../../api/api.js";

import { deleteTicket, editTicket, clearTicketForm } from "./editTicket.js";
import { printTicket } from "./printTicket.js";
import { listMyTickets } from "./listMyTickets.js";
import { showPaymentModal } from "../pay/pay.js";
import {
  verifyTicketAndShowModal,
  cancelTicket,
  transferTicket
} from "./ticketTransfer.js";

export interface Ticket {
  ticketid: string | number;
  name: string;
  price: number;
  quantity: number;
  currency?: string;
  color?: string;
  seatstart?: number;
  seatend?: number;
  [key: string]: unknown;
}

export interface PaymentResult {
  success: boolean;
  [key: string]: unknown;
}

export interface ApiFetchTicketResponse {
  success?: boolean;
  message?: string;
  data?: Ticket[];
  [key: string]: unknown;
}

/* ────────── Helpers ────────── */
function formatCurrency(minorValue: number, currencyCode = "INR"): string {
  const code = currencyCode.toUpperCase();
  const divisor = code === "JPY" ? 1 : 100;

  const localeMap: Record<string, string> = {
    INR: "en-IN",
    USD: "en-US",
    EUR: "de-DE",
    GBP: "en-GB"
  };
  const locale = localeMap[code] || navigator.language || "en-US";

  return new Intl.NumberFormat(locale, {
    style: "currency",
    currency: code
  }).format((minorValue || 0) / divisor);
}

function formatSeatRange(ticket: Ticket): string {
  const seatStart = Number(ticket.seatstart);
  const seatEnd = Number(ticket.seatend);

  if (Number.isFinite(seatStart) && Number.isFinite(seatEnd) && seatStart > 0 && seatEnd > 0) {
    return `Seats ${seatStart}–${seatEnd}`;
  }

  return "Seat allocation not specified";
}

/* ────────── Ticket Card ────────── */
function createTicketCard(
  ticket: Ticket,
  eventId: string | number,
  isCreator: boolean,
  isLoggedIn: boolean,
  onRefresh?: () => void
): HTMLElement {
  const card = TicketCard({
    isl: isLoggedIn,
    seatstart: ticket.seatstart,
    seatend: ticket.seatend,
    creator: isCreator,
    name: ticket.name,
    price: formatCurrency(ticket.price, ticket.currency),
    quantity: ticket.quantity,
    color: ticket.color || "#a3a3a349",
    attributes: { "data-ticket-id": String(ticket.ticketid) },
    onClick: () => {
      if (!isLoggedIn || isCreator) {
        return;
      }

      const availableQuantity = Math.max(0, Number(ticket.quantity) || 0);
      if (availableQuantity <= 0) {
        Notify("This ticket is sold out.", {
          type: "warning",
          dismissible: true
        });
        return;
      }

      const currencyCode = String(ticket.currency || "INR").toUpperCase();
      const unitPrice = Number(ticket.price) || 0;
      const quantityInput = createElement("input", {
        type: "number",
        min: 1,
        max: availableQuantity,
        step: 1,
        value: 1,
        style: { width: "100%", maxWidth: "120px" }
      });

      const totalLabel = createElement("p", {
        style: { margin: "8px 0 0", fontWeight: "bold" }
      }, [`Total: ${formatCurrency(unitPrice * 1, currencyCode)}`]);

      quantityInput.addEventListener("input", () => {
        const quantity = parseInt(quantityInput.value, 10);
        const safeQuantity = Number.isInteger(quantity) && quantity > 0 ? quantity : 0;
        const cappedQuantity = Math.min(safeQuantity, availableQuantity);
        totalLabel.textContent = `Total: ${formatCurrency(unitPrice * cappedQuantity, currencyCode)}`;
      });

      const summary = createElement("div", {
        style: {
          display: "grid",
          gap: "8px",
          background: "#f8f8f8",
          border: "1px solid #e5e5e5",
          borderRadius: "8px",
          padding: "12px",
          marginBottom: "12px"
        }
      }, [
        createElement("p", { style: { margin: "0", fontWeight: "600" } }, [`${ticket.name}`]),
        createElement("p", { style: { margin: "0" } }, [`Price: ${formatCurrency(unitPrice, currencyCode)}`]),
        createElement("p", { style: { margin: "0" } }, [`Available: ${availableQuantity}`]),
        createElement("p", { style: { margin: "0" } }, [formatSeatRange(ticket)]),
        totalLabel
      ]);

      const wrapper = createElement("div", { class: "modal-form-group" }, [
        summary,
        createElement("label", { style: { display: "grid", gap: "6px" } }, [
          "Quantity",
          quantityInput
        ])
      ]);

      const modal = Modal({
        title: `Buy ${ticket.name}`,
        content: wrapper,
        actions: () =>
          createElement("div", { class: "modal-actions" }, [
            Button({
              title: "Confirm Purchase",
              classes: "buttonx primary",
              events: {
                click: () => {
                  void (async () => {
                    const quantity = parseInt(quantityInput.value, 10);

                    if (
                      !Number.isInteger(quantity) ||
                      quantity < 1 ||
                      quantity > availableQuantity
                    ) {
                      Notify(
                        `Enter a valid quantity (1-${availableQuantity}).`,
                        { type: "warning", dismissible: true }
                      );
                      return;
                    }

                    modal.close();

                    try {
                      const paymentResult = await showPaymentModal({
                        paymentType: "purchase",
                        entityType: "ticket",
                        entityId: ticket.ticketid,
                        entityName: ticket.name
                      }) as PaymentResult | null;

                      if (!paymentResult?.success) {
                        Notify("Payment cancelled or failed.", {
                          type: "error",
                          dismissible: true
                        });
                        return;
                      }

                      const resp = await apiFetch<ApiFetchTicketResponse>(
                        `/ticket/event/${eventId}/${ticket.ticketid}/confirm-purchase`,
                        "POST",
                        { quantity }
                      );

                      if (resp?.success) {
                        Notify(`Ticket${quantity > 1 ? "s" : ""} purchased successfully!`, {
                          type: "success",
                          dismissible: true
                        });
                      } else {
                        Notify(resp?.message || "Purchase failed.", {
                          type: "error",
                          dismissible: true
                        });
                      }
                    } catch (err: unknown) {
                      console.error("Ticket purchase failed:", err);
                      const errorMessage = err instanceof Error ? err.message : String(err);
                      Notify(`Purchase failed: ${errorMessage}`, {
                        type: "error",
                        dismissible: true
                      });
                    }
                  })();
                }
              }
            }),
            Button({
              title: "Cancel",
              classes: "buttonx",
              events: { click: () => modal.close() }
            })
          ])
      });
    }
  });

  if (isCreator) {
    const actions = createElement("div", {
      class: "hflex-sb",
      style: { padding: "0 0.5rem" }
    });

    actions.append(
      Button({
        title: "Edit",
        classes: "buttonx primary",
        events: { click: () => editTicket(ticket.ticketid, eventId, onRefresh) }
      }),
      Button({
        title: "Delete",
        classes: "buttonx delete-btn",
        events: { click: () => deleteTicket(ticket.ticketid, eventId, onRefresh) }
      })
    );
    card.append(actions);
  }

  return card;
}

export function displayNewTicket(
  ticketData: Ticket,
  ticketList: HTMLElement,
  isCreator = false,
  isLoggedIn = false,
  eventId: string | number,
  onRefresh?: () => void
): void {
  ticketList.append(createTicketCard(ticketData, eventId, isCreator, isLoggedIn, onRefresh));
}

export async function displayTickets(
  ticketContainer: HTMLElement,
  eventId: string | number,
  isCreator: boolean,
  isLoggedIn: boolean
): Promise<void> {
  let tickets: Ticket[] = [];

  const handleRefresh = () => {
    void displayTickets(ticketContainer, eventId, isCreator, isLoggedIn);
  };

  try {
    const resp = await apiFetch<ApiFetchTicketResponse>(`/ticket/event/${eventId}`);
    tickets = resp?.data ?? [];
  } catch (err) {
    console.error("Failed to load tickets:", err);
    ticketContainer.replaceChildren(
      createElement("p", {}, ["Error loading tickets."])
    );
    return;
  }

  ticketContainer.replaceChildren(createElement("h2", {}, ["Tickets"]));
  const actionsCon = createElement("div", { class: "hvflex" });

  if (!isCreator && tickets.length > 0) {
    actionsCon.append(
      Button({
        title: "Verify Ticket",
        classes: "buttonx action-btn",
        events: { click: () => { void verifyTicketAndShowModal(eventId); } }
      }),
      Button({
        title: "Print Ticket",
        classes: "buttonx action-btn",
        events: { click: () => { void printTicket(eventId); } }
      }),
      Button({
        title: "Cancel Ticket",
        classes: "buttonx action-btn",
        events: { click: () => { void cancelTicket(eventId); } }
      }),
      Button({
        title: "Transfer Ticket",
        classes: "buttonx action-btn",
        events: { click: () => { void transferTicket(eventId); } }
      }),
      Button({
        title: "My Tickets",
        classes: "buttonx action-btn",
        events: { click: () => { void listMyTickets(eventId); } }
      })
    );
  }

  const ticketListDiv = createElement("div", { class: "hvflex gap20", id: "ticket-list" });

  if (isCreator) {
    ticketContainer.append(
      Button({
        title: "Add Tickets",
        id: "add-ticket-btn",
        classes: "buttonx",
        events: { click: () => { void addTicketForm(eventId, ticketListDiv); } }
      })
    );
  }

  if (tickets.length > 0) {
    tickets.forEach((t) =>
      ticketListDiv.append(createTicketCard(t, eventId, isCreator, isLoggedIn, handleRefresh))
    );
  } else {
    ticketListDiv.append(
      createElement("p", {}, ["No tickets available for this event."])
    );
  }

  ticketContainer.append(actionsCon, ticketListDiv);
}

/* ────────── Add Ticket API & Form ────────── */
async function handleAddTicketSubmit(
  form: HTMLFormElement,
  eventId: string | number,
  ticketList: HTMLElement,
  modalInstance: { close: () => void } | null
): Promise<void> {
  const formData = new FormData(form);

  const nameValue = formData.get("name");
  const currencyValue = formData.get("currency");
  const colorValue = formData.get("color");

  const payload = {
    name: typeof nameValue === "string" ? nameValue.trim() : "",
    price: Number(formData.get("price")),
    quantity: Number(formData.get("quantity")),
    currency: typeof currencyValue === "string" ? currencyValue : "INR",
    color: typeof colorValue === "string" ? colorValue : "#f3f3f3",
    seatstart: Number(formData.get("seatstart") || 0),
    seatend: Number(formData.get("seatend") || 0)
  };

  if (
    !payload.name ||
    payload.price <= 0 ||
    payload.quantity <= 0 ||
    payload.seatstart > payload.seatend
  ) {
    Notify("Please enter valid ticket details.", {
      type: "warning",
      dismissible: true,
      duration: 3000
    });
    return;
  }

  try {
    const ticket = await apiFetch<Ticket>(`/ticket/event/${eventId}`, "POST", payload);

    if (ticket?.ticketid) {
      Notify("Ticket added successfully.", {
        type: "success",
        dismissible: true,
        duration: 3000
      });

      displayNewTicket(ticket, ticketList, true, true, eventId);
      clearTicketForm();
      modalInstance?.close?.();
    } else {
      Notify("Failed to add ticket.", { type: "error", dismissible: true });
    }
  } catch (err) {
    console.error("Error adding ticket:", err);
    Notify("Error adding ticket.", { type: "error", dismissible: true });
  }
}

export function addTicketForm(eventId: string | number, ticketList: HTMLElement): void {
  const form = createElement<HTMLFormElement>("form", { id: "add-ticket-form" }, []);

  const fields = [
    { label: "Ticket Name", type: "text", id: "ticket-name", name: "name", required: true },
    { label: "Ticket Price (minor unit)", type: "number", id: "ticket-price", name: "price", required: true },
    { label: "Quantity", type: "number", id: "ticket-quantity", name: "quantity", required: true },
    { label: "Seat Start", type: "number", id: "seat-start", name: "seatstart" },
    { label: "Seat End", type: "number", id: "seat-end", name: "seatend" }
  ];

  fields.forEach((f) => form.append(createFormGroup(f)));

  /* Currency Select */
  const currencySelect = createElement("select", {
    id: "ticket-currency",
    name: "currency",
    required: true
  });

  ["INR", "USD", "EUR", "GBP", "CAD", "AUD", "JPY"].forEach((c) =>
    currencySelect.append(createElement("option", { value: c }, [c]))
  );

  form.append(
    createElement("div", { class: "form-group" }, [
      createElement("label", { for: "ticket-currency" }, ["Currency"]),
      currencySelect
    ])
  );

  /* Color Input */
  form.append(
    createElement("div", { class: "form-group" }, [
      createElement("label", { for: "ticket-color" }, ["Ticket Color"]),
      createElement("input", {
        id: "ticket-color",
        name: "color",
        type: "color",
        value: "#f3f3f3"
      })
    ])
  );

  const modal = Modal({
    title: "Add Ticket",
    content: form,
    actions: () =>
      createElement("div", { class: "modal-actions" }, [
        Button({
          title: "Add Ticket",
          classes: "buttonx primary",
          events: { click: () => { form.requestSubmit(); } }
        }),
        Button({
          title: "Cancel",
          classes: "buttonx",
          events: { click: () => modal.close() }
        })
      ])
  });

  form.addEventListener("submit", (e) => {
    e.preventDefault();
    void handleAddTicketSubmit(form, eventId, ticketList, modal);
  });
}