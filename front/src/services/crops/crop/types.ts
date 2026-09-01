export interface AvailabilityDay {
  enabled?: boolean;
  from?: string;
  to?: string;
}

export type AvailabilityMap = Record<string, AvailabilityDay>;

export interface CropListing {
  cropid: string;
  farmid: string;
  farmName?: string;
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
  harvestDate?: string;
  plantedDate?: string;
  lastSoldAt?: string;
  availability?: AvailabilityMap;
  phone?: string;
  tags?: string[];
}

export interface CropApiResponse {
  success: boolean;
  name?: string;
  category?: string;
  total?: number;
  listings?: CropListing[];
}

export interface FilterValues {
  location: string;
  breed: string;
  minPrice: number | null;
  maxPrice: number | null;
  minQty: number | null;
  maxQty: number | null;
  harvestDate: string | null;
}

export interface SetupFilterInteractionsParams {
  filterForm: HTMLFormElement;
  toggleFiltersBtn: HTMLElement;
  listings: CropListing[];
  onFiltered: (data: CropListing[]) => void;
}
