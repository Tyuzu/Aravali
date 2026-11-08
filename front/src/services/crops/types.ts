export interface PriceHistory {
  date?: string | number | Date;
  price?: number;
  [key: string]: unknown;
}

export interface AvailabilityDay {
  enabled?: boolean;
  from?: string;
  to?: string;
  [key: string]: unknown;
}

export type AvailabilitySchedule = Record<string, AvailabilityDay>;

export interface Farm {
  farmid?: string | number;
  id?: string | number;
  name?: string;
  location?: string;
  description?: string;
  owner?: string;
  createdBy?: string | number;
  createdby?: string | number;
  contact?: string;
  phone?: string;
  practice?: string;
  social?: string;
  banner?: string;
  photo?: string;
  media?: string[];
  availability?: AvailabilitySchedule;
  crops?: Crop[];
  updatedAt?: string | number | Date;
  createdAt?: string | number | Date;
  tags?: string[];
  avgRating?: number;
  reviewCount?: number;
  favoritesCount?: number;
  [key: string]: unknown;
}

export interface Crop {
  cropid?: string | number;
  id?: string | number;
  farmid?: string | number;
  farmId?: string | number;
  name?: string;
  category?: string;
  price?: number;
  discount?: number;
  unit?: string;
  quantity?: number;
  availableQtyKg?: number;
  banner?: string;
  notes?: string;
  tags?: string[];
  harvestDate?: string | number | Date;
  HarvestDate?: string | number | Date;
  plantedDate?: string | number | Date;
  lastSoldAt?: string | number | Date;
  expiryDate?: string | number | Date;
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  createdby?: string;
  createdBy?: string;
  farmName?: string;
  farmname?: string;
  outOfStock?: boolean;
  featured?: boolean;
  history?: PriceHistory[];
  priceHistory?: PriceHistory[];
  [key: string]: unknown;
}

export interface CropListing {
  cropid?: string | number;
  farmid?: string | number;
  farmName?: string;
  farmname?: string;
  breed?: string;
  banner?: string;
  location?: string;
  pricePerKg?: number;
  unit?: string;
  availableQtyKg?: number;
  inventoryValue?: number;
  outOfStock?: boolean;
  featured?: boolean;
  avgRating?: number;
  reviewCount?: number;
  favoritesCount?: number;
  harvestDate?: string | number | Date;
  plantedDate?: string | number | Date;
  lastSoldAt?: string | number | Date;
  availability?: AvailabilitySchedule;
  phone?: string;
  tags?: string[];
  [key: string]: unknown;
}
