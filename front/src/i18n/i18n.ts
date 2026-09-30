import { setState } from "../state/state.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type TranslationDictionary = {
  [key: string]: string | TranslationDictionary;
};

export type InterpolationVars = {
  count?: number;
  [key: string]: any;
};

/* =========================================================
   STATE & CACHE
========================================================= */

let translations: TranslationDictionary = {};
let currentLang: string = "en";
let activeRequest: number = 0;

const cache = new Map<string, Promise<TranslationDictionary>>();
let cachedPluralRules: Intl.PluralRules | null = null;

const SUPPORTED_LANGS: readonly string[] = ["en", "es", "fr", "hi", "ar", "jp"];
const FALLBACK_LANG: string = "en";

/* =========================================================
   HELPERS & LOCALE NORMALIZATION
========================================================= */

/**
 * Normalizes input language strings to internal app codes ('jp' instead of 'ja').
 */
function normalizeLanguageCode(lang: string): string {
  const raw = (lang || "").trim().toLowerCase();
  if (!raw) return FALLBACK_LANG;

  if (raw === "ja" || raw === "jp") return "jp";

  if (SUPPORTED_LANGS.includes(raw)) return raw;

  const base = raw.split("-")[0];
  if (base === "ja") return "jp";
  if (SUPPORTED_LANGS.includes(base)) return base;

  return FALLBACK_LANG;
}

/**
 * Maps internal app language codes to valid BCP 47 locale tags for Intl APIs.
 */
function toBCP47Locale(lang: string): string {
  return lang === "jp" ? "ja" : lang;
}

/**
 * Fetch language JSON with basic error boundary and cache invalidation.
 */
function fetchTranslations(lang: string): Promise<TranslationDictionary> {
  return fetch(`/i18n/${lang}.json`, { cache: "no-cache" })
    .then((res) => {
      if (!res.ok) throw new Error(`Failed to load ${lang}: ${res.status}`);
      return res.json() as Promise<TranslationDictionary>;
    })
    .catch((err) => {
      cache.delete(lang);
      throw err;
    });
}

/**
 * Load and apply target language dictionary with active request cancellation guard.
 */
async function loadTranslations(lang: string): Promise<void> {
  const requestId = ++activeRequest;

  try {
    if (!cache.has(lang)) {
      cache.set(lang, fetchTranslations(lang));
    }

    const data = await cache.get(lang);

    // Stale request guard: Ignore out-of-order async completions
    if (requestId !== activeRequest) return;

    if (!data || typeof data !== "object" || Array.isArray(data)) {
      throw new Error(`Invalid translation data structure for ${lang}`);
    }

    translations = data;
    currentLang = lang;
    cachedPluralRules = null; // Invalidate cached plural rules for new language

    localStorage.setItem("lang", lang);
    setState("lang", lang);
    setState("isI18nLoaded", true);
  } catch (err) {
    if (requestId !== activeRequest) return;

    console.error(`Failed to load translations for "${lang}"`, err);

    if (lang !== FALLBACK_LANG) {
      return loadTranslations(FALLBACK_LANG);
    }

    translations = {};
    currentLang = FALLBACK_LANG;
    cachedPluralRules = null;
  }
}

/* =========================================================
   PUBLIC API
========================================================= */

export async function setLanguage(lang: string): Promise<void> {
  const targetLang = normalizeLanguageCode(lang);
  await loadTranslations(targetLang);
}

export function detectLanguage(): string {
  const saved = localStorage.getItem("lang");
  if (saved) return normalizeLanguageCode(saved);

  const langs = navigator.languages || [navigator.language];
  for (const lang of langs) {
    const normalized = normalizeLanguageCode(lang);
    if (SUPPORTED_LANGS.includes(normalized)) {
      return normalized;
    }
  }

  return FALLBACK_LANG;
}

export const getCurrentLanguage = (): string => currentLang;
// Export alias for components expecting getLanguage()
export const getLanguage = getCurrentLanguage;

/**
 * Safely resolve nested property values from translation objects using dot notation.
 */
function getNested(obj: TranslationDictionary, path: string): unknown {
  if (!obj || typeof obj !== "object") return undefined;

  return path.split(".").reduce<unknown>((value, key) => {
    if (value === null || value === undefined || typeof value !== "object") {
      return undefined;
    }
    return (value as Record<string, unknown>)[key];
  }, obj);
}

/**
 * Translate key with interpolation and pluralization support.
 */
export function t(key: string, vars: InterpolationVars = {}, fallback: string = ""): string {
  const isDev = (import.meta as any).env?.DEV;

  if (typeof key !== "string" || !key.trim()) {
    if (isDev) console.warn("Missing or non-string translation key:", key);
    return fallback || "";
  }

  let template: unknown = getNested(translations, key);

  // Pluralization handling
  if (typeof vars.count === "number") {
    if (!cachedPluralRules) {
      // Convert internal 'jp' code to valid BCP 47 'ja' code for Intl.PluralRules
      cachedPluralRules = new Intl.PluralRules(toBCP47Locale(currentLang));
    }
    const rule = cachedPluralRules.select(vars.count);
    const pluralKey = `${key}.${rule}`;
    const plural = getNested(translations, pluralKey);

    if (typeof plural === "string") {
      template = plural;
    }
  }

  // Fallback handling
  if (typeof template !== "string") {
    if (isDev) {
      console.warn(`Missing translation key: "${key}"`, {
        language: currentLang,
        value: template
      });
    }
    template = fallback || key;
  }

  // String variable interpolation ({year}, {name}, etc.)
  return (template as string).replace(/\{(\w+)\}/g, (_, variable: string) => {
    return Object.prototype.hasOwnProperty.call(vars, variable)
      ? String(vars[variable])
      : `{${variable}}`;
  });
}

/**
 * Initialize i18n system on app boot.
 */
export async function initI18n(): Promise<void> {
  await setLanguage(detectLanguage());
}