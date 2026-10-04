import { apiFetch } from "../../api/api.js";
import type { EntityItem, EntityType } from "./types.js";

function normalizeEntityItems(items: EntityItem[] | unknown): EntityItem[] {
  if (!Array.isArray(items)) return [];

  return items.map((item) => {
    const entityId = item.entity_id ?? item.entityId ?? item.postid ?? item.postId ?? item.id ?? item.item_id ?? item.itemId;
    const createdAt = item.created_at ?? item.createdAt ?? new Date(0);

    return {
      ...item,
      entity_id: entityId ?? item.entity_id,
      entityId: entityId ?? item.entityId,
      created_at: createdAt,
      createdAt: createdAt,
      userId: item.userId ?? item.userid,
      userid: item.userid ?? item.userId,
      item_id: item.item_id ?? item.itemId,
      itemId: item.itemId ?? item.item_id,
      image_url: item.image_url ?? item.imageUrl,
      imageUrl: item.imageUrl ?? item.image_url,
      entity_type: item.entity_type ?? item.entityType,
      entityType: item.entityType ?? item.entity_type,
      item_type: item.item_type ?? item.itemType,
      itemType: item.itemType ?? item.item_type,
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
