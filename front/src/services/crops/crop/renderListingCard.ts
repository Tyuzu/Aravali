import { createElement } from "../../../components/createElement";
import Button from "../../../components/base/Button";
import { navigate } from "../../../routes/navigate";
import { getState } from "../../../state/state";
import { createQuantityControl } from "./components/quantityControl";
import { handleAddToCart as addToCartController } from "./controllers/cartController";

// --- Types & Interfaces ---

export interface CropListingItem {
  cropid: string;
  farmid: string;
  farmName?: string;
  location?: string;
  breed?: string;
  pricePerKg?: number;
  [key: string]: unknown;
}

const MAX_QUANTITY = 99;

/**
 * Renders a listing card with reactive quantity control and cart action.
 */
export function renderListingCard(listing: CropListingItem): HTMLElement {
  // quantity control component
  const quantityCtrl = createQuantityControl(1, MAX_QUANTITY);
  const quantityWrapper = quantityCtrl.element;
  // 2. Navigation Elements
  const farmUrl = `/farm/${listing.farmid}`;
  const farmLink = createElement(
    "a",
    {
      href: farmUrl,
      events: {
        click: (event: Event): void => {
          event.preventDefault();
          navigate(farmUrl);
        }
      }
    },
    [listing.farmName ?? "Unknown farm"]
  );

  // 3. Cart Handler (delegated to controller)
  const handleAddToCart = async (): Promise<void> => {
    await addToCartController(listing.cropid as string | number, quantityCtrl.getQuantity, quantityCtrl.setDisabled, (resp: unknown) => {
      console.debug("Cart updated:", resp);
    });
  };

  // 4. Integrated Button Component
  const addToCartButton = Button({
    title: "Add To Cart",
    id: "a2c-crop",
    events: { click: handleAddToCart },
    classes: "buttonx"
  });

  // 5. Structure Assembly
  return createElement(
    "div",
    { class: "listing-card" },
    [
      farmLink,
      createElement("p", {}, [`Location: ${listing.location ?? "N/A"}`]),
      createElement("p", {}, [`Breed: ${listing.breed ?? "N/A"}`]),
      createElement("p", {}, [
        `Price: ₹${listing.pricePerKg !== undefined ? listing.pricePerKg.toLocaleString("en-IN") : "N/A"} per kg`
      ]),
      createElement("label", {}, ["Quantity (kg):"]),
      quantityWrapper,
      addToCartButton
    ]
  ) as HTMLElement;
}

export default renderListingCard;