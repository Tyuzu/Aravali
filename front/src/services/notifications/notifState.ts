import { getState, setState } from "../../state/state.js";

export function syncUnreadNotificationState(count: number): void {
  const numericCount = Number(count);
  const nextCount = Math.max(0, Number.isFinite(numericCount) ? numericCount : 0);

  setState("unreadNotifications", nextCount, true);
}

export function decrementUnreadNotificationState(): void {
  const currentCount = Number(getState("unreadNotifications") || 0);
  const safeCount = Number.isFinite(currentCount) ? currentCount : 0;

  syncUnreadNotificationState(safeCount - 1);
}

export function incrementUnreadNotificationState(step: number = 1): void {
  const currentCount = Number(getState("unreadNotifications") || 0);
  const safeCount = Number.isFinite(currentCount) ? currentCount : 0;

  syncUnreadNotificationState(safeCount + step);
}