import Datex from "../../components/base/Datex.js";
import { createElement, ChildInput } from "../../components/createElement.js";
import Modal from "../../components/ui/Modal.js";
import { Button } from "../../components/base/Button.js";
import { cancelTicketRequest, fetchMyTickets, type UserTicket } from "./api.js";
import { printTicketPDF } from "./printTicket.js";

export async function listMyTickets(eventid: string | number): Promise<void> {
  const container = createElement("div", {}, []);

  try {
    const tickets = await fetchMyTickets(eventid);

    if (!tickets || tickets.length === 0) {
      container.append(
        createElement(
          "p",
          { style: { color: "#888", fontStyle: "italic" } },
          ["No tickets found for this event."]
        )
      );
    } else {
      tickets.sort((a, b) => {
        const order: Record<string, number> = { Active: 0, Transferred: 1, Cancelled: 2 };
        const aStatus = a.status ?? "Active";
        const bStatus = b.status ?? "Active";
        return (order[aStatus] ?? 9) - (order[bStatus] ?? 9);
      });

      tickets.forEach((ticket) => {
        const status = ticket.status ?? "Active";
        let statusColor = "green";
        if (status === "Cancelled") {
          statusColor = "red";
        }
        if (status === "Transferred") {
          statusColor = "orange";
        }

        const actionButtons: HTMLElement[] = [];
        const uniqueCode = ticket.uniquecode || ticket.uniqueCode || "";

        if (!ticket.canceled) {
          actionButtons.push(
            Button({
              title: "Print",
              events: {
                click: () => {
                  if (!uniqueCode) return;
                  printTicketPDF(eventid, uniqueCode);
                }
              },
              classes: "buttonx"
            }),
            Button({
              title: "Cancel",
              events: {
                click: async () => {
                  if (!confirm("Cancel this ticket?")) {
                    return;
                  }

                  try {
                    if (!uniqueCode) {
                      alert("Ticket code is missing.");
                      return;
                    }
                    await cancelTicketRequest(eventid, uniqueCode);
                    listMyTickets(eventid);
                  } catch {
                    alert("Failed to cancel ticket.");
                  }
                }
              },
              classes: "buttonx delete-btn"
            }),
            Button({
              title: "Transfer",
              events: {
                click: () =>
                  alert("Transfer modal not yet implemented.")
              },
              classes: "buttonx"
            })
          );
        }

        const meta: ChildInput[] = [];

        if (ticket.refundstatus) {
          meta.push(
            createElement(
              "p",
              { style: { color: "#444" } },
              [`Refund Status: ${ticket.refundstatus}`]
            )
          );
        }

        if (ticket.transferredto) {
          meta.push(
            createElement(
              "p",
              { style: { color: "#444" } },
              [`Transferred To: ${ticket.transferredto}`]
            )
          );
        }

        const ticketBox = createElement(
          "div",
          {
            style: {
              border: "1px solid #ccc",
              borderRadius: "8px",
              padding: "12px 14px",
              marginBottom: "12px",
              background: "#fafafa"
            }
          },
          [
            createElement("h4", {}, [`Ticket ID: ${ticket.ticketid}`]),
            createElement("p", {}, [`Unique Code: ${uniqueCode || "Unavailable"}`]),
            createElement("p", {}, [`Buyer: ${ticket.buyername || ticket.buyerName || "Unknown"}`]),
            createElement("p", {}, [`Purchase Date: ${Datex(ticket.purchasedate ?? ticket.purchaseDate ?? new Date())}`]),
            createElement(
              "p",
              { style: { fontWeight: "bold", color: statusColor } },
              [`Status: ${status}`]
            ),
            ...meta,
            actionButtons.length > 0
              ? createElement(
                  "div",
                  { style: { marginTop: "10px", display: "flex", gap: "8px" } },
                  actionButtons
                )
              : createElement("div", {}, [])
          ]
        );

        container.append(ticketBox);
      });
    }
  } catch (err) {
    console.error(err);
    container.append(
      createElement(
        "p",
        { style: { color: "red" } },
        ["Failed to load your tickets."]
      )
    );
  }

  Modal({
    title: "My Tickets",
    content: container,
    size: "large"
  });
}