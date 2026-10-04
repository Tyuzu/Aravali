export type VendorLike = {
    vendorid?: string | number;
    vendorId?: string | number;
    vendor_id?: string | number;
    id?: string | number;
    name?: string;
    full_name?: string;
    fullName?: string;
    business_name?: string;
    businessName?: string;
    title?: string;
    profile_image?: string | null;
    profileImage?: string | null;
    rating_count?: number | string | null;
    ratingCount?: number | string | null;
    rating?: number | string | null;
    category?: string;
    [key: string]: unknown;
};

export function normalizeVendorList(response: any): any[] {
    if (Array.isArray(response)) {
        return response;
    }

    if (Array.isArray(response?.vendors)) {
        return response.vendors;
    }

    if (Array.isArray(response?.data)) {
        return response.data;
    }

    if (Array.isArray(response?.items)) {
        return response.items;
    }

    if (Array.isArray(response?.results)) {
        return response.results;
    }

    return [];
}

export function getVendorId(vendor: any): string | null {
    const value = vendor?.vendorid ?? vendor?.vendor_id ?? vendor?.vendorId ?? vendor?.id ?? null;
    return value === null || value === undefined ? null : String(value);
}

export function getVendorName(vendor: any): string {
    return (
        vendor?.name ??
        vendor?.full_name ??
        vendor?.fullName ??
        vendor?.fullname ??
        vendor?.business_name ??
        vendor?.businessName ??
        vendor?.title ??
        "Unnamed Vendor"
    );
}

export function getVendorProfileImage(vendor: VendorLike | null | undefined): string | null {
    return (vendor?.profile_image ?? vendor?.profileImage ?? null) as string | null;
}

export function getVendorRatingCount(vendor: VendorLike | null | undefined): number {
    const raw = vendor?.rating_count ?? vendor?.ratingCount ?? vendor?.rating ?? 0;
    const numeric = Number(raw);
    return Number.isFinite(numeric) ? numeric : 0;
}

export function normalizeErrorMessage(error: unknown): string {
    if (!error) {
        return "";
    }

    if (typeof error === "string") {
        return error;
    }

    if (error instanceof Error) {
        return error.message || "";
    }

    const errObj = error as any;
    return errObj.message || errObj.error || errObj.details || errObj.msg || "";
}

export function isValidEmail(email?: string): boolean {
    if (!email) {
        return true;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

export function formatRequestStatus(status?: string): string {
    if (!status) return "Pending";
    switch (String(status).toLowerCase()) {
        case "pending": return "Request Pending";
        case "accepted":
        case "hired": return "Already Hired ✓";
        case "completed": return "Completed";
        case "cancelled": return "Cancelled";
        default: return status;
    }
}