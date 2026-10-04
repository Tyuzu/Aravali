export interface MenuItem {
  menuid: string | number;
  placeid?: string | number;
  placeId?: string | number;
  name: string;
  price: number;
  discount?: number;
  stock: number;
  menu_pic?: string;
  menuPic?: string;
  description?: string;
  created_at?: string | number | Date;
  createdAt?: string | number | Date;
  updated_at?: string | number | Date;
  updatedAt?: string | number | Date;
  userid?: string | number;
  userId?: string | number;
  [key: string]: unknown;
}

export interface MenuApiResponse<T = unknown> {
  success?: boolean;
  data?: T;
  message?: string;
  [key: string]: unknown;
}

export interface StockResponse {
  stock: number;
  [key: string]: unknown;
}
