import "../../../css/inistyles/analytics.css";
import { createElement } from "../../components/createElement.js";
import { getAnalytics } from "./api.js";
import Notify from "../../components/ui/Notify.js";
import Datex from "../../components/base/Datex.js";

/* ────────────────────── Interfaces & Types ────────────────────── */

export interface EngagementMetrics {
  [key: string]: string | number;
}

export interface InsightsData {
  [key: string]: string | number;
}

export interface AnalyticsMetrics {
  [key: string]: number | string;
}

export interface AnalyticsData {
  name?: string;
  type?: string;
  lastUpdated?: string | number | Date;
  metrics?: AnalyticsMetrics;
  trend?: (number | string)[];
  engagement?: EngagementMetrics;
  insights?: InsightsData;
  topLocations?: (string | number)[];
}

export interface RenderAnalyticsPageParams {
  container: HTMLElement | null;
  isLoggedIn: boolean;
  entityType?: string;
  entityId?: string | number | null;
}

/* ────────────────────── Utility Helpers ────────────────────── */

/**
 * Format system keys (e.g. "page_views" or "bounceRate") into human-friendly titles ("Page Views", "Bounce Rate")
 */
function beautifyKey(key: string): string {
  if (!key) return "";
  return key
    .replace(/_/g, " ")
    .replace(/([A-Z])/g, " $1")
    .trim()
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

/* ────────────────────── Render Sub-components ────────────────────── */

function formatMetricValue(value: number | string): string {
  if (typeof value === "number") {
    if (Number.isInteger(value)) return String(value);
    return value.toFixed(1);
  }

  return String(value ?? "0");
}

function getTrendLabels(): string[] {
  const labels = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
  return labels;
}

function renderAnalytics(data: AnalyticsData): HTMLElement {
  const metrics = data.metrics && typeof data.metrics === "object" ? data.metrics : {};
  const trend = Array.isArray(data.trend) ? data.trend : [];
  const engagement = data.engagement && typeof data.engagement === "object" ? data.engagement : {};
  const insights = data.insights && typeof data.insights === "object" ? data.insights : {};
  const topLocations = Array.isArray(data.topLocations) ? data.topLocations : [];
  const lastUpdated = data.lastUpdated ? Datex(data.lastUpdated) : "";

  const metricEntries = Object.entries(metrics).slice(0, 8);
  const numericTrend = trend.map((v) => Math.max(0, Number(v) || 0));
  const maxTrendVal = Math.max(...numericTrend, 1);

  const headerChildren: (HTMLElement | null)[] = [
    createElement("span", { class: "analytics-kicker" }, [data.type ? data.type.toUpperCase() : "ENTITY"]),
    createElement("h2", {}, [`Analytics for ${data.name || "Unknown"}`]),
    createElement("p", {}, [`Entity Type: ${data.type || "N/A"}`]),
    lastUpdated ? createElement("small", {}, ["Last Updated: ", String(lastUpdated)]) : null
  ];

  const header = createElement(
    "div",
    { class: "analytics-header" },
    headerChildren.filter((child): child is HTMLElement => child !== null)
  );

  const summaryCards = createElement(
    "div",
    { class: "analytics-summary-cards" },
    metricEntries.map(([key, value]) =>
      createElement("div", { class: "analytics-card" }, [
        createElement("span", { class: "analytics-card-label" }, [beautifyKey(key)]),
        createElement("strong", { class: "analytics-card-value" }, [formatMetricValue(value)])
      ])
    )
  );

  const trendBars = createElement(
    "div",
    { class: "analytics-trend-bars" },
    numericTrend.map((val, index) => {
      const heightPct = (val / maxTrendVal) * 100;
      return createElement("div", { class: "analytics-trend-column" }, [
        createElement("span", { class: "analytics-trend-value" }, [formatMetricValue(val)]),
        createElement("span", {
          class: "analytics-trend-bar",
          style: { height: `${Math.max(12, heightPct)}%` },
          title: `Value: ${val}`
        }),
        createElement("small", { class: "analytics-trend-label" }, [getTrendLabels()[index] || `D${index + 1}`])
      ]);
    })
  );

  const trendSection = createElement("div", { class: "analytics-panel analytics-trend-panel" }, [
    createElement("h3", {}, ["7-Day Trend"]),
    trendBars
  ]);

  const engagementSection = Object.keys(engagement).length
    ? createElement("div", { class: "analytics-panel analytics-stats-panel" }, [
        createElement("h3", {}, ["Engagement"]),
        createElement(
          "div",
          { class: "analytics-stat-grid" },
          Object.entries(engagement).map(([k, v]) =>
            createElement("div", { class: "analytics-stat" }, [
              createElement("span", { class: "analytics-stat-label" }, [beautifyKey(k)]),
              createElement("strong", { class: "analytics-stat-value" }, [String(v)])
            ])
          )
        )
      ])
    : null;

  const insightsSection = Object.keys(insights).length
    ? createElement("div", { class: "analytics-panel analytics-insights-panel" }, [
        createElement("h3", {}, ["Insights"]),
        createElement(
          "div",
          { class: "analytics-insight-list" },
          Object.entries(insights).map(([k, v]) =>
            createElement("div", { class: "analytics-insight-item" }, [
              createElement("span", { class: "analytics-insight-kicker" }, [beautifyKey(k)]),
              createElement("p", {}, [String(v)])
            ])
          )
        )
      ])
    : null;

  const topLocationsSection = topLocations.length
    ? createElement("div", { class: "analytics-panel analytics-locations-panel" }, [
        createElement("h3", {}, ["Top Locations"]),
        createElement(
          "div",
          { class: "analytics-location-list" },
          topLocations.map((loc, index) =>
            createElement("div", { class: "analytics-location-chip" }, [`${index + 1}. ${String(loc)}`])
          )
        )
      ])
    : null;

  const sections: (HTMLElement | null)[] = [
    header,
    summaryCards,
    trendSection,
    engagementSection,
    insightsSection,
    topLocationsSection
  ];

  return createElement(
    "div",
    { class: "analytics-page" },
    sections.filter((section): section is HTMLElement => section !== null)
  );
}

/* ────────────────────── Main Entry ────────────────────── */

export async function renderAnalyticsPage({
  container,
  isLoggedIn,
  entityType = "events",
  entityId = null
}: RenderAnalyticsPageParams): Promise<void> {
  if (!container) return;

  // Clear previous DOM state safely
  container.replaceChildren();

  if (!isLoggedIn) {
    Notify("Please log in to view analytics.", { type: "warning", duration: 3000, dismissible: true });
    return;
  }

  try {
    const data = (await getAnalytics(entityType, entityId)) as AnalyticsData | null;

    if (!data || !data.metrics) {
      Notify("No analytics data found.", { type: "warning", dismissible: true });
      return;
    }

    container.appendChild(renderAnalytics(data));
  } catch (err) {
    Notify("Failed to load analytics data.", { type: "error", dismissible: true });
    console.error("Analytics rendering error:", err);
  }
}