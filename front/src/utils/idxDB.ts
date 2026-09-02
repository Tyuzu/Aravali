/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type SystemLogType = "success" | "error" | "info" | "warning";

export interface SystemLogEntry {
  id: string | number;
  title: string;
  message: string;
  type: SystemLogType;
  isRead: boolean;
  createdAt: string;
  [key: string]: unknown;
}

export interface AddSystemLogOptions {
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
const DB_VERSION = 2;
const STORE_NAME = "system_logs";
const MAX_LOGS = 200;

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
  const nextValue = { ...value };

  return {
    id: nextValue.id ?? `log-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    title: String(nextValue.title ?? "System Message"),
    message: String(nextValue.message ?? ""),
    type: normalizeLogType(nextValue.type),
    isRead: Boolean(nextValue.isRead),
    createdAt: nextValue.createdAt ? new Date(nextValue.createdAt).toISOString() : new Date().toISOString(),
    ...nextValue,
  } as T & SystemLogEntry;
}

function generateLogId(): string {
  const hasCrypto = typeof crypto !== "undefined" && typeof crypto.randomUUID === "function";
  if (hasCrypto) {
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

/**
 * Initializes and opens the IndexedDB database instance.
 * Reuses a single connection for the life of the app to avoid unnecessary reopen churn.
 */
function openDB(): Promise<IDBDatabase> {
  if (!isIndexedDBAvailable()) {
    return Promise.reject(new Error("IndexedDB is not available in this browser environment."));
  }

  if (dbPromise) {
    return dbPromise;
  }

  dbPromise = new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = (event: IDBVersionChangeEvent) => {
      const db = (event.target as IDBOpenDBRequest).result;

      if (!db.objectStoreNames.contains(STORE_NAME)) {
        const store = db.createObjectStore(STORE_NAME, { keyPath: "id" });
        store.createIndex("createdAt", "createdAt", { unique: false });
        store.createIndex("isRead", "isRead", { unique: false });
        store.createIndex("type", "type", { unique: false });
      } else {
        const store = request.transaction?.objectStore(STORE_NAME);
        if (store && !store.indexNames.contains("createdAt")) {
          store.createIndex("createdAt", "createdAt", { unique: false });
        }
        if (store && !store.indexNames.contains("isRead")) {
          store.createIndex("isRead", "isRead", { unique: false });
        }
        if (store && !store.indexNames.contains("type")) {
          store.createIndex("type", "type", { unique: false });
        }
      }
    };

    request.onsuccess = () => {
      const db = request.result;

      db.onversionchange = () => {
        db.close();
      };

      resolve(db);
    };

    request.onerror = () => {
      dbPromise = null;
      reject(request.error ?? new Error("Failed to open IndexedDB."));
    };
  });

  return dbPromise;
}

async function pruneOldLogs(): Promise<void> {
  try {
    const logs = await getAll<SystemLogEntry>();
    if (logs.length <= MAX_LOGS) return;

    const oldest = sortLogs(logs, "asc").slice(0, logs.length - MAX_LOGS);
    await Promise.all(oldest.map((entry) => remove(entry.id)));
  } catch (error) {
    console.warn("Failed to trim old IndexedDB notices:", safeErrorMessage(error));
  }
}

/* =========================================================
   DATA ACCESS LAYER (CRUD)
========================================================= */

/**
 * Retrieve a single log by ID.
 */
export async function get<T extends SystemLogEntry = SystemLogEntry>(
  id: IDBValidKey
): Promise<T | undefined> {
  if (!isIndexedDBAvailable()) return undefined;

  try {
    const db = await openDB();

    return await new Promise<T | undefined>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readonly");
      const store = tx.objectStore(STORE_NAME);
      const request = store.get(id);

      request.onsuccess = () => resolve((request.result as T | undefined) ?? undefined);
      request.onerror = () => reject(request.error ?? new Error("Failed to fetch log."));
      tx.onabort = () => reject(tx.error ?? new Error("IndexedDB read aborted."));
    });
  } catch (error) {
    console.warn(`Failed to read IndexedDB log ${String(id)}:`, safeErrorMessage(error));
    return undefined;
  }
}

/**
 * Store a new log or update an existing log.
 */
export async function set<T extends SystemLogEntry = SystemLogEntry>(
  value: T
): Promise<IDBValidKey> {
  if (!value || value.id == null) {
    throw new Error("IndexedDB log must contain an id.");
  }

  if (!isIndexedDBAvailable()) {
    return value.id;
  }

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

/**
 * Explicit update alias.
 */
export async function update<T extends SystemLogEntry = SystemLogEntry>(
  value: T
): Promise<IDBValidKey> {
  return set(value);
}

/**
 * Explicit put alias.
 */
export async function put<T extends SystemLogEntry = SystemLogEntry>(
  value: T
): Promise<IDBValidKey> {
  return set(value);
}

/**
 * Delete a single log by ID.
 */
export async function remove(id: IDBValidKey): Promise<void> {
  if (!isIndexedDBAvailable()) return;

  const db = await openDB();

  await new Promise<void>((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readwrite");
    const request = tx.objectStore(STORE_NAME).delete(id);

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error ?? new Error("Failed to delete log."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB delete aborted."));
  });
}

/**
 * Retrieve all system log entries with optional filtering/sorting.
 */
export async function getAll<T extends SystemLogEntry = SystemLogEntry>(
  options: GetAllOptions = {}
): Promise<T[]> {
  if (!isIndexedDBAvailable()) return [];

  try {
    const db = await openDB();
    const items = await new Promise<T[]>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readonly");
      const request = tx.objectStore(STORE_NAME).getAll();

      request.onsuccess = () => resolve((request.result as T[]) ?? []);
      request.onerror = () => reject(request.error ?? new Error("Failed to fetch logs."));
      tx.onabort = () => reject(tx.error ?? new Error("IndexedDB read aborted."));
    });

    let filtered = items.map((item) => normalizeLogEntry(item) as T);

    if (options.unreadOnly) {
      filtered = filtered.filter((item) => !item.isRead);
    }

    if (options.type) {
      filtered = filtered.filter((item) => item.type === normalizeLogType(options.type));
    }

    if (options.since) {
      const sinceTime = new Date(options.since).getTime();
      if (!Number.isNaN(sinceTime)) {
        filtered = filtered.filter((item) => new Date(item.createdAt).getTime() >= sinceTime);
      }
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

/**
 * Count all stored logs.
 */
export async function count(): Promise<number> {
  if (!isIndexedDBAvailable()) return 0;

  const db = await openDB();

  return new Promise((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readonly");
    const request = tx.objectStore(STORE_NAME).count();

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error ?? new Error("Failed to count logs."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB count aborted."));
  });
}

/**
 * Clear all entries from the store.
 */
export async function clear(): Promise<void> {
  if (!isIndexedDBAvailable()) return;

  const db = await openDB();

  await new Promise<void>((resolve, reject) => {
    const tx = db.transaction(STORE_NAME, "readwrite");
    const request = tx.objectStore(STORE_NAME).clear();

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error ?? new Error("Failed to clear logs."));
    tx.onabort = () => reject(tx.error ?? new Error("IndexedDB clear aborted."));
  });
}

/**
 * Convenience wrapper to mark all stored logs as read.
 */
export async function markAllAsRead(): Promise<number> {
  const unread = await getAll<SystemLogEntry>({ unreadOnly: true });
  if (!unread.length) return 0;

  await Promise.all(unread.map((entry) => set({ ...entry, isRead: true })));
  return unread.length;
}

/* =========================================================
   DOMAIN HELPERS
========================================================= */

/**
 * Helper to add a new system log.
 */
export async function addSystemLog({
  title,
  message,
  type = "info"
}: AddSystemLogOptions): Promise<SystemLogEntry> {
  const logItem: SystemLogEntry = normalizeLogEntry({
    id: generateLogId(),
    title: title?.trim() || "System Message",
    message: message?.trim() || "No details provided.",
    type: normalizeLogType(type),
    isRead: false,
    createdAt: new Date().toISOString(),
  });

  await set(logItem);
  await pruneOldLogs();

  return logItem;
}