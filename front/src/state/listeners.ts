import type { AppState, StateListener } from "./types";

export function createListenerRegistry<State extends AppState>(
  getState: () => State,
  getStateValue: (key: string, state: State) => any
) {
  const listeners = new Map<string, Set<StateListener>>();
  const deepListeners = new Map<string, Set<StateListener>>();
  const notifyQueue = new Set<string>();
  let notifyPending = false;

  function getValueByPath(path: string, source: any = getState()): any {
    if (!path) {
      return source;
    }
    return path.split(".").reduce((current, part) => current?.[part], source);
  }

  function scheduleNotify(key: string): void {
    if (!key) {
      return;
    }

    notifyQueue.add(key);
    if (notifyPending) {
      return;
    }

    notifyPending = true;
    queueMicrotask(() => {
      try {
        const queuedKeys = new Set(notifyQueue);
        const currentState = getState();

        for (const queueKey of queuedKeys) {
          const value = queueKey.includes(".")
            ? getValueByPath(queueKey, currentState)
            : getStateValue(queueKey, currentState);
          const fns = listeners.get(queueKey);
          if (!fns) {
            continue;
          }

          for (const fn of [...fns]) {
            try {
              fn(value, currentState);
            } catch (error) {
              console.error("[STATE] State listener failed:", error);
            }
          }
        }

        const deepCalls = new Map<string, Set<StateListener>>();
        for (const [path, fns] of deepListeners) {
          let affected = false;
          for (const queueKey of queuedKeys) {
            if (
              queueKey === path ||
              queueKey.startsWith(`${path}.`) ||
              path.startsWith(`${queueKey}.`)
            ) {
              affected = true;
              break;
            }
          }

          if (!affected) {
            continue;
          }

          deepCalls.set(path, new Set(fns));
        }

        for (const [path, fns] of deepCalls) {
          const value = getValueByPath(path, currentState);
          for (const fn of fns) {
            try {
              fn(value, currentState);
            } catch (error) {
              console.error("[STATE] Deep state listener failed:", error);
            }
          }
        }
      } finally {
        notifyQueue.clear();
        notifyPending = false;
      }
    });
  }

  function subscribe(key: string, fn: StateListener): () => void {
    if (typeof fn !== "function") {
      throw new TypeError("State subscriber must be a function.");
    }

    const rootKey = String(key).split(".")[0];
    if (!rootKey) {
      return () => undefined;
    }

    if (String(key).includes(".")) {
      return subscribeDeep(String(key), fn);
    }

    if (!listeners.has(key)) {
      listeners.set(key, new Set());
    }

    listeners.get(key)!.add(fn);
    return () => unsubscribe(key, fn);
  }

  function unsubscribe(key: string, fn: StateListener): void {
    const set = listeners.get(key);
    if (!set) {
      return;
    }
    set.delete(fn);
    if (set.size === 0) {
      listeners.delete(key);
    }
  }

  function subscribeDeep(path: string, fn: StateListener): () => void {
    if (typeof fn !== "function") {
      throw new TypeError("State subscriber must be a function.");
    }

    if (!deepListeners.has(path)) {
      deepListeners.set(path, new Set());
    }

    deepListeners.get(path)!.add(fn);
    return () => unsubscribeDeep(path, fn);
  }

  function unsubscribeDeep(path: string, fn: StateListener): void {
    const set = deepListeners.get(path);
    if (!set) {
      return;
    }
    set.delete(fn);
    if (set.size === 0) {
      deepListeners.delete(path);
    }
  }

  function clearAllListeners(): void {
    listeners.clear();
    deepListeners.clear();
    notifyQueue.clear();
    notifyPending = false;
  }

  return {
    scheduleNotify,
    subscribe,
    unsubscribe,
    subscribeDeep,
    unsubscribeDeep,
    clearAllListeners
  };
}
