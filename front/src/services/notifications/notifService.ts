import { apiFetch } from "../../api/api.js";

export interface NotificationItem {
  id?: string | number;
  notificationid?: string | number;
  title?: string;
  type?: string;
  message?: string;
  createdAt?: string | number | Date;
  isRead?: boolean;
}

export async function getNotifications(): Promise<NotificationItem[]> {
  try {
    const response: any = await apiFetch("/notifs", "GET");
    const rawList: NotificationItem[] = Array.isArray(response)
      ? response
      : response?.notifications || response?.data || [];

    return rawList.sort((a, b) => new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime());
  } catch (error) {
    console.error("Failed to fetch notifications:", error);
    return [];
  }
}

export async function markNotificationAsRead(id: string | number | undefined): Promise<unknown> {
  if (!id) throw new Error("Notification ID is required.");
  try {
    return await apiFetch(`/notifs/notif/${id}/read`, "PUT");
  } catch (error) {
    console.error(`Failed to mark notification ${id} as read:`, error);
    throw error;
  }
}

export async function markAllNotificationsAsRead(): Promise<unknown> {
  try {
    return await apiFetch("/notifs/read-all", "PUT");
  } catch (error) {
    console.error("Failed to mark all notifications as read:", error);
    throw error;
  }
}

export async function clearAllNotifications(): Promise<unknown> {
  try {
    return await apiFetch("/notifs", "DELETE");
  } catch (error) {
    console.error("Failed to clear notifications:", error);
    throw error;
  }
}