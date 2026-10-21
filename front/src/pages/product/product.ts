
import "../../../css/farmstyles/productpage.css";
import "../../../css/subpages/product.css";
import { displayProduct } from "../../services/products/productPage.js";

export async function Product(
  isLoggedIn: boolean,
  productTypeOrParams: string | Record<string, any>,
  productIdOrContainer: string | HTMLElement,
  maybeContainer?: HTMLElement
): Promise<void> {
  // Normalize arguments because router may call with (auth, params, container)
  let productType = "product";
  let productId = "";
  let contentContainer: HTMLElement | null = null;

  if (typeof productTypeOrParams === "string") {
    // Called as (isLoggedIn, type, id, container)
    productType = productTypeOrParams;
    productId = typeof productIdOrContainer === "string" ? productIdOrContainer : "";
    contentContainer = maybeContainer ?? (productIdOrContainer instanceof HTMLElement ? productIdOrContainer : null);
  } else if (productTypeOrParams && typeof productTypeOrParams === "object") {
    // Called as (isLoggedIn, params, container)
    const params = productTypeOrParams as Record<string, any>;
    productType = params.type || params.t || "product";
    productId = params.id || params.productid || "";
    contentContainer = (productIdOrContainer instanceof HTMLElement) ? productIdOrContainer : null;
  }

  if (!contentContainer) {
    throw new Error("Product page received invalid container");
  }

  contentContainer.innerHTML = "";
  await displayProduct(Boolean(isLoggedIn), productType, productId, contentContainer);
}
