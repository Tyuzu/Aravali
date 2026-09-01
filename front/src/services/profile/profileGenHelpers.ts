import { formatDate } from "./profileHelpers.js";

/* ============================================================
    TYPE DEFINITIONS
============================================================ */

export interface UserProfile {
    userid?: string | number;
    username?: string;
    name?: string;
    bio?: string;
    is_following?: boolean;
    last_login?: string | Date;
    is_verified?: boolean;
    wallet_balance?: number;
    followerscount?: number;
    followscount?: number;
    [key: string]: unknown;
}

export interface InfoItem {
    label: string;
    value: string | Node;
}

export interface StatItem {
    label: string;
    value: number | string;
}

/* ============================================================
    DOM HELPERS & COMPONENTS
============================================================ */

/**
 * Appends child nodes to a parent element safely.
 */
import { createProfileDetails } from "./components/createProfileDetails.js";
import { createStatistics } from "./components/createStatistics.js";
import { createProfileActions } from "./components/createProfileActions.js";
import { createProfileInfo } from "./components/createProfileInfo.js";

function appendChildren(parent: HTMLElement, ...children: (Node | null | undefined)[]): void {
    children.forEach(child => {
        if (child instanceof Node) {
            parent.appendChild(child);
        } else if (child !== null && child !== undefined) {
            console.error("Invalid child passed to appendChildren:", child);
        }
    });
}

/* ============================================================
    FORMATTERS & HELPERS
============================================================ */

/**
 * Formats a date string into a Datex component instance or formatted string
 */
export function formatDateUtil(dateString?: string | Date | null): HTMLElement | string | null {
    if (!dateString) return null;
    try {
        const safe = dateString instanceof Date ? dateString.toISOString() : (dateString as string);
        return formatDate(safe);
    } catch (error) {
        console.error("Error formatting date with Datex:", error);
        return new Date(dateString).toLocaleString();
    }
}

/**
 * Capitalizes the first letter of a string
 */
export function capitalize(string: string = ""): string {
    if (!string) return "";
    return string.charAt(0).toUpperCase() + string.slice(1);
}

/* ============================================================
    LOADING INDICATORS
============================================================ */

/**
 * Renders a loading message to a target container or default #content element
 */
export function showLoadingMessage(message: string, containerId: string = "content"): void {
    removeLoadingMessage();

    const container = document.getElementById(containerId);
    if (!container) {
        console.warn(`Container #${containerId} not found to show loading message.`);
        return;
    }

    const loadingMsg = document.createElement("p");
    loadingMsg.id = "loading-msg";
    loadingMsg.className = "loading-message";
    loadingMsg.textContent = message;

    container.appendChild(loadingMsg);
}

/**
 * Removes active loading message from the DOM
 */
export function removeLoadingMessage(): void {
    const loadingMsg = document.getElementById("loading-msg");
    if (loadingMsg) {
        loadingMsg.remove();
    }
}

/* ============================================================
    MEDIA PREVIEWS
============================================================ */

/**
 * Previews an image file selection on a target image element
 */
export function previewAvatar(event: Event, previewId: string = "profile-picture-preview"): void {
    const target = event.target as HTMLInputElement | null;
    const file = target?.files?.[0];
    const preview = document.getElementById(previewId) as HTMLImageElement | null;

    if (!preview) return;

    if (file) {
        // Revoke previous Object URL to prevent memory leaks if re-uploading
        if (preview.dataset.objectUrl) {
            URL.revokeObjectURL(preview.dataset.objectUrl);
        }

        const objectUrl = URL.createObjectURL(file);
        preview.src = objectUrl;
        preview.style.display = "block";
        preview.dataset.objectUrl = objectUrl;
    }
}

export {
    createProfileDetails,
    createProfileActions,
    createProfileInfo,
    createStatistics,
    appendChildren
};