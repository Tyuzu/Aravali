import { addToCart, isValidCartQuantity } from "../../../cart/addToCart";
import { getState } from "../../../../state/state";
import Notify from "../../../../components/ui/Notify";

export async function handleAddToCart(
  listingId: string | number,
  getQuantity: () => number,
  setDisabled: (v: boolean) => void,
  onCartUpdated?: (resp: unknown) => void
): Promise<boolean> {
  if (!isValidCartQuantity(getQuantity())) {
    console.error("Invalid cart quantity:", getQuantity());
    return false;
  }

  const isLoggedIn = Boolean(getState("token"));

  setDisabled(true);
  try {
    const success = await addToCart({
      itemId: listingId,
        itemType: "crop",
      quantity: getQuantity(),
      isLoggedIn,
      onCartUpdated: (response: unknown) => {
        try {
          onCartUpdated?.(response);
        } catch {}
      }
    });

    return Boolean(success);
  } catch (err: unknown) {
    Notify((err as Error)?.message || "Failed to add to cart", { type: "error", dismissible: true });
    return false;
  } finally {
    setDisabled(false);
  }
}
