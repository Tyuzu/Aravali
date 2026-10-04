export interface BookingItem {
  id?: string;
  bookingid?: string;
  userid?: string;
  userId?: string;
  entityType?: string;
  entitytype?: string;
  entityId?: string;
  entityid?: string;
  slotId?: string;
  slotid?: string;
  tierId?: string;
  tierid?: string;
  tierName?: string;
  tiername?: string;
  date?: string;
  start?: string;
  end?: string;
  seats?: number;
  pricePaid?: number;
  price_paid?: number;
  status?: string;
  createdAt?: string | number;
  created_at?: string | number;
  [key: string]: unknown;
}

export interface PricingTier {
  id?: string;
  entityType?: string;
  entitytype?: string;
  entityId?: string;
  entityid?: string;
  name?: string;
  price?: number;
  capacity?: number;
  timeRange?: [string, string];
  time_range?: [string, string];
  daysOfWeek?: number[];
  days_of_week?: number[];
  features?: string[];
  createdAt?: string | number;
  created_at?: string | number;
  [key: string]: unknown;
}

export interface BookingSlot {
  id?: string;
  entityType?: string;
  entitytype?: string;
  entityId?: string;
  entityid?: string;
  date?: string;
  start?: string;
  end?: string;
  capacity?: number;
  tierId?: string;
  tierid?: string;
  tierName?: string;
  tiername?: string;
  createdAt?: string | number;
  created_at?: string | number;
  [key: string]: unknown;
}

export interface CreateBookingPayload {
  seats?: number | string;
  userid?: string;
  userId?: string;
  date?: string;
  start?: string;
  end?: string;
  slotId?: string;
  slotid?: string;
  tierId?: string;
  tierid?: string;
  pricePaid?: number;
  price_paid?: number;
  entityType?: string;
  entitytype?: string;
  entityId?: string;
  entityid?: string;
  [key: string]: unknown;
}
