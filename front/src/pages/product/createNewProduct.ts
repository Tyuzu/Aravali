import { createProduct } from "../../services/products/createOrEditProduct.js";

export async function CreateProduct(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";
  createProduct(contentContainer, isLoggedIn);
}
