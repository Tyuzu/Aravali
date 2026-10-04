import { apiFetch } from "../../api/api.js";

export interface ToggleLikeResponse {
    liked: boolean;
    count: number;
    userid?: string | number;
    userId?: string | number;
    activity_type?: string;
    activityType?: string;
    entity_id?: string | number;
    entityId?: string | number;
    entity_type?: string;
    entityType?: string;
    [key: string]: unknown;
}

export interface FollowEntityResponse {
    ok?: boolean;
    success?: boolean;
    userid?: string | number;
    userId?: string | number;
    follows?: Array<string | number>;
    followers?: Array<string | number>;
    [key: string]: unknown;
}

export async function toggleLike(entityType: string, entityId: string | number): Promise<ToggleLikeResponse> {
    try {
        const path = `/likes/${entityType}/like/${entityId}`;
        return (await apiFetch(path, "PUT")) as ToggleLikeResponse;
    } catch (err) {
        console.error("beats.api.toggleLike error:", err);
        return { liked: false, count: 0 } as ToggleLikeResponse;
    }
}

export async function followEntity(apiPath: string, entityId: string | number, method: "PUT" | "DELETE"): Promise<FollowEntityResponse> {
    try {
        const endpoint = `${apiPath}${entityId}`;
        return (await apiFetch(endpoint, method)) as FollowEntityResponse;
    } catch (err) {
        console.error("beats.api.followEntity error:", err);
        return { ok: false };
    }
}

export default {
    toggleLike,
    followEntity,
};
