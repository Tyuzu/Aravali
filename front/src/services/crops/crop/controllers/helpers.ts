import type { AvailabilityMap } from "../types";

export function formatAvailability(availability?: AvailabilityMap): string {
  if (!availability || typeof availability !== "object") {
    return "N/A";
  }

  const activeDays = Object.entries(availability)
    .filter(([_, value]) => value && value.enabled)
    .map(([day, value]) => {
      const capitalized = day.charAt(0).toUpperCase() + day.slice(1);
      return `${capitalized}: ${value.from || ""}-${value.to || ""}`;
    });

  return activeDays.length > 0 ? activeDays.join(", ") : "Closed";
}

export function formatRelativeDate(dateString?: string): string {
  if (!dateString) return "N/A";

  const date = new Date(dateString);
  if (isNaN(date.getTime())) return "N/A";

  const diffDays = Math.floor((Date.now() - date.getTime()) / (1000 * 60 * 60 * 24));

  if (diffDays <= 0) return "Today";
  if (diffDays === 1) return "1 day ago";
  return `${diffDays} days ago`;
}

export function getStockStatus(qty: number): string {
  if (qty <= 0) return "Out of Stock";
  if (qty < 10) return "Low";
  return "In Stock";
}
