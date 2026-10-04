export interface AnalyticsMetricMap {
  [key: string]: number | string;
}

export interface AnalyticsItem {
  id?: string;
  name?: string;
  type?: string;
  metrics?: Record<string, number | string>;
  trend?: Array<number | string>;
  topLocations?: Array<string | number>;
  engagement?: Record<string, string | number>;
  insights?: Record<string, string | number>;
  lastUpdated?: string | number | Date;
  last_updated?: string | number | Date;
  [key: string]: unknown;
}

export type AnalyticsData = AnalyticsItem;
