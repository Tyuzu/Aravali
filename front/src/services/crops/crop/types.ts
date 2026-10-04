import type { AvailabilitySchedule, CropListing as SharedCropListing } from "../types.js";

export type AvailabilityMap = AvailabilitySchedule;
export type CropListing = SharedCropListing;

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
