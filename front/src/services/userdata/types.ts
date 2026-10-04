export type EntityType =
  | "userhome"
  | "place"
  | "event"
  | "feedpost"
  | "media"
  | "ticket"
  | "merch"
  | "review"
  | "comment"
  | "like"
  | "favourite"
  | "booking"
  | "blogpost"
  | "collection";

export interface EntityItem {
  userid?: string | number;
  userId?: string | number;
  entity_id?: string | number;
  entityId?: string | number;
  item_id?: string | number;
  itemId?: string | number;
  id?: string | number;
  postid?: string | number;
  postId?: string | number;
  entity_type?: string;
  entityType?: string;
  item_type?: string;
  itemType?: string;
  created_at?: string | number | Date | null;
  createdAt?: string | number | Date | null;
  image_url?: string | null;
  imageUrl?: string | null;
  caption?: string;
  [key: string]: unknown;
}

export interface TabStructure {
  mainTabButton: HTMLDivElement;
  tabSection: HTMLDivElement;
  childTabs: HTMLDivElement[];
  tabContentContainers: HTMLDivElement[];
}

export interface MainTabsResult {
  mainTabContainer: HTMLDivElement;
  mainTabButtons: HTMLDivElement;
  mainTabContents: HTMLDivElement;
}