export interface CartItem {
  id?: string;
  itemId?: string;
  itemid?: string;
  itemType?: string;
  itemtype?: string;
  entityId?: string | number;
  entityid?: string | number;
  entityType?: string;
  entitytype?: string;
  itemName?: string;
  itemname?: string;
  category?: string;
  quantity?: number;
  price?: number;
  discount?: number;
  unit?: string;
  addedAt?: string | number;
  added_at?: string | number;
  updatedAt?: string | number;
  updated_at?: string | number;
  userid?: string;
  userId?: string;
  [key: string]: unknown;
}

export interface CartOrder {
  orderId?: string;
  orderid?: string;
  orderType?: string;
  ordertype?: string;
  userid?: string;
  userId?: string;
  items?: Record<string, CartItem[]> | Record<string, unknown>;
  address?: string;
  paymentMethod?: string;
  paymentmethod?: string;
  status?: string;
  createdAt?: string | number;
  created_at?: string | number;
  subtotal?: number;
  discount?: number;
  tax?: number;
  delivery?: number;
  total?: number;
  name?: string;
  phone?: string;
  [key: string]: unknown;
}

export interface CheckoutSession {
  userid?: string;
  userId?: string;
  items?: Record<string, CartItem[]>;
  address?: string;
  total?: number;
  subtotal?: number;
  tax?: number;
  delivery?: number;
  discount?: number;
  paymentMethod?: string;
  paymentmethod?: string;
  paymentDetails?: unknown;
  createdAt?: string | number;
  created_at?: string | number;
  [key: string]: unknown;
}
