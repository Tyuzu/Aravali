/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type SystemLogType = "success" | "error" | "info" | "warning";

export interface SystemLogEntry {
  id: string | number;
  userId: string;
  title: string;
  message: string;
  type: SystemLogType;
  isRead: boolean;
  createdAt: string;
  [key: string]: unknown;
}

export interface AddSystemLogOptions {
  userId: string;
  title: string;
  message: string;
  type?: SystemLogType | string;
}

export interface GetAllOptions {
  order?: "asc" | "desc";
  limit?: number;
  unreadOnly?: boolean;
  type?: SystemLogType | string;
  since?: string | number | Date;
}

/* =========================================================
   CONSTANTS & DATABASE INIT
========================================================= */

const DB_NAME = "AppNotificationsDB";
const DB_VERSION = 3;
const STORE_NAME = "system_logs";
const MAX_LOGS_PER_USER = 200;

let dbPromise: Promise<IDBDatabase> | null = null;

function isIndexedDBAvailable(): boolean {
  return typeof indexedDB !== "undefined";
}

function normalizeLogType(type?: SystemLogType | string): SystemLogType {
  const nextType = (type || "info").toLowerCase();
  if (nextType === "success" || nextType === "error" || nextType === "warning") {
    return nextType;
  }
  return "info";
}

function normalizeLogEntry<T extends Partial<SystemLogEntry>>(value: T): T & SystemLogEntry {
  if (!value.userId) {
    throw new Error("IndexedDB log entry must include a valid userId.");
  }

  return {
    id: value.id ?? `log-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    userId: String(value.userId),
    title: String(value.title ?? "System Message"),
    message: String(value.message ?? ""),
    type: normalizeLogType(value.type),
    isRead: Boolean(value.isRead),
    createdAt: value.createdAt ? new Date(value.createdAt).toISOString() : new Date().toISOString(),
    ...value,
  } as T & SystemLogEntry;
}

function generateLogId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return `log-${crypto.randomUUID()}`;
  }
  return `log-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function sortLogs<T extends SystemLogEntry>(items: T[], order: "asc" | "desc" = "desc"): T[] {
  return [...items].sort((a, b) => {
    const aTime = new Date(a.createdAt || 0).getTime();
    const bTime = new Date(b.createdAt || 0).getTime();
    return order === "asc" ? aTime - bTime : bTime - aTime;
  });
}

function safeErrorMessage(error: unknown): string {
  if (error instanceof Error) return error.message;
  return String(error ?? "Unknown IndexedDB error");
}

function openDB(): Promise<IDBDatabase> {
  if (!isIndexedDBAvailable()) {
    return Promise.reject(new Error("IndexedDB is not available in this browser environment."));
  }

  if (dbPromise) return dbPromise;

  dbPromise = new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = (event: IDBVersionChangeEvent) => {
      const db = (event.target as IDBOpenDBRequest).result;
      let store: IDBObjectStore;

      if (!db.objectStoreNames.contains(STORE_NAME)) {
        store = db.createObjectStore(STORE_NAME, { keyPath: "id" });
      } else {
        store = request.transaction!.objectStore(STORE_NAME);
      }

      if (!store.indexNames.contains("userId")) store.createIndex("userId", "userId", { unique: false });
      if (!store.indexNames.contains("createdAt")) store.createIndex("createdAt", "createdAt", { unique: false });
      if (!store.indexNames.contains("isRead")) store.createIndex("isRead", "isRead", { unique: false });
      if (!store.indexNames.contains("type")) store.createIndex("type", "type", { unique: false });
    };

    request.onsuccess = () => {
      const db = request.result;
      db.onversionchange = () => db.close();
      resolve(db);
    };

    request.onerror = () => {
      dbPromise = null;
      reject(request.error ?? new Error("Failed to open IndexedDB."));
    };
  });

  return dbPromise;
}

async function pruneOldLogs(userId: string): Promise<void> {
  try {
    const logs = await getAll<SystemLogEntry>(userId);
    if (logs.length <= MAX_LOGS_PER_USER) return;

    const oldest = sortLogs(logs, "asc").slice(0, logs.length - MAX_LOGS_PER_USER);
    await Promise.all(oldest.map((entry) => remove(entry.id, userId)));
  } catch (error) {
    console.warn("Failed to trim old IndexedDB notices:", safeErrorMessage(error));
  }
}

/* =========================================================
   DATA ACCESS LAYER (CRUD)
========================================================= */

export async function get<T extends SystemLogEntry = SystemLogEntry>(
  id: IDBValidKey,
  userId: string
): Promise<T | undefined> {
  if (!isIndexedDBAvailable() || !userId) return undefined;

  try {
    const db = await openDB();
    const log = await new Promise<T | undefined>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readonly");
      const store = tx.objectStore(STORE_NAME);
      const request = store.get(id);

      request.onsuccess = () => resolve((request.result as T | undefined) ?? undefined);
      request.onerror = () => reject(request.error ?? new Error("Failed to fetch log."));
      tx.onabort = () => reject(tx.error ?? new Error("IndexedDB read aborted."));
    });

    return log && log.userId === userId ? log : undefined;
  } catch (error) {
    console.warn(`Failed to read IndexedDB log ${String(id)}:`, safeErrorMessage(error));
    return undefined;
  }
}

export async function set<T extends SystemLogEntry = SystemLogEntry>(
  value: T
): Promise<IDBValidKey> {
  if (!value || value.id == null) throw new Error("IndexedDB log must contain an id.");
  if (!value.userId) throw new Error("IndexedDB log must contain a userId.");
  if (!isIndexedDBAvailable()) return value.id;

  const normalizedValue = normalizeLogEntry(value);
  const db = await openDB();

  return new Promise((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readwrite");
    const store = tx.objectStore(STORE_NAME);
    const request = store.put(normalizedValue);

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error ?? new Error("Failed to save log."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB write aborted."));
  });
}

export async function update<T extends SystemLogEntry = SystemLogEntry>(value: T): Promise<IDBValidKey> {
  return set(value);
}

export async function put<T extends SystemLogEntry = SystemLogEntry>(value: T): Promise<IDBValidKey> {
  return set(value);
}

export async function remove(id: IDBValidKey, userId: string): Promise<void> {
  if (!isIndexedDBAvailable() || !userId) return;

  const existing = await get(id, userId);
  if (!existing) return;

  const db = await openDB();

  return new Promise<void>((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readwrite");
    const request = tx.objectStore(STORE_NAME).delete(id);

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error ?? new Error("Failed to delete log."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB delete aborted."));
  });
}

export async function getAll<T extends SystemLogEntry = SystemLogEntry>(
  userId: string,
  options: GetAllOptions = {}
): Promise<T[]> {
  if (!isIndexedDBAvailable() || !userId) return [];

  try {
    const db = await openDB();
    const items = await new Promise<T[]>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readonly");
      const store = tx.objectStore(STORE_NAME);
      const index = store.index("userId");
      const request = index.getAll(IDBKeyRange.only(userId));

      request.onsuccess = () => resolve((request.result as T[]) ?? []);
      request.onerror = () => reject(request.error ?? new Error("Failed to fetch logs."));
      tx.onabort = () => reject(tx.error ?? new Error("IndexedDB read aborted."));
    });

    let filtered = items.map((item) => normalizeLogEntry(item) as T);

    if (options.unreadOnly) filtered = filtered.filter((item) => !item.isRead);
    if (options.type) filtered = filtered.filter((item) => item.type === normalizeLogType(options.type));
    if (options.since) {
      const sinceTime = new Date(options.since).getTime();
      if (!Number.isNaN(sinceTime)) filtered = filtered.filter((item) => new Date(item.createdAt).getTime() >= sinceTime);
    }

    filtered = sortLogs(filtered, options.order ?? "desc");

    if (typeof options.limit === "number" && options.limit > 0) {
      filtered = filtered.slice(0, options.limit);
    }

    return filtered;
  } catch (error) {
    console.warn("Failed to read IndexedDB logs:", safeErrorMessage(error));
    return [];
  }
}

export async function count(userId: string): Promise<number> {
  if (!isIndexedDBAvailable() || !userId) return 0;
  const db = await openDB();

  return new Promise((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readonly");
    const index = tx.objectStore(STORE_NAME).index("userId");
    const request = index.count(IDBKeyRange.only(userId));

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error ?? new Error("Failed to count logs."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB count aborted."));
  });
}

/**
 * Optimized single-transaction cursor cleanup for user logs.
 */
export async function clear(userId: string): Promise<void> {
  if (!isIndexedDBAvailable() || !userId) return;
  const db = await openDB();

  return new Promise<void>((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readwrite");
    const store = tx.objectStore(STORE_NAME);
    const index = store.index("userId");
    const request = index.openCursor(IDBKeyRange.only(userId));

    request.onsuccess = (event) => {
      const cursor = (event.target as IDBRequest<IDBCursorWithValue>).result;
      if (cursor) {
        cursor.delete();
        cursor.continue();
      }
    };

    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error ?? new Error("Failed to clear user logs."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB clear aborted."));
  });
}

export async function markAllAsRead(userId: string): Promise<number> {
  if (!userId) return 0;

  const unread = await getAll<SystemLogEntry>(userId, { unreadOnly: true });
  if (!unread.length) return 0;

  await Promise.all(unread.map((entry) => set({ ...entry, isRead: true })));
  return unread.length;
}

export async function addSystemLog({
  userId,
  title,
  message,
  type = "info"
}: AddSystemLogOptions): Promise<SystemLogEntry> {
  if (!userId) throw new Error("Cannot add system log without a valid userId.");

  const logItem: SystemLogEntry = normalizeLogEntry({
    id: generateLogId(),
    userId,
    title: title?.trim() || "System Message",
    message: message?.trim() || "No details provided.",
    type: normalizeLogType(type),
    isRead: false,
    createdAt: new Date().toISOString(),
  });

  await set(logItem);
  await pruneOldLogs(userId);

  return logItem;
}