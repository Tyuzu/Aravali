import { createFormGroupWithValidation } from "../../components/form/createFormGroupEnhanced.js";
import { createElement } from "../../components/createElement.js";

export function createFarm(isLoggedIn: boolean, contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  if (!isLoggedIn) {
    const note = document.createElement("div");
    note.className = "empty-state";
    note.textContent = "Please sign in to create a farm profile.";
    contentContainer.appendChild(note);
    return;
  }

  const form = createElement("form", { class: "farm-profile-form" });

  // --- 1. Basic Farm Information ---
  const nameGroup = createFormGroupWithValidation({
    type: "text",
    id: "farm-name",
    name: "farmName",
    label: "Farm Name",
    required: true,
    placeholder: "e.g., Green Valley Organic Farm",
    validator: (val) =>
      String(val).trim().length < 3 ? "Farm name must be at least 3 characters long." : null
  });

  const typeGroup = createFormGroupWithValidation({
    type: "select",
    id: "farm-type",
    name: "farmType",
    label: "Primary Farming Type",
    required: true,
    options: [
      { value: "crop", label: "Crop & Produce" },
      { value: "livestock", label: "Livestock & Dairy" },
      { value: "mixed", label: "Mixed Farming" },
      { value: "orchard", label: "Orchard & Vineyard" },
      { value: "greenhouse", label: "Greenhouse / Controlled Environment" }
    ],
    placeholder: "Select farm type..."
  });

  const sizeGroup = createFormGroupWithValidation({
    type: "number",
    id: "farm-size",
    name: "acreage",
    label: "Total Acreage",
    required: true,
    placeholder: "e.g., 50",
    validator: (val) => (Number(val) <= 0 ? "Acreage must be greater than zero." : null)
  });

  const locationGroup = createFormGroupWithValidation({
    type: "text",
    id: "farm-location",
    name: "location",
    label: "Location / Region",
    required: true,
    placeholder: "e.g., Sonoma County, CA"
  });

  // --- 2. Practices & Certification ---
  const practicesGroup = createFormGroupWithValidation({
    type: "select",
    id: "farming-practices",
    name: "practices",
    label: "Farming Practices",
    multiple: true,
    options: [
      "Certified Organic",
      "Regenerative Agriculture",
      "Permaculture",
      "Hydroponic / Aquaponic",
      "Conventional",
      "Pesticide-Free"
    ],
    placeholder: "Select practices..."
  });

  const bioGroup = createFormGroupWithValidation({
    type: "textarea",
    id: "farm-bio",
    name: "bio",
    label: "Farm Story / Description",
    required: true,
    placeholder: "Tell customers or partners about your history, mission, and what makes your farm unique..."
  });

  // --- 3. Farm Gallery / Photos ---
  const imageGroup = createFormGroupWithValidation({
    type: "file",
    id: "farm-images",
    name: "images",
    label: "Farm Photos",
    accept: "image/*",
    multiple: true
  });

  // --- 4. Submit Button ---
  const submitBtn = createElement(
    "button",
    { type: "submit", class: "btn-submit" },
    ["Create Farm Profile"]
  );

  // Append all groups to the form
  form.append(
    nameGroup,
    typeGroup,
    sizeGroup,
    locationGroup,
    practicesGroup,
    bioGroup,
    imageGroup,
    submitBtn
  );

  // --- 5. Submission & Validation Handler ---
  form.addEventListener("submit", (e) => {
    e.preventDefault();

    const groups = form.querySelectorAll<HTMLElement>(".form-group");
    let isFormValid = true;

    groups.forEach((groupEl) => {
      const inputChild = groupEl.querySelector("input, select, textarea");
      if (inputChild && typeof (inputChild as any).validate === "function") {
        const isValid = (inputChild as any).validate();
        if (!isValid) isFormValid = false;
      }
    });

    if (!isFormValid) {
      console.warn("Farm profile submission blocked due to validation errors.");
      return;
    }

    const formData = new FormData(form);
    const farmPayload = Object.fromEntries(formData.entries());

    console.log("Submitting Farm Profile:", farmPayload);
    alert("Farm profile successfully created!");
  });

  contentContainer.appendChild(form);
}