export type RefundStatus = "pending" | "approved" | "rejected" | "completed" | "none";

export type OrderType =
  | "farm"
  | "regular"
  | "product"
  | "merch"
  | "ticket"
  | "subscription"
  | "menu";

export interface RefundRequest {
  id: string;
  order_id?: string;
  orderId?: string;
  userid?: string;
  userId?: string;
  amount: number;
  reason?: string;
  status: RefundStatus;
  created_at?: string;
  createdAt?: string;
  order_type?: OrderType;
  orderType?: OrderType;
  review_notes?: string;
  reviewNotes?: string;
  reviewed_by?: string;
  reviewedBy?: string;
  [key: string]: unknown;
}

export interface OrderItem {
  entityName?: string;
  itemName?: string;
  quantity?: number;
  price?: number;
}

export interface OrderItemsStructure {
  products?: OrderItem[];
  [category: string]: OrderItem[] | undefined;
}

/**
 * Clean domain model for an Order.
 * Raw API payloads with duplicate/casing variants should be 
 * converted to this canonical shape via `normalizeOrders()`.
 */
export interface Order {
  orderId: string;
  orderid?: string;
  orderType: OrderType | string;
  ordertype?: string;
  createdAt: string | number;
  created_at?: string | number;
  status: string;
  paymentMethod: string;
  paymentmethod?: string;
  address: string;
  total: number;
  subtotal?: number;
  discount?: number;
  tax?: number;
  delivery?: number;
  approvedBy?: string[];
  approvedby?: string[];
  farmId?: string;
  farmid?: string;
  userid?: string;
  userId?: string;
  items: OrderItemsStructure;
  refundStatus?: RefundStatus;
  [key: string]: unknown;
}

export interface OrderFilters {
  status: string;
  date: string;
}

export interface OrderPageState {
  orders: Order[];
  filters: OrderFilters;
  currentPage: number;
  expandedOrders: Set<string>;
  loading?: boolean;
}

export type StateMutatedCallback = () => void;