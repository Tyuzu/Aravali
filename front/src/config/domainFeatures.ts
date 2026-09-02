// src/config/domainFeatures.ts

export type FeatureKey = string;

export interface DomainMetadata {
  title: string;
  description: string;
  theme?: string;
  logo?: string;
  favicon?: string;
}

export const DOMAIN_FEATURE_MAP: Record<string, FeatureKey[]> = {
  "farms.netlify.app": ["core", "farms", "chats"],
  "places.netlify.app": ["core", "places", "chats", "events", "baito"],
  "events.netlify.app": ["core", "events"],
  "baito.netlify.app": ["core", "baito", "chats", "places"],
  "chats.netlify.app": ["core", "chats"],
  "admin.netlify.app": ["core", "admin"],
  "social.netlify.app": ["core", "social", "chats"],

  // Local development, staging, and main hub access
  "localhost": ["ALL"],
  "127.0.0.1": ["ALL"],
  "indium.netlify.app": ["ALL"]
};

export const DOMAIN_METADATA: Record<string, DomainMetadata> = {
  "farms.netlify.app": {
    title: "FarmHub",
    description: "Discover local farms, fresh produce, and community crops.",
    theme: "theme-green",
    logo: "/logos/farm.svg",
    favicon: "/favicons/farm.ico"
  },
  "places.netlify.app": {
    title: "Places",
    description: "Explore local spots and community landmarks.",
    theme: "theme-emerald",
    logo: "/logos/places.svg",
    favicon: "/favicons/places.ico"
  },
  "events.netlify.app": {
    title: "EventPulse",
    description: "Find concerts, gatherings, and local experiences.",
    theme: "theme-purple",
    logo: "/logos/events.svg",
    favicon: "/favicons/events.ico"
  },
  "baito.netlify.app": {
    title: "BaitoJobs",
    description: "Part-time gigs, local hiring, and quick work.",
    theme: "theme-amber",
    logo: "/logos/baito.svg",
    favicon: "/favicons/baito.ico"
  },
  "chats.netlify.app": {
    title: "MereChat",
    description: "Connect and message directly with community members.",
    theme: "theme-blue",
    logo: "/logos/chats.svg",
    favicon: "/favicons/chats.ico"
  },
  "social.netlify.app": {
    title: "Community Posts",
    description: "Share updates, fan media, and community stories.",
    theme: "theme-pink",
    logo: "/logos/social.svg",
    favicon: "/favicons/social.ico"
  },
  "admin.netlify.app": {
    title: "System Admin",
    description: "Platform management and system controls.",
    theme: "theme-slate",
    logo: "/logos/admin.svg",
    favicon: "/favicons/admin.ico"
  },

  // Defaults for Local Development & Main Domain
  "localhost": {
    title: "Dev Suite",
    description: "Local development suite with full access.",
    theme: "theme-default",
    logo: "/logos/main.svg",
    favicon: "/favicon.ico"
  },
  "127.0.0.1": {
    title: "Dev Suite",
    description: "Local development suite with full access.",
    theme: "theme-default",
    logo: "/logos/main.svg",
    favicon: "/favicon.ico"
  },
  "netlify.app": {
    title: "Main Network",
    description: "All-in-one community hub.",
    theme: "theme-default",
    logo: "/logos/main.svg",
    favicon: "/favicon.ico"
  }
};

/** Gets the feature override parameter from URL if present (`?feature=events`). */
function getUrlFeatureOverride(): string | null {
  if (typeof window === "undefined") return null;
  const urlParams = new URLSearchParams(window.location.search);
  const override = urlParams.get("feature");
  return override ? override.toLowerCase() : null;
}

/** Returns allowed features for the current hostname or URL override. */
export function getCurrentAllowedFeatures(): FeatureKey[] {
  const featureOverride = getUrlFeatureOverride();
  if (featureOverride) {
    return ["core", featureOverride];
  }

  const hostname = typeof window !== "undefined" ? window.location.hostname : "localhost";
  return DOMAIN_FEATURE_MAP[hostname] || ["ALL"];
}

/** Returns metadata object for current hostname or target URL override. */
export function getActiveDomainMetadata(): DomainMetadata {
  const featureOverride = getUrlFeatureOverride();
  if (featureOverride) {
    const direct = DOMAIN_METADATA[featureOverride];
    if (direct) return direct;
    const dotted = DOMAIN_METADATA[`${featureOverride}.netlify.app`];
    if (dotted) return dotted;
  }

  const hostname = typeof window !== "undefined" ? window.location.hostname : "localhost";
  return DOMAIN_METADATA[hostname] || DOMAIN_METADATA["localhost"] || {
    title: "App Network",
    description: "Community platform",
    theme: "theme-default",
    logo: "/logos/main.svg",
    favicon: "/favicon.ico"
  };
}