import { apiFetch } from "../../api/api.js";
import type { EntityItem, EntityType } from "./types.js";

function normalizeEntityItems(items: EntityItem[] | unknown): EntityItem[] {
  if (!Array.isArray(items)) return [];

  return items.map((item) => {
    const entityId = item.entity_id ?? item.postid ?? item.id;
    return {
      ...item,
      entity_id: entityId ?? item.entity_id,
      created_at: item.created_at ?? new Date(0)
    };
  });
}

export async function fetchUserProfileData(
  username: string,
  entityType: EntityType
): Promise<EntityItem[]> {
  if (typeof username !== "string" || !username.trim()) {
    throw new Error("Username is required.");
  }

  if (typeof entityType !== "string" || !entityType.trim()) {
    throw new Error("Entity type is required.");
  }

  const encodedUsername = encodeURIComponent(username.trim());
  const encodedEntityType = encodeURIComponent(entityType.trim());

  try {
    const response = await apiFetch<EntityItem[]>(
      `/user/${encodedUsername}/data?entity_type=${encodedEntityType}`,
      "GET"
    );
    return normalizeEntityItems(response);
  } catch (error) {
    console.error(`Error fetching ${entityType} data for user:`, error);
    throw error;
  }
}

export async function fetchOtherUserProfileData(
  userId: string,
  entityType: EntityType
): Promise<EntityItem[]> {
  if (typeof userId !== "string" || !userId.trim()) {
    throw new Error("User ID is required.");
  }

  if (typeof entityType !== "string" || !entityType.trim()) {
    throw new Error("Entity type is required.");
  }

  const encodedUserId = encodeURIComponent(userId.trim());
  const encodedEntityType = encodeURIComponent(entityType.trim());

  try {
    const response = await apiFetch<EntityItem[]>(
      `/user/${encodedUserId}/udata?entity_type=${encodedEntityType}`,
      "GET"
    );
    return normalizeEntityItems(response);
  } catch (error) {
    console.error(`Error fetching other user's ${entityType} data:`, error);
    throw error;
  }
}

export default {
  fetchUserProfileData,
  fetchOtherUserProfileData
};
