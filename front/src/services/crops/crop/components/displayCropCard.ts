import Imagex from "../../../../components/base/Imagex";
import { createElement } from "../../../../components/createElement";
import { resolveImagePath, EntityType, PictureType } from "../../../../utils/imagePaths";
import Button from "../../../../components/base/Button";
import { navigate } from "../../../../routes/navigate";

export interface CropCardData {
    name: string;
    price: number;
    unit: string;
    quantity: number;
    banner?: string;
    discount?: number;
    farmName?: string;
    cropid?: string | number;
    farmid?: string | number;
}

/**
 * Creates a standalone crop product card component.
 */
export function displayCropCard(crop: CropCardData): HTMLElement {
    const card = createElement("div", { class: "crop-card" }) as HTMLElement;

    // 1. Image Header
    if (crop.banner) {
        const imageSrc = resolveImagePath(EntityType.CROP, PictureType.BANNER, crop.banner);
        const cropImg = Imagex({
            src: imageSrc,
            alt: crop.name || "Crop Image",
            classes: "crop-card-image",
            loading: "lazy"
        }) as HTMLElement;

        card.appendChild(cropImg);
    }

    // 2. Pricing Calculations (robust handling for percent or fractional inputs)
    const discountRaw = Number(crop.discount ?? 0);
    const basePrice = Number(crop.price || 0);
    // If discountRaw is between 0 and 1 it's treated as a fraction (e.g. 0.15 = 15%).
    // Otherwise treat it as a percentage and clamp to [0,100].
    const discountFraction = discountRaw > 0 && discountRaw <= 1
        ? discountRaw
        : Math.min(Math.max(discountRaw, 0), 100) / 100;
    const hasDiscount = discountFraction > 0;
    const discountedPrice = hasDiscount
        ? Number((basePrice * (1 - discountFraction)).toFixed(2))
        : basePrice;

    // 3. Card Elements
    const elements: (HTMLElement | null)[] = [
        createElement("h4", { class: "crop-card-title" }, [crop.name || "Unnamed Crop"]),
        createElement("p", { class: "crop-card-price" }, [
            `💰 ₹${discountedPrice.toLocaleString("en-IN")} per ${crop.unit || "unit"}`
        ]),
        hasDiscount
            ? createElement(
                  "p",
                  { class: "crop-card-discount", style: { color: "#e53935", fontWeight: "bold" } },
                  [`${Number((discountFraction * 100).toFixed(0))}% OFF`]
              )
            : null,
        createElement("p", { class: "crop-card-stock" }, [`📦 In Stock: ${crop.quantity ?? 0}`]),
        createElement("p", { class: "crop-card-farm" }, [`👨‍🌾 Farm: ${crop.farmName || "Unknown"}`]),
        crop.cropid
            ? Button({
                  title: "View Details",
                  classes: "buttonx primary crop-card-btn",
                  events: {
                      click: (e: Event) => {
                          e.stopPropagation();
                          navigate(`/crop/${crop.cropid}`);
                      }
                  }
              })
            : null
    ];

    // Filter out null nodes and append to card container
    card.append(...(elements.filter((el): el is HTMLElement => Boolean(el))));

    return card;
}

export default displayCropCard;