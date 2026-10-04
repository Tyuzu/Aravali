export interface ItineraryVisit {
  start_time?: string;
  startTime?: string;
  end_time?: string;
  endTime?: string;
  location?: string;
  transport?: string;
  [key: string]: unknown;
}

export interface ItineraryDay {
  date?: string;
  visits?: ItineraryVisit[];
  [key: string]: unknown;
}

export interface ItineraryItem {
  itineraryid?: string | number;
  id?: string | number;
  userid?: string | number;
  userId?: string | number;
  name?: string;
  description?: string;
  start_date?: string;
  startDate?: string;
  end_date?: string;
  endDate?: string;
  status?: string;
  published?: boolean;
  forked_from?: string | null;
  forkedFrom?: string | null;
  days?: ItineraryDay[];
  [key: string]: unknown;
}

export interface ItineraryApiResponse<T = unknown> {
  success?: boolean;
  data?: T;
  message?: string;
  [key: string]: unknown;
}
