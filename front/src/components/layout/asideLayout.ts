import "../../../css/layout/aside.css";
import { createElement } from "../../components/createElement.js";
import { adspace } from "../../services/ads/newads.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type PrimitiveContent = Node | string | number | boolean | null | undefined;

/**
 * Union type supporting up to 4 levels of array nesting.
 * Avoids infinite type instantiation errors in TypeScript compiler.
 */
export type NestedContent =
  | PrimitiveContent
  | PrimitiveContent[]
  | PrimitiveContent[][]
  | PrimitiveContent[][][]
  | PrimitiveContent[][][][];

export interface AsideSectionInput {
  title?: string;
  content?: NestedContent;
  className?: string;
}

export type SectionType = Node | AsideSectionInput | null | undefined;

export interface AsideContentOptions {
  title?: string;
  actions?: NestedContent;
  sections?: SectionType[];
  children?: NestedContent;
  showAd?: boolean;
  page?: string;
  adPosition?: string;
  adPlacement?: "top" | "middle" | "bottom";
  adOptions?: Record<string, any>;
  asContainer?: boolean;
}

/* =========================================================
   HELPERS
========================================================= */

/**
 * Normalizes mixed nested inputs into a flat array of valid DOM Nodes.
 * Bypasses recursive type resolution on Array.prototype.flat using explicit casting.
 */
const normalizeContent = (content: NestedContent): Node[] => {
  if (content == null || content === false) return [];

  const flatItems = (
    Array.isArray(content) ? (content as any[]).flat(Infinity) : [content]
  ) as PrimitiveContent[];

  return flatItems
    .filter((item): item is NonNullable<PrimitiveContent> => item != null && item !== false && item !== "")
    .map((item) => (item instanceof Node ? item : document.createTextNode(String(item))));
};

/**
 * Creates structured sections or elements inside an aside layout.
 */
function renderSection(section: SectionType): HTMLElement | Node | null {
  if (!section) return null;
  if (section instanceof Node) return section;

  const children: Node[] = [];

  if (section.title) {
    children.push(createElement("h3", { class: "aside-section-title" }, [section.title]));
  }

  if (section.content) {
    children.push(...normalizeContent(section.content));
  }

  const className = ["aside-section", section.className].filter(Boolean).join(" ");
  return createElement("section", { class: className }, children);
}

/* =========================================================
   MAIN COMPONENT BUILDER
========================================================= */

/**
 * Reusable sidebar element builder with title, actions, sections, custom content, and ad placement.
 */
export function createAsideContent({
  title = "Actions",
  actions = [],
  sections = [],
  children = [],
  showAd = true,
  page,
  adPosition = "aside",
  adPlacement = "top",
  adOptions = {},
  asContainer = false
}: AsideContentOptions = {}): HTMLElement | Node[] {
  // 1. Resolve optional ad node
  const adNode: Node | null = showAd ? adspace(adPosition, page, adOptions) : null;

  // 2. Build title and actions
  const titleNode = title ? createElement("h2", { class: "aside-title" }, [title]) : null;

  const normalizedActions = normalizeContent(actions);
  const actionsContainer = normalizedActions.length > 0
    ? createElement("div", { class: "aside-actions" }, normalizedActions)
    : null;

  // 3. Process sections & children
  const renderedSections = sections.map(renderSection).filter((sec): sec is Node => sec !== null);
  const normalizedChildren = normalizeContent(children);

  // 4. Assemble components based on ad placement
  const contentNodes: Node[] = [];

  if (adPlacement === "top" && adNode) contentNodes.push(adNode);
  if (titleNode) contentNodes.push(titleNode);
  if (actionsContainer) contentNodes.push(actionsContainer);
  if (adPlacement === "middle" && adNode) contentNodes.push(adNode);
  contentNodes.push(...renderedSections, ...normalizedChildren);
  if (adPlacement === "bottom" && adNode) contentNodes.push(adNode);

  // 5. Return container element or array of nodes
  if (asContainer) {
    return createElement("aside", { class: "aside-container" }, contentNodes);
  }

  return contentNodes;
}