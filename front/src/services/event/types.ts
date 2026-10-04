export interface EventCoordinates {
  latitude?: number;
  longitude?: number;
  lat?: number;
  lng?: number;
  [key: string]: unknown;
}

export interface EventContactInfo {
  email?: string;
  phone?: string;
  organizer_name?: string;
  organizerName?: string;
  [key: string]: unknown;
}

export interface EventNewsItem {
  id?: string | number;
  title?: string;
  content?: string;
  timestamp?: string | Date;
  [key: string]: unknown;
}

export interface EventPollOption {
  text?: string;
  votes?: number;
  [key: string]: unknown;
}

export interface EventPoll {
  id?: string | number;
  question?: string;
  options?: EventPollOption[];
  [key: string]: unknown;
}

export interface EventLostFoundItem {
  id?: string | number;
  type?: string;
  description?: string;
  contact?: string;
  [key: string]: unknown;
}

export interface EventItem {
  id?: string | number;
  eventid: string | number;
  title?: string;
  name?: string;
  description?: string;
  date: string | Date;
  placeid?: string | number;
  placename?: string;
  location?: string;
  coords?: EventCoordinates;
  creatorid?: string | number;
  category?: string;
  banner?: string;
  seating?: unknown;
  website_url?: string;
  status?: string;
  tags?: string[];
  created_at?: string | Date | null;
  updated_at?: string | Date | null;
  organizer_name?: string;
  organizer_contact?: string;
  artists?: string[];
  published?: string | boolean;
  external?: boolean;
  externallink?: string;
  contactInfo?: EventContactInfo;
  news?: EventNewsItem[];
  polls?: EventPoll[];
  lostfound?: EventLostFoundItem[];
  prices?: number[];
  currency?: string;
  [key: string]: unknown;
}

export interface EventDetailData extends EventItem {
  eventid: string | number;
  social_links?: Record<string, string>;
  custom_fields?: Record<string, string | number>;
}

export interface EventData extends EventItem {
  success?: boolean;
  error?: string;
}
