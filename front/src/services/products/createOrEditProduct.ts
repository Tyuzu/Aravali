import { createFormGroupWithValidation } from "../../components/form/createFormGroupEnhanced.js";
import { createElement } from "../../components/createElement.js";

export function createProduct(container: HTMLElement, isLoggedIn: boolean): void {
  // Clear container
  container.innerHTML = "";

  if (!isLoggedIn) {
    const notice = createElement("div", { class: "auth-warning" }, [
      "You must be logged in to create a product or list a tool for sale."
    ]);
    container.appendChild(notice);
    return;
  }

  const form = createElement("form", { class: "product-form" });

  // --- 1. Basic Information ---
  const typeGroup = createFormGroupWithValidation({
    type: "select",
    id: "product-type",
    name: "productType",
    label: "Listing Type",
    required: true,
    options: [
      { value: "product", label: "Standard Product (Physical / Digital)" },
      { value: "tool", label: "Tool / Equipment (For Rent or Sale)" }
    ],
    value: "product"
  });

  const titleGroup = createFormGroupWithValidation({
    type: "text",
    id: "product-title",
    name: "title",
    label: "Product / Tool Name",
    required: true,
    placeholder: "e.g., Professional Cordless Drill or Handmade Ceramic Mug",
    validator: (val) => (String(val).trim().length < 3 ? "Title must be at least 3 characters long." : null)
  });

  const descGroup = createFormGroupWithValidation({
    type: "textarea",
    id: "product-desc",
    name: "description",
    label: "Description",
    required: true,
    placeholder: "Provide details about condition, features, dimensions, or usage instructions..."
  });

  const priceGroup = createFormGroupWithValidation({
    type: "number",
    id: "product-price",
    name: "price",
    label: "Price ($)",
    required: true,
    placeholder: "0.00",
    validator: (val) => (Number(val) <= 0 ? "Price must be greater than zero." : null)
  });

  // --- 2. Conditional Container for Tool-Specific Fields ---
  const conditionalContainer = createElement("div", { class: "conditional-tool-fields", style: "display: none;" });

  const conditionGroup = createFormGroupWithValidation({
    type: "select",
    id: "tool-condition",
    name: "condition",
    label: "Condition",
    options: ["Brand New", "Like New", "Good", "Fair", "For Parts"]
  });

  const categoryGroup = createFormGroupWithValidation({
    type: "text",
    id: "tool-category",
    name: "category",
    label: "Category / Power Source",
    placeholder: "e.g., Cordless, Pneumatic, Hand Tool"
  });

  conditionalContainer.appendChild(conditionGroup);
  conditionalContainer.appendChild(categoryGroup);

  // Toggle dynamic fields when listing type changes
  const typeSelectInput = typeGroup.querySelector("select");
  if (typeSelectInput) {
    typeSelectInput.addEventListener("change", () => {
      if (typeSelectInput.value === "tool") {
        conditionalContainer.style.display = "block";
      } else {
        conditionalContainer.style.display = "none";
      }
    });
  }

  // --- 3. Media Upload ---
  const imageGroup = createFormGroupWithValidation({
    type: "file",
    id: "product-image",
    name: "images",
    label: "Product Images",
    accept: "image/*",
    multiple: true
  });

  // --- 4. Submit Button ---
  const submitBtn = createElement("button", { type: "submit", class: "btn-submit" }, ["Publish Listing"]);

  // Append everything to form
  form.append(
    typeGroup,
    titleGroup,
    descGroup,
    priceGroup,
    conditionalContainer,
    imageGroup,
    submitBtn
  );

  // --- 5. Form Submission Handler ---
  form.addEventListener("submit", (e) => {
    e.preventDefault();

    // Validate all custom components inside the form
    const validatableElements = form.querySelectorAll<any>(".form-group");
    let isFormValid = true;

    validatableElements.forEach((groupEl) => {
      // Find elements that have our custom `.validate()` attached
      const inputChild = groupEl.querySelector("input, select, textarea");
      if (inputChild && typeof (inputChild as any).validate === "function") {
        const isValid = (inputChild as any).validate();
        if (!isValid) isFormValid = false;
      }
    });

    if (!isFormValid) {
      console.warn("Form validation failed.");
      return;
    }

    // Collect FormData or construct payload
    const formData = new FormData(form);
    const productPayload = Object.fromEntries(formData.entries());

    console.log("Submitting Product/Tool Payload:", productPayload);
    alert("Listing created successfully!");
  });

  container.appendChild(form);
}