export type ItemType = "product" | "tool";

export interface CategoryOption {
  value: string;
  label: string;
}

export interface FarmItem {
  productid?: string | number;
  productId?: string | number;
  userid?: string | number;
  userId?: string | number;
  name?: string;
  category?: string;
  price?: number | string;
  discount?: number | string;
  quantity?: number | string;
  unit?: string;
  sku?: string | null;
  availableFrom?: string | Date | null;
  availableTo?: string | Date | null;
  description?: string;
  featured?: boolean;
  banner?: string | null;
  photo?: string | null;
  images?: string | string[] | null;
  type?: string;
  seller?: {
    name?: string;
    contact?: string;
    [key: string]: unknown;
  } | null;
  createdAt?: string | number | Date | null;
  updatedAt?: string | number | Date | null;
  [key: string]: unknown;
}

export interface ItemPayload {
  name: string;
  category: string;
  price: number | string;
  discount?: number | string;
  quantity: number | string;
  unit?: string;
  sku?: string | null;
  availableFrom?: string | Date | null;
  availableTo?: string | Date | null;
  description?: string;
  featured?: boolean;
  type?: string;
  [key: string]: unknown;
}

export interface DisplayItemsOptions {
  limit?: number;
  offset?: number;
  search?: string;
  category?: string;
  sort?: string;
}

export interface UserState {
  userid?: string | number;
  userId?: string | number;
  [key: string]: unknown;
}