import { createElement } from "../../../../components/createElement";
import Imagex from "../../../../components/base/Imagex";
import { resolveImagePath, PictureType, EntityType } from "../../../../utils/imagePaths";
import Button from "../../../../components/base/Button";
import { navigate } from "../../../../routes/navigate";

export interface CropCardData {
  name: string;
  minPrice?: number;
  maxPrice?: number;
  availableCount?: number;
  unit?: string;
  banner?: string;
  tags?: string[];
  seasonMonths?: number[];
  price?: number;
  quantity?: number;
  farmName?: string;
}

function formatPrice(value?: number): string {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 0
  }).format(value || 0);
}

function formatPriceRange(min?: number, max?: number): string {
  return `${formatPrice(min)} - ${formatPrice(max)}`;
}

function isSeasonal(crop: CropCardData): boolean {
  const currentMonth = new Date().getMonth() + 1;
  return Array.isArray(crop.seasonMonths) && crop.seasonMonths.includes(currentMonth);
}

function formatCropSlug(name: string): string {
  return name.toLowerCase().replace(/\s+/g, "_");
}

export function renderCropCard(crop: CropCardData, mode: "catalogue" | "listing" = "catalogue"): HTMLElement {
  const card = createElement("div", { class: "crop-card" });
  const cropSlug = formatCropSlug(crop.name);

  card.addEventListener("click", () => navigate(`/crop/${cropSlug}`));

  const img = Imagex({
    src: resolveImagePath(EntityType.CROP, PictureType.THUMB, crop.banner),
    alt: crop.name,
    class: "crop-card-image",
    loading: "lazy"
  });

  const title = createElement("h4", {}, [crop.name]);

  if (mode === "catalogue") {
    const info = createElement("p", { class: "crop-info" }, [
      `${formatPriceRange(crop.minPrice, crop.maxPrice)} per ${crop.unit} • ${crop.availableCount} listings`
    ]);

    const inSeason = isSeasonal(crop);
    const seasonLabel = inSeason ? "🟢 In Season" : "🔴 Off Season";
    const seasonClass = inSeason ? "in-season" : "off-season";

    const season = createElement("p", { class: `season-indicator ${seasonClass}` }, [seasonLabel]);
    const tags = createElement(
      "div",
      { class: "tag-wrap" },
      (crop.tags || []).map(tag => createElement("span", { class: "tag-pill" }, [tag]))
    );

    const btn = Button({
      title: "View Farms",
      type: "button",
      events: { click: () => navigate(`/crop/${cropSlug}`) },
      classes: "buttonx"
    });

    const contentWrapper = createElement("div", { class: "nimgcon" }) as HTMLElement;
    contentWrapper.append(title, info, season, tags, btn);
    card.append(img, contentWrapper);
  } else {
    card.append(
      title,
      createElement("p", {}, [`💰 ${formatPrice(crop.price)} per ${crop.unit}`]),
      createElement("p", {}, [`📦 In Stock: ${crop.quantity}`]),
      createElement("p", {}, [`👨‍🌾 Farm: ${crop.farmName || "Unknown"}`])
    );
  }

  return card;
}
