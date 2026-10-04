export type ActivityEventType =
  | "pageview"
  | "click"
  | "scroll"
  | "input_focus"
  | "time_on_page"
  | "button_click"
  | "purchase"
  | string;

export interface ActivityItem {
  placeId?: string | number;
  placeid?: string | number;
  action?: string;
  performedBy?: string;
  performedby?: string;
  timestamp?: string | number | Date;
  details?: string;
  ipAddress?: string;
  ipaddress?: string;
  deviceInfo?: string;
  deviceinfo?: string;
  activityid?: string | number;
  activityId?: string | number;
  userid?: string | number;
  userId?: string | number;
  activity_type?: string;
  activityType?: string;
  entity_id?: string | number;
  entityId?: string | number;
  entity_type?: string | null;
  entityType?: string | null;
  [key: string]: unknown;
}

export interface ActivityAnalyticsEvent {
  type: ActivityEventType;
  data?: Record<string, unknown>;
  ts?: number;
}

export interface EnqueuedActivityEvent extends ActivityAnalyticsEvent {
  ts: number;
}

export interface ActivityBatchMetadata {
  lang: string;
  platform: string;
  referrer: string;
  url: string;
  ua: string;
  screen: string;
  session: string;
  user: string;
}

export interface ActivityBatchPayload {
  meta: ActivityBatchMetadata;
  events: EnqueuedActivityEvent[];
}
