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
  orderId: string;
  userId: string;
  amount: number;
  reason?: string;
  status: RefundStatus;
  createdAt?: string;
  orderType?: OrderType;
  reviewNotes?: string;
  reviewedBy?: string;
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
  orderType: OrderType | string;
  createdAt: string | number;
  status: string;
  paymentMethod: string;
  address: string;
  total: number;
  subtotal?: number;
  discount?: number;
  tax?: number;
  delivery?: number;
  approvedBy?: string[];
  farmId?: string;
  items: OrderItemsStructure;
  refundStatus?: RefundStatus;
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