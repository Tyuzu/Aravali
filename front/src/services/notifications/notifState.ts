import { getState, setState } from "../../state/state.js";

export function getUserId(): string | null {
  try {
    const userStr = localStorage.getItem("user");
    if (!userStr) return null;
    const user = JSON.parse(userStr);
    return user.id || user._id || userStr;
  } catch {
    return localStorage.getItem("user");
  }
}

export function syncUnreadNotificationState(count: number): void {
  const nextCount = Math.max(0, Number.isFinite(count) ? Number(count) : 0);
  setState("unreadNotifications", nextCount, true);
}

export function decrementUnreadNotificationState(): void {
  const currentCount = Number(getState("unreadNotifications") || 0);
  syncUnreadNotificationState(Math.max(0, currentCount - 1));
}
