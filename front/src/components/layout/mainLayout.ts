import { createElement } from "../../components/createElement.js";
import { adspace } from "../../services/ads/newads.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type PrimitiveContent = Node | string | number | boolean | null | undefined;

/**
 * Finite nested union type to avoid infinite type instantiation depth during build.
 */
export type ContentInput =
  | PrimitiveContent
  | PrimitiveContent[]
  | PrimitiveContent[][]
  | PrimitiveContent[][][]
  | PrimitiveContent[][][][];

export interface MainLayoutConfig {
  mainContent?: ContentInput;
  asideContent?: ContentInput;
  pageClass?: string;
  page?: string;
  showMainAd?: boolean;
  mainAdPosition?: string;
  mainAdPlacement?: "top" | "bottom" | string;
  mainAdOptions?: Record<string, unknown>;
}

/* =========================================================
   HELPERS
========================================================= */

/**
 * Normalizes mixed inputs into a flat array of valid DOM Nodes.
 */
const normalizeContent = (content: ContentInput): Node[] => {
  if (content == null || content === false) return [];

  const rawArray = Array.isArray(content) ? content : [content];

  // Cast to any[] before .flat() to prevent recursive compiler unwrapping errors
  const flatItems = (rawArray as any[]).flat(10) as PrimitiveContent[];

  return flatItems
    .filter((item): item is NonNullable<PrimitiveContent> => item != null && item !== false && item !== "")
    .map((item) => (item instanceof Node ? item : document.createTextNode(String(item))));
};

/* =========================================================
   MAIN LAYOUT BUILDER
========================================================= */

/**
 * Creates a standard two-column page structure with a main content area and an aside sidebar.
 */
export function createMainLayout({
  mainContent = [],
  asideContent = [],
  pageClass = "page-layout",
  page,
  showMainAd = false,
  mainAdPosition = "main-bottom",
  mainAdPlacement = "bottom",
  mainAdOptions = {}
}: MainLayoutConfig = {}): HTMLDivElement {
  // 1. Resolve optional main ad node
  const mainAdNode: Node | null = showMainAd
    ? adspace(mainAdPosition, page, mainAdOptions)
    : null;

  // 2. Normalize main content
  const normalizedMain = normalizeContent(mainContent);

  // 3. Assemble main section children cleanly based on placement
  const finalMainContent: Node[] = [];

  if (mainAdPlacement === "top" && mainAdNode) {
    finalMainContent.push(mainAdNode);
  }

  finalMainContent.push(...normalizedMain);

  if (mainAdPlacement === "bottom" && mainAdNode) {
    finalMainContent.push(mainAdNode);
  }

  // 4. Construct layout containers
  const containerClass = ["two-column", pageClass].filter(Boolean).join(" ");

  const main = createElement("main", { class: "layout-main" }, finalMainContent);
  const aside = createElement("aside", { class: "layout-aside" }, normalizeContent(asideContent));

  return createElement("div", { class: containerClass }, [main, aside]) as HTMLDivElement;
}