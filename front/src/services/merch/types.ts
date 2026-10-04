export interface MerchItem {
  merchid: string | number;
  name: string;
  slug?: string;
  sku?: string;
  category?: string;
  price: number;
  discount?: number;
  stock: number | string;
  stock_status?: string;
  stockStatus?: string;
  merch_pic?: string;
  merchPic?: string;
  gallery?: string[];
  entity_id?: string | number;
  entityId?: string | number;
  entity_type?: string;
  entityType?: string;
  description?: string;
  short_desc?: string;
  shortDesc?: string;
  rating?: number;
  review_count?: number;
  reviewCount?: number;
  weight?: number;
  dimensions?: string;
  tags?: string[];
  created_at?: string | number | Date;
  createdAt?: string | number | Date;
  updated_at?: string | number | Date;
  updatedAt?: string | number | Date;
  deleted_at?: string | number | Date | null;
  userid?: string | number;
  userId?: string | number;
  [key: string]: unknown;
}

export interface MerchApiResponse<T = unknown> {
  success?: boolean;
  data?: T;
  message?: string;
  [key: string]: unknown;
}
