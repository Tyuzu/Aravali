export interface PlaceCoordinates {
  latitude?: number;
  longitude?: number;
  lat?: number;
  lng?: number;
  [key: string]: unknown;
}

export interface Place {
  placeid: string | number;
  placeId?: string | number;
  createdBy?: string | number;
  creatorid?: string | number;
  name?: string;
  short_desc?: string;
  description?: string;
  place?: string;
  capacity?: number | string;
  date?: string | Date | null;
  address?: string;
  category?: string;
  banner?: string;
  website_url?: string;
  website?: string;
  status?: string;
  accessibility_info?: string;
  social_media_links?: string[];
  socialLinks?: Record<string, string>;
  tags?: string[];
  custom_fields?: Record<string, unknown>;
  created_at?: string | Date | null;
  updated_at?: string | Date | null;
  city?: string;
  country?: string;
  zipCode?: string;
  jobs?: string;
  location?: PlaceCoordinates | Record<string, unknown>;
  coordinates?: PlaceCoordinates;
  phone?: string;
  isopen?: boolean;
  distance?: number;
  views?: number;
  reviewcount?: number;
  updatedBy?: string;
  deletedAt?: string | Date | null;
  amenities?: string[];
  events?: string[];
  operatinghours?: string[];
  keywords?: string[];
  [key: string]: unknown;
}
