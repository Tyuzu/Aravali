import { apiFetch } from "../../api/api.js";
import type { Ticket as TicketData, TicketPayload, UserTicket } from "./types.js";

export type { TicketData, TicketPayload, UserTicket };
export type TicketStatus = UserTicket["status"] extends string ? UserTicket["status"] : string;

export async function fetchTicketData(
  ticketId: string | number,
  eventId: string | number
): Promise<TicketData> {
  return await apiFetch<TicketData>(`/ticket/event/${eventId}/${ticketId}`, "GET");
}

export async function updateTicketRequest(
  ticketId: string | number,
  eventId: string | number,
  payload: TicketPayload
): Promise<void> {
  await apiFetch<void>(`/ticket/event/${eventId}/${ticketId}`, "PUT", payload);
}

export async function deleteTicketRequest(
  ticketId: string | number,
  eventId: string | number
): Promise<void> {
  await apiFetch<void>(`/ticket/event/${eventId}/${ticketId}`, "DELETE");
}

export async function fetchMyTickets(eventid: string | number): Promise<UserTicket[]> {
  return await apiFetch<UserTicket[]>(`/ticket/mytickets/${eventid}`, "GET");
}

export async function cancelTicketRequest(
  eventid: string | number,
  uniquecode: string
): Promise<void> {
  await apiFetch<void>(`/ticket/cancel/${eventid}`, "POST", { uniquecode });
}

export default {
  fetchTicketData,
  updateTicketRequest,
  deleteTicketRequest,
  fetchMyTickets,
  cancelTicketRequest
};
