import { createElement } from "../../../components/createElement.js";
import Button from "../../../components/base/Button.js";
import { displayCropCard } from "./displayCropCard.js";
import type { Crop, PriceHistory } from "../types.js";

function getAgeInDays(dateValue?: string | number | Date): number {
  if (dateValue === undefined || dateValue === null || dateValue === "") {
    return 0;
  }

  const d = dateValue instanceof Date ? dateValue : new Date(dateValue);
  if (Number.isNaN(d.getTime())) {
    return 0;
  }

  return Math.floor((Date.now() - d.getTime()) / (1000 * 60 * 60 * 24));
}

function getNumericValue(value: unknown): number {
  return typeof value === "number" && !Number.isNaN(value) ? value : 0;
}

function createPriceHistoryToggle(history: PriceHistory[]): HTMLElement {
  const block = createElement(
    "pre",
    { class: "price-history hidden" },
    [history.map((price) => `${price.date}: ₹${price.price}`).join("\n")]
  ) as HTMLElement;

  const toggle = Button({
    title: "📈 Show Price History",
    classes: "buttonx",
    events: {
      click: () => block.classList.toggle("hidden")
    }
  });

  return createElement("div", { class: "price-history-wrapper" }, [toggle, block]) as HTMLElement;
}

export function renderFarmCropCard(
  crop: Crop,
  farmName: string,
  farmId: string | number,
  editcon: HTMLElement,
  isLoggedIn: boolean,
  isCreator: boolean,
  createUserControls: (
    crop: Crop,
    farmName: string,
    farmId: string | number,
    isLoggedIn: boolean
  ) => HTMLElement[],
  createCreatorControls: (crop: Crop, farmId: string | number, editcon: HTMLElement) => HTMLElement[]
): HTMLElement {
  const harvestDate = crop?.HarvestDate || crop?.harvestDate;
  const harvestAge = harvestDate ? `${getAgeInDays(harvestDate)} days old` : "Unknown age";

  let expiryNotice: HTMLElement | null = null;

  if (crop?.expiryDate) {
    const daysUntilExpiry = -getAgeInDays(crop.expiryDate);

    if (daysUntilExpiry <= 0) {
      expiryNotice = createElement("span", { class: "badge badge-expired" }, ["⚠️ Expired"]) as HTMLElement;
    } else if (daysUntilExpiry <= 2) {
      expiryNotice = createElement(
        "span",
        { class: "badge badge-warning" },
        [`⚠️ Expires in ${daysUntilExpiry} days`]
      ) as HTMLElement;
    }
  }

  const originalVal = getNumericValue(crop?.price);
  const discountVal = getNumericValue(crop?.discount);
  const controls = isCreator
    ? createCreatorControls(crop, farmId, editcon)
    : createUserControls(crop, farmName, farmId, isLoggedIn);

  const card = displayCropCard({
    name: crop?.name || "Crop",
    price: originalVal,
    unit: crop?.unit || "unit",
    quantity: Number(crop?.quantity ?? 0),
    banner: crop?.banner,
    discount: discountVal,
    farmName,
    cropid: crop?.cropid,
    farmid: farmId
  });

  const cardChildren = [
    createElement(
      "div",
      { class: "crop-header" },
      [
        createElement(
          "h4",
          { style: "display: inline-block; margin-right: 8px;" },
          [crop?.name || "Crop"]
        ),
        crop?.category
          ? createElement(
              "span",
              {
                class: "badge-category",
                style:
                  "background: var(--color-fg); color: var(--color-bg); font-size: 11px; padding: 2px 8px; border-radius: 12px; text-transform: uppercase;"
              },
              [crop.category]
            )
          : null
      ].filter((child): child is HTMLElement => child !== null)
    ),
    createElement("p", {}, [`📦 Stock: ${crop?.quantity ?? 0} ${crop?.unit || ""}`]),
    createElement("p", {}, [`🕓 Harvested: ${harvestAge}`]),
    createElement("p", {}, [(crop?.quantity ?? 0) > 0 ? "✅ Available" : "❌ Out of Stock"]),
    expiryNotice,
    crop?.history && crop.history.length > 1 ? createPriceHistoryToggle(crop.history) : null,
    ...controls
  ].flat().filter((child): child is HTMLElement => child !== null);

  card.append(...cardChildren);
  return card;
}

export default renderFarmCropCard;
