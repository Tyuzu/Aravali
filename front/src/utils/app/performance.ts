import type { ErrorTracker } from "./errors.js";

/* =========================================================
   GLOBAL TYPE DECLARATIONS
========================================================= */
type PerformanceMetricDetails = Record<string, unknown> & {
  name: string;
  entryType: string;
  value: number;
  duration: number; // Retained for backward compatibility
  threshold: number;
};

declare global {
  interface Window {
    __errorTracker?: ErrorTracker;
  }
}

/* =========================================================
   PERFORMANCE MONITORING
========================================================= */
export function setupPerformanceMonitoring(): void {
  if (typeof window === "undefined" || !("PerformanceObserver" in window)) return;

  try {
    const supportedTypes = PerformanceObserver.supportedEntryTypes || [];
    const targetTypes = ["navigation", "longtask", "largest-contentful-paint"] as const;

    // Filter only the types supported by the current browser
    const validTypes = targetTypes.filter((type) => supportedTypes.includes(type));
    if (validTypes.length === 0) return;

    // Instantiate a single observer to handle all supported performance metrics
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        processPerformanceEntry(entry);
      }
    });

    // Observe each valid type with buffering enabled
    for (const type of validTypes) {
      observer.observe({ type, buffered: true });
    }
  } catch (error: unknown) {
    console.warn("[PERF] Monitoring initialization failed:", error);
  }
}

function shouldIgnorePerformanceEntry(entry: PerformanceEntry): boolean {
  const name = String(entry.name || "").toLowerCase();
  return /placehold\.co|doubleclick|googlesyndication|ads|adserver|googletagmanager/i.test(name);
}

export function processPerformanceEntry(entry: PerformanceEntry): void {
  if (shouldIgnorePerformanceEntry(entry)) {
    return;
  }

  switch (entry.entryType) {
    case "longtask":
      if (entry.duration > 200) {
        reportPerformanceIssue(entry, entry.duration, 200);
      }
      break;

    case "largest-contentful-paint":
      // Note: LCP entry.duration is always 0; startTime represents the render time.
      if (entry.startTime > 3000) {
        reportPerformanceIssue(entry, entry.startTime, 3000);
      }
      break;

    case "navigation":
      if (entry.duration > 3000) {
        reportPerformanceIssue(entry, entry.duration, 3000);
      }
      break;
  }
}

export function reportPerformanceIssue(
  entry: PerformanceEntry,
  metricValue: number,
  threshold: number
): void {
  const roundedValue = Math.round(metricValue);

  const details: PerformanceMetricDetails = {
    name: entry.name || "unknown",
    entryType: entry.entryType,
    value: roundedValue,
    duration: roundedValue,
    threshold
  };

  console.warn(`[PERF] Slow ${entry.entryType}:`, details);

  if (window.__errorTracker?.trackMetric) {
    try {
      window.__errorTracker.trackMetric("performance_degradation", details);
    } catch (error: unknown) {
      console.warn("[PERF] Metric tracking failed:", error);
    }
  }
}