import { setState } from "../state/state.js";

type TranslationDictionary = {
  [key: string]: string | TranslationDictionary;
};

type InterpolationVars = {
  count?: number;
  [key: string]: any;
};

let translations: TranslationDictionary = {};
let currentLang: string = "en";
let activeRequest: number = 0;
const cache = new Map<string, Promise<TranslationDictionary>>();
let cachedPluralRules: Intl.PluralRules | null = null;

const SUPPORTED_LANGS: readonly string[] = ["en", "es", "fr", "hi", "ar", "jp"];
const FALLBACK_LANG: string = "en";
const LEGACY_TRANSLATION_ALIASES: Record<string, readonly string[]> = {
  "nav.baitos": ["nav.baito"],
  "nav.baito": ["nav.baitos"],
  "nav.dash": ["nav.dashboard", "common.dashboard"],
  "nav.more": ["common.more"],
  "footer.copyright": ["common.copyright"],
  "footer.lang.en": ["common.lang.en"],
  "footer.lang.es": ["common.lang.es"],
  "footer.lang.fr": ["common.lang.fr"],
  "footer.lang.hi": ["common.lang.hi"],
  "footer.lang.ar": ["common.lang.ar"],
  "footer.lang.jp": ["common.lang.jp"]
};

function normalizeLanguageCode(lang: string): string {
  const raw = (lang || "").trim().toLowerCase();
  if (!raw) {
    return FALLBACK_LANG;
  }

  if (raw === "ja" || raw === "jp") {
    return "jp";
  }

  if (SUPPORTED_LANGS.includes(raw)) {
    return raw;
  }

  const base = raw.split("-")[0];
  if (base === "ja") {
    return "jp";
  }
  if (SUPPORTED_LANGS.includes(base)) {
    return base;
  }

  return FALLBACK_LANG;
}

function fetchTranslations(lang: string): Promise<TranslationDictionary> {
  return fetch(`/i18n/${lang}.json`, {
    cache: "no-cache"
  })
    .then((res) => {
      if (!res.ok) throw new Error(`Failed to load ${lang}: ${res.status}`);
      return res.json() as Promise<TranslationDictionary>;
    })
    .catch((err) => {
      cache.delete(lang);
      throw err;
    });
}

async function loadTranslations(lang: string): Promise<void> {
  const requestId = ++activeRequest;
  try {
    if (!cache.has(lang)) {
      cache.set(lang, fetchTranslations(lang));
    }
    const data = await cache.get(lang);
    if (requestId !== activeRequest) return;
    if (!data || typeof data !== "object" || Array.isArray(data)) {
      throw new Error(`Invalid translation data for ${lang}`);
    }
    translations = data;
    currentLang = lang;
    cachedPluralRules = null;
    localStorage.setItem("lang", lang);
    setState("lang", lang);
  } catch (err) {
    cache.delete(lang);
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

export async function setLanguage(lang: string): Promise<void> {
  const targetLang = normalizeLanguageCode(lang);
  await loadTranslations(targetLang);
}

export function detectLanguage(): string {
  const saved = localStorage.getItem("lang");
  if (saved) {
    return normalizeLanguageCode(saved);
  }

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

function getNested(obj: TranslationDictionary, path: string): unknown {
  if (!obj || typeof obj !== "object") {
    return undefined;
  }

  const segments = path.split(".");
  let current: unknown = obj;

  for (let index = 0; index < segments.length; index += 1) {
    const segment = segments[index];
    if (current === null || current === undefined || typeof current !== "object") {
      return undefined;
    }

    const record = current as Record<string, unknown>;
    const exactMatch = record[segment];
    if (Object.prototype.hasOwnProperty.call(record, segment)) {
      current = exactMatch;
      continue;
    }

    const remainder = segments.slice(index).join(".");
    if (Object.prototype.hasOwnProperty.call(record, remainder)) {
      return record[remainder];
    }

    return undefined;
  }

  return current;
}

function resolveTranslation(key: string): unknown {
  const direct = getNested(translations, key);
  if (typeof direct === "string") {
    return direct;
  }

  const aliases = LEGACY_TRANSLATION_ALIASES[key] ?? [];
  for (const alias of aliases) {
    const value = getNested(translations, alias);
    if (typeof value === "string") {
      return value;
    }
  }

  return direct;
}

export function t(key: string, vars: InterpolationVars = {}, fallback: string = ""): string {
  const isDev = (import.meta as any).env?.DEV;

  if (typeof key !== "string" || !key.trim()) {
    if (isDev) {
      console.warn("Missing or non-string translation key:", key);
    }
    return fallback || "";
  }

  let template: unknown = resolveTranslation(key);

  if (typeof vars.count === "number") {
    if (!cachedPluralRules) {
      cachedPluralRules = new Intl.PluralRules(currentLang);
    }
    const rule = cachedPluralRules.select(vars.count);
    const pluralKey = `${key}.${rule}`;
    const plural = getNested(translations, pluralKey);
    if (typeof plural === "string") {
      template = plural;
    }
  }

  if (typeof template !== "string") {
    const hasLoadedTranslations = Object.keys(translations || {}).length > 0;
    if (isDev && hasLoadedTranslations && !fallback) {
      console.warn(`Missing or non-string translation key: ${key}`, {
        language: currentLang,
        value: template,
        availableRootKeys: Object.keys(translations || {}),
      });
    }
    template = fallback || key;
  }

  // Explicit type assertion ensures TypeScript recognizes template as a string
  return (template as string).replace(/\{(\w+)\}/g, (_, variable: string) => {
    return Object.prototype.hasOwnProperty.call(vars, variable)
      ? String(vars[variable])
      : `{${variable}}`;
  });
}

export async function initI18n(): Promise<void> {
  await setLanguage(detectLanguage());
}