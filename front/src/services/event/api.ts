import { apiFetch } from "../../api/api.js";

export async function fetchEventById(eventId: string | number): Promise<any> {
  return await apiFetch(`/events/event/${eventId}`);
}

export async function fetchEventsPage(page: number = 1, limit: number = 1000): Promise<any> {
  return await apiFetch(`/events/events?page=${page}&limit=${limit}`);
}

export async function createEventRequest(formData: FormData): Promise<any> {
  return await apiFetch(`/events/event`, "POST", formData);
}

export async function updateEventRequest(eventId: string | number, formData: FormData): Promise<any> {
  return await apiFetch(`/events/event/${eventId}`, "PUT", formData);
}

export async function fetchEventMerch(eventId: string | number): Promise<any> {
  return await apiFetch(`/merch/event/${eventId}`);
}

export async function fetchFaqsForEvent(eventId: string | number): Promise<any[]> {
  const response = await apiFetch<any[] | { data?: any[] }>(`/faqs/event/${eventId}`);
  const items = Array.isArray(response) ? response : Array.isArray(response?.data) ? response.data : [];

  return items.map((faq: any) => {
    const raw = typeof faq?.content === "string" ? faq.content.trim() : "";
    const lines = raw ? raw.split(/\r?\n/).map((line: string) => line.trim()).filter(Boolean) : [];

    if (lines.length > 1) {
      return {
        ...faq,
        title: faq?.title || faq?.question || lines[0],
        content: lines.slice(1).join("\n") || lines[0]
      };
    }

    return {
      ...faq,
      title: faq?.title || faq?.question || "FAQ",
      content: raw || faq?.answer || ""
    };
  });
}

export async function createFaqForEvent(eventId: string | number, title: string, content: string): Promise<any> {
  const faqContent = title && content ? `${title}\n${content}` : (title || content || "");
  return await apiFetch(`/faqs/event/${eventId}`, "POST", { content: faqContent });
}

export async function fetchLostAndFoundForEvent(eventId: string | number): Promise<any> {
  return await apiFetch(`/events/${eventId}/lostfound`);
}

export async function createLostAndFoundItem(eventId: string | number, payload: Record<string, any>): Promise<any> {
  return await apiFetch(`/events/${eventId}/lostfound`, "POST", payload);
}

export default {
  fetchEventById,
  fetchEventsPage,
  createEventRequest,
  updateEventRequest,
  fetchEventMerch,
  fetchFaqsForEvent,
  createFaqForEvent,
  fetchLostAndFoundForEvent,
  createLostAndFoundItem
};
