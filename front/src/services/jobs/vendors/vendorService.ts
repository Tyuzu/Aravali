import { apiFetch } from "../../../api/api.js";

export type VendorRecord = {
    vendorid?: string | number;
    vendorId?: string | number;
    vendor_id?: string | number;
    userid?: string | number;
    userId?: string | number;
    name?: string;
    full_name?: string;
    business_name?: string;
    category?: string;
    description?: string | null;
    email?: string | null;
    phone?: string | null;
    location?: string | null;
    rating?: number | string | null;
    rating_count?: number | string | null;
    ratingCount?: number | string | null;
    profile_image?: string | null;
    profileImage?: string | null;
    portfolio?: string[];
    verified?: boolean;
    available?: boolean;
    created_at?: string | number | Date | null;
    createdAt?: string | number | Date | null;
    updated_at?: string | number | Date | null;
    updatedAt?: string | number | Date | null;
    hiringid?: string | number;
    hiringId?: string | number;
    eventid?: string | number;
    eventId?: string | number;
    hired_at?: string | number | Date | null;
    hiredAt?: string | number | Date | null;
    status?: string;
    [key: string]: unknown;
};

function requireId(value: unknown, label: string): void {
    if (value === null || value === undefined || value === "") {
        throw new Error(`${label} is required.`);
    }
}

export async function fetchVendors(): Promise<VendorRecord[]> {
    return apiFetch<VendorRecord[]>('/vendors', 'GET');
}

export async function fetchVendor(vendorId: string): Promise<any> {
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/${vendorId}`, "GET");
}

export async function fetchEventVendors(eventId: string): Promise<any> {
    requireId(eventId, "Event ID");
    return apiFetch(`/vendors/events/${eventId}`, "GET");
}

export async function createVendor(payload: Record<string, any>): Promise<any> {
    return apiFetch("/vendors", "POST", payload);
}

export async function updateVendor(vendorId: string, payload: Record<string, any>): Promise<any> {
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/${vendorId}`, "PATCH", payload);
}

export async function deleteVendor(vendorId: string): Promise<any> {
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/${vendorId}`, "DELETE");
}

export async function hireEventVendor(eventId: string, vendorId: string): Promise<any> {
    requireId(eventId, "Event ID");
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/events/${eventId}/hire`, "POST", {
        vendorid: vendorId
    });
}

export async function removeEventVendor(eventId: string, vendorId: string): Promise<any> {
    requireId(eventId, "Event ID");
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/events/${eventId}/vendor/${vendorId}`, "DELETE");
}

export async function fetchAvailability(vendorId: string): Promise<any> {
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/${vendorId}/availability`, "GET");
}

export async function getMyVendorRequests(): Promise<any> {
    return apiFetch(`/vendors/me/requests`, "GET");
}

export async function updateVendorHiringStatus(hiringId: string, status: string): Promise<any> {
    requireId(hiringId, "Hiring ID");
    requireId(status, "Status");
    return apiFetch(`/vendors/hiring/${hiringId}/status`, "PATCH", {
        status,
    });
}

export async function createAvailability(vendorId: string, slot: Record<string, any>): Promise<any> {
    requireId(vendorId, "Vendor ID");
    return apiFetch(`/vendors/${vendorId}/availability`, "POST", slot);
}

export async function deleteAvailability(vendorId: string, slotId: string): Promise<any> {
    requireId(vendorId, "Vendor ID");
    requireId(slotId, "Slot ID");
    return apiFetch(`/vendors/${vendorId}/availability/${slotId}`, "DELETE");
}