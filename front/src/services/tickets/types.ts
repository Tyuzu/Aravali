export type TicketId = string | number;
export type DateLike = string | number | Date;
export type TicketStatus = "Active" | "Transferred" | "Cancelled" | string;

export interface TicketSeat {
  id?: string;
  _id?: string;
  entity_id?: string;
  entity_type?: string;
  seat_number?: string;
  userid?: string;
  status?: string;
  [key: string]: unknown;
}

export interface Ticket {
  ticketid: TicketId;
  eventid?: TicketId;
  eventID?: TicketId;
  name?: string;
  price?: number;
  currency?: string;
  color?: string;
  quantity?: number;
  entity_id?: string;
  entity_type?: string;
  available?: number;
  total?: number;
  created_at?: DateLike;
  createdAt?: DateLike;
  updated_at?: DateLike;
  updatedAt?: DateLike;
  description?: string;
  sold?: number;
  seatstart?: number;
  seatend?: number;
  seats?: TicketSeat[];
  [key: string]: unknown;
}

export interface TicketPayload {
  name: string;
  price: number;
  quantity: number;
  currency: string;
  color: string;
  seatstart: number;
  seatend: number;
  [key: string]: unknown;
}

export interface PurchasedTicket {
  ticketid: TicketId;
  eventid?: TicketId;
  eventID?: TicketId;
  userid?: TicketId;
  userID?: TicketId;
  buyerName?: string;
  buyername?: string;
  uniquecode?: string;
  uniqueCode?: string;
  purchasedate?: DateLike;
  purchaseDate?: DateLike;
  status?: TicketStatus;
  canceled?: boolean;
  canceledAt?: DateLike;
  canceledat?: DateLike;
  canceledReason?: string;
  cancelledreason?: string;
  transferred?: boolean;
  transferredTo?: string;
  transferredto?: string;
  refundstatus?: string;
  [key: string]: unknown;
}

export type UserTicket = PurchasedTicket & {
  buyername?: string;
  purchasedate?: DateLike;
  transferredto?: string;
  refundstatus?: string;
  canceled?: boolean;
  status?: TicketStatus;
};
