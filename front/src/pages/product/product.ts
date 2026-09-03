
import "../../../css/farmstyles/productpage.css";
import { displayProduct } from "../../services/products/productPage.js";

export async function Product(
  isLoggedIn: boolean,
  productType: string,
  productId: string,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  displayProduct(isLoggedIn, productType, productId, contentContainer);
}
