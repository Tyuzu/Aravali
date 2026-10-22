import { createFormGroupWithValidation } from "../../components/form/createFormGroupEnhanced.js";
import { createElement } from "../../components/createElement.js";

export function createRecipe(container: HTMLElement, isLoggedIn: boolean): void {
  // Clear container
  container.innerHTML = "";

  if (!isLoggedIn) {
    const notice = createElement("div", { class: "auth-warning" }, [
      "You must be logged in to post and share a recipe."
    ]);
    container.appendChild(notice);
    return;
  }

  const form = createElement("form", { class: "recipe-form" });

  // --- 1. Basic Recipe Details ---
  const titleGroup = createFormGroupWithValidation({
    type: "text",
    id: "recipe-title",
    name: "title",
    label: "Recipe Title",
    required: true,
    placeholder: "e.g., Authentic Spaghetti Carbonara",
    validator: (val) =>
      String(val).trim().length < 3 ? "Title must be at least 3 characters long." : null
  });

  const descGroup = createFormGroupWithValidation({
    type: "textarea",
    id: "recipe-description",
    name: "description",
    label: "Brief Description",
    required: true,
    placeholder: "Give a quick summary of this dish, its flavor profile, or origin..."
  });

  const categoryGroup = createFormGroupWithValidation({
    type: "select",
    id: "recipe-category",
    name: "category",
    label: "Category",
    required: true,
    options: [
      { value: "breakfast", label: "Breakfast & Brunch" },
      { value: "main", label: "Main Course" },
      { value: "side", label: "Side Dish" },
      { value: "dessert", label: "Dessert" },
      { value: "beverage", label: "Beverage" },
      { value: "snack", label: "Snack & Appetizer" }
    ],
    placeholder: "Select a category..."
  });

  // --- 2. Cooking Metadata (Time & Servings) ---
  const prepTimeGroup = createFormGroupWithValidation({
    type: "number",
    id: "prep-time",
    name: "prepTime",
    label: "Prep Time (minutes)",
    required: true,
    placeholder: "15",
    validator: (val) => (Number(val) < 0 ? "Time cannot be negative." : null)
  });

  const cookTimeGroup = createFormGroupWithValidation({
    type: "number",
    id: "cook-time",
    name: "cookTime",
    label: "Cook Time (minutes)",
    required: true,
    placeholder: "30",
    validator: (val) => (Number(val) < 0 ? "Time cannot be negative." : null)
  });

  const servingsGroup = createFormGroupWithValidation({
    type: "number",
    id: "recipe-servings",
    name: "servings",
    label: "Servings",
    required: true,
    placeholder: "4",
    validator: (val) => (Number(val) <= 0 ? "Servings must be at least 1." : null)
  });

  const timeServingsWrapper = createElement("div", { class: "form-row-three" });
  timeServingsWrapper.append(prepTimeGroup, cookTimeGroup, servingsGroup);

  // --- 3. Dynamic Ingredients Section ---
  const ingredientsSection = createElement("div", { class: "recipe-dynamic-section" });
  const ingredientsTitle = createElement("h3", {}, ["Ingredients"]);
  const ingredientsContainer = createElement("div", { class: "ingredients-list" });

  const addIngredientRow = () => {
    const row = createElement("div", { class: "ingredient-row" });

    const ingredientInput = createFormGroupWithValidation({
      type: "text",
      name: "ingredients[]",
      placeholder: "e.g., 200g Guanciale or 2 cups Flour",
      required: true
    });

    const removeBtn = createElement(
      "button",
      { type: "button", class: "btn-remove" },
      ["×"]
    );
    removeBtn.addEventListener("click", () => row.remove());

    row.append(ingredientInput, removeBtn);
    ingredientsContainer.appendChild(row);
  };

  const addIngredientBtn = createElement(
    "button",
    { type: "button", class: "btn-secondary" },
    ["+ Add Ingredient"]
  );
  addIngredientBtn.addEventListener("click", addIngredientRow);

  // Start with one default ingredient input
  addIngredientRow();
  ingredientsSection.append(ingredientsTitle, ingredientsContainer, addIngredientBtn);

  // --- 4. Dynamic Preparation Steps Section ---
  const stepsSection = createElement("div", { class: "recipe-dynamic-section" });
  const stepsTitle = createElement("h3", {}, ["Instructions"]);
  const stepsContainer = createElement("div", { class: "steps-list" });

  const addStepRow = () => {
    const stepCount = stepsContainer.children.length + 1;
    const row = createElement("div", { class: "step-row" });

    const stepInput = createFormGroupWithValidation({
      type: "textarea",
      name: "instructions[]",
      label: `Step ${stepCount}`,
      placeholder: "Describe this step in detail...",
      required: true
    });

    const removeBtn = createElement(
      "button",
      { type: "button", class: "btn-remove" },
      ["Remove Step"]
    );
    removeBtn.addEventListener("click", () => {
      row.remove();
      // Renumber steps
      Array.from(stepsContainer.children).forEach((child, idx) => {
        const label = child.querySelector("label");
        if (label) label.textContent = `Step ${idx + 1}`;
      });
    });

    row.append(stepInput, removeBtn);
    stepsContainer.appendChild(row);
  };

  const addStepBtn = createElement(
    "button",
    { type: "button", class: "btn-secondary" },
    ["+ Add Step"]
  );
  addStepBtn.addEventListener("click", addStepRow);

  // Start with one default step input
  addStepRow();
  stepsSection.append(stepsTitle, stepsContainer, addStepBtn);

  // --- 5. Recipe Photos ---
  const imageGroup = createFormGroupWithValidation({
    type: "file",
    id: "recipe-image",
    name: "images",
    label: "Recipe Cover Photo",
    accept: "image/*",
    multiple: true
  });

  // --- 6. Submit Button ---
  const submitBtn = createElement(
    "button",
    { type: "submit", class: "btn-submit" },
    ["Publish Recipe"]
  );

  // Append all sections to the form
  form.append(
    titleGroup,
    descGroup,
    categoryGroup,
    timeServingsWrapper,
    ingredientsSection,
    stepsSection,
    imageGroup,
    submitBtn
  );

  // --- 7. Validation & Submission Handler ---
  form.addEventListener("submit", (e) => {
    e.preventDefault();

    // Validate all input components in the form
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
      console.warn("Recipe submission blocked due to validation errors.");
      return;
    }

    // Collect dynamic inputs
    const formData = new FormData(form);
    const ingredients = Array.from(
      form.querySelectorAll<HTMLInputElement>("input[name='ingredients[]']")
    ).map((el) => el.value);

    const instructions = Array.from(
      form.querySelectorAll<HTMLTextAreaElement>("textarea[name='instructions[]']")
    ).map((el) => el.value);

    const recipePayload = {
      title: formData.get("title"),
      description: formData.get("description"),
      category: formData.get("category"),
      prepTime: Number(formData.get("prepTime")),
      cookTime: Number(formData.get("cookTime")),
      servings: Number(formData.get("servings")),
      ingredients,
      instructions
    };

    console.log("Submitting Recipe Payload:", recipePayload);
    alert("Recipe published successfully!");
  });

  container.appendChild(form);
}