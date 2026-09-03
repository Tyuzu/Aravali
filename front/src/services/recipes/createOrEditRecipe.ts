import { createElement } from "../../components/createElement.js";
import Button from "../../components/base/Button.js";
import { createFormGroup } from "../../components/form/createFormGroupEnhanced.js";
import { navigate } from "../../routes/navigate.js";
import { saveRecipeRequest } from "./api.js";
import { Recipe, IngredientAlternative } from "./types/recipe.js";
import Modal from "../../components/ui/Modal.js";
import { fetchFarmItems, fetchFarmCategories } from "../products/api.js";
import Imagex from "../../components/base/Imagex.js";
import { resolveImagePath, EntityType, PictureType } from "../../utils/imagePaths.js";

let ingredientPickerCounter = 0;

type FormMode = "create" | "edit";

export function createRecipe(container: HTMLElement): void {
  renderRecipeForm(container, "create", null);
}

export function editRecipe(container: HTMLElement, recipe: Recipe): void {
  renderRecipeForm(container, "edit", recipe);
}

function renderRecipeForm(
  container: HTMLElement,
  mode: FormMode = "create",
  recipe: Recipe | null = null
): void {
  container.replaceChildren();

  const form = createElement("form", {
    class: "create-section",
    enctype: "multipart/form-data",
  }) as HTMLFormElement;

  const titleGroup = createFormGroup({
    label: "Recipe Title",
    type: "text",
    id: "title",
    value: recipe?.title || "",
    placeholder: "Enter recipe title",
    required: true,
  });

  const descriptionGroup = createFormGroup({
    label: "Description",
    type: "textarea",
    id: "description",
    value: recipe?.description || "",
    placeholder: "Short summary of the dish",
    required: true,
    additionalProps: { rows: 3 },
  });

  const cuisineGroup = createFormGroup({
    label: "Cuisine",
    type: "text",
    id: "cuisine",
    value: recipe?.cuisine || "",
    placeholder: "e.g. Italian, Indian, Mexican",
  });

  const dietaryGroup = createFormGroup({
    label: "Dietary / Allergen Info (comma-separated)",
    type: "text",
    id: "dietary",
    value: (recipe?.dietary || []).join(", "),
    placeholder: "e.g. Vegan, Gluten-Free, Nut-Free",
  });

  const cookTimeGroup = createFormGroup({
    label: "Cook Time",
    type: "text",
    id: "cookTime",
    value: recipe?.cookTime || "",
    placeholder: "e.g. 30 mins",
  });

  const servingsGroup = createFormGroup({
    label: "Servings",
    type: "number",
    id: "servings",
    value: recipe?.servings !== undefined ? String(recipe.servings) : "",
    placeholder: "e.g. 4",
    additionalProps: { min: 1 },
  });

  const portionGroup = createFormGroup({
    label: "Portion Size / Servings Scaling",
    type: "text",
    id: "portionSize",
    value: recipe?.portionSize || "",
    placeholder: "e.g. 1 cup per person",
  });

  const seasonGroup = createFormGroup({
    label: "Season / Occasion",
    type: "text",
    id: "season",
    value: recipe?.season || "",
    placeholder: "e.g. Winter, Christmas",
  });

  const tagsGroup = createFormGroup({
    label: "Tags (comma-separated)",
    type: "text",
    id: "tags",
    value: (recipe?.tags || []).join(", "),
    placeholder: "e.g. spicy, vegan, south indian",
  });

  const difficultyGroup = createFormGroup({
    label: "Difficulty",
    type: "select",
    id: "difficulty",
    value: recipe?.difficulty || "",
    options: [
      { value: "", label: "Select difficulty" },
      { value: "Easy", label: "Easy" },
      { value: "Medium", label: "Medium" },
      { value: "Hard", label: "Hard" },
    ],
  });

  // --- Ingredients ---
  const ingredientsGroup = createElement("div", { class: "form-group" });
  ingredientsGroup.appendChild(createElement("label", {}, ["Ingredients"]));
  const ingredientsList = createElement("div", { id: "ingredients-list" });
  ingredientsGroup.appendChild(ingredientsList);

  function addIngredientRow(
    name: string = "",
    quantity: number | string = "",
    unit: string = "",
    alternatives: IngredientAlternative[] = [],
    linkedItemId: string | null = null,
    linkedItemType: string | null = null
  ): void {
    const altStr = alternatives
      .map((a) => [a.name, a.itemId, a.type].join("|"))
      .join(",");

    const removeBtn = Button({
      title: "−",
      classes: "remove-btn",
      type: "button",
      events: {
        click: (e: Event) => {
          e.preventDefault();
          row.remove();
        },
      },
    });

    const hiddenItemId = createElement("input", { type: "hidden", name: "ingredientItemId[]", value: linkedItemId || "" });
    const hiddenItemType = createElement("input", { type: "hidden", name: "ingredientType[]", value: linkedItemType || "" });

    const linkedLabel = createElement("span", { class: "linked-item" }, [linkedItemId ? "Linked" : "Not linked"]);

    const linkBtn = Button({
      title: "Link Item",
      type: "button",
      classes: "small-button",
      events: {
        click: async (e: Event) => {
          e.preventDefault();
          const pickId = `ingredient-picker-${++ingredientPickerCounter}`;
          const searchInput = createElement("input", { type: "search", placeholder: "Search products...", id: `${pickId}-input`, autocomplete: "off" }) as HTMLInputElement;
          const categorySelect = createElement("select", { id: `${pickId}-category` }) as HTMLSelectElement;
          const results = createElement("div", { class: "picker-results" });

          // populate categories
          try {
            const cats = await fetchFarmCategories("product");
            categorySelect.appendChild(createElement("option", { value: "" }, ["All categories"]));
            cats.forEach((c: string) => categorySelect.appendChild(createElement("option", { value: c }, [c])));
          } catch (err) {
            categorySelect.appendChild(createElement("option", { value: "" }, ["All categories"]));
          }
          let typingTimer: number | undefined;
          let searchSeq = 0;
          let active = true;

          const modalRefHolder: { ref?: any } = {};

          const doSearch = async () => {
            window.clearTimeout(typingTimer);
            typingTimer = window.setTimeout(async () => {
              const q = searchInput.value.trim();
              const cat = (categorySelect as HTMLSelectElement).value || "";
              results.replaceChildren();
              if (!q && !cat) return;
              const seq = ++searchSeq;
              try {
                const resp = await fetchFarmItems("product", { search: q, limit: 20, category: cat });
                if (!active || seq !== searchSeq) return; // ignore out-of-order results
                const items = resp.items || [];
                if (items.length === 0) {
                  results.appendChild(createElement("p", {}, ["No products found."]));
                  return;
                }
                items.forEach((it: any) => {
                  const thumbSrc = resolveImagePath(EntityType.PRODUCT, PictureType.THUMB, it.banner || (Array.isArray(it.images) ? it.images[0] : it.images));
                  const thumb = Imagex({ src: thumbSrc, alt: it.name || "", classes: "picker-thumb" });
                  const meta = createElement("div", { class: "picker-meta" }, [createElement("div", { class: "picker-name" }, [it.name || String(it.productid || it.productId || "")]), it.category ? createElement("div", { class: "picker-cat" }, [it.category]) : null].filter(Boolean) as HTMLElement[]);
                  const itemBtn = createElement("button", { type: "button", class: "picker-item" }, [thumb, meta]) as HTMLButtonElement;
                  itemBtn.addEventListener("click", (ev: Event) => {
                    ev.preventDefault();
                    (hiddenItemId as HTMLInputElement).value = it.productid || it.productId || it.id || "";
                    (hiddenItemType as HTMLInputElement).value = "product";
                    linkedLabel.replaceChildren(it.name || String(it.productid || it.productId || ""));
                    cleanupAndClose();
                  });
                  results.appendChild(itemBtn);
                });
              } catch (err) {
                if (!active) return;
                results.appendChild(createElement("p", {}, ["Search failed"]));
              }
            }, 250);
          };

          const onInput = () => doSearch();
          const onCategoryChange = () => doSearch();

          const cleanupAndClose = () => {
            active = false;
            window.clearTimeout(typingTimer);
            searchInput.removeEventListener("input", onInput);
            categorySelect.removeEventListener("change", onCategoryChange);
            try {
              modalRefHolder.ref?.close();
            } catch (_) {}
          };

          const modalRef = Modal({
            title: "Select Product",
            content: createElement("div", {}, [searchInput, categorySelect, results]),
            size: "medium",
            actions: () => Button({ title: "Close", type: "button", events: { click: () => cleanupAndClose() } })
          });
          modalRefHolder.ref = modalRef;

          searchInput.addEventListener("input", onInput);
          categorySelect.addEventListener("change", onCategoryChange);
        }
      }
    }) as HTMLButtonElement;

    const row = createElement("div", { class: "ingredient-row hvflex" }, [
      createElement("input", {
        type: "text",
        name: "ingredientName[]",
        placeholder: "Name",
        value: name,
        required: true,
      }),
      createElement("input", {
        type: "number",
        name: "ingredientQuantity[]",
        placeholder: "Qty",
        step: "any",
        value: quantity !== "" ? String(quantity) : "",
        required: true,
      }),
      createElement("input", {
        type: "text",
        name: "ingredientUnit[]",
        placeholder: "Unit",
        value: unit,
        required: true,
      }),
      createElement("input", {
        type: "text",
        name: "ingredientAlternatives[]",
        placeholder: "Alternatives (name|itemId|type, ...)",
        value: altStr,
      }),
      hiddenItemId,
      hiddenItemType,
      linkedLabel,
      linkBtn,
      removeBtn,
    ]);

    ingredientsList.appendChild(row);
  }

  if (recipe?.ingredients?.length) {
    recipe.ingredients.forEach((ing) =>
      addIngredientRow(
        ing.name,
        ing.quantity,
        ing.unit,
        ing.alternatives || [],
        (ing.itemId as unknown as string) || null,
        (ing.type as unknown as string) || null
      )
    );
  } else {
    addIngredientRow();
  }

  const addIngredientBtn = Button({
    title: "Add Ingredient",
    type: "button",
    events: {
      click: (e: Event) => {
        e.preventDefault();
        addIngredientRow();
      },
    },
  });

  ingredientsGroup.appendChild(addIngredientBtn);

  // --- Steps ---
  const initialStepsText = (recipe?.steps || [])
    .map((step) => (typeof step === "object" ? step.text : step))
    .join("\n");

  const stepsGroup = createFormGroup({
    label: "Steps",
    type: "textarea",
    id: "steps",
    value: initialStepsText,
    placeholder: "Each step on a new line",
    required: true,
    additionalProps: { rows: 6 },
  });

  // --- Video URL ---
  const videoGroup = createFormGroup({
    label: "Video URL / Cooking Tutorial",
    type: "url",
    id: "videoUrl",
    value: recipe?.videoUrl || "",
    placeholder: "e.g. https://www.youtube.com/...",
  });

  // --- Notes ---
  const notesGroup = createFormGroup({
    label: "Recipe Notes / Tips",
    type: "textarea",
    id: "notes",
    value: recipe?.notes || "",
    placeholder: "Extra tips, variations, or notes",
    additionalProps: { rows: 3 },
  });

  // --- Submit Button ---
  const submitBtn = Button({
    title: mode === "edit" ? "Update Recipe" : "Create Recipe",
    type: "submit",
  });

  form.addEventListener("submit", async (e: SubmitEvent) => {
    e.preventDefault();
    const formData = new FormData(form);
    const endpoint =
      mode === "edit" ? `/recipes/recipe/${recipe?.recipeid}` : "/recipes";
    const method = mode === "edit" ? "PUT" : "POST";
    const submitBtnEl = submitBtn as HTMLButtonElement;
    const originalText = submitBtnEl.textContent || submitBtnEl.innerText || "";
    try {
      submitBtnEl.disabled = true;
      form.setAttribute("aria-busy", "true");
      submitBtnEl.textContent = "Saving...";

      const result = await saveRecipeRequest(formData, mode, recipe?.recipeid);
      if (mode === "create") {
        form.reset();
      }

      const modalRefHolder: { ref?: any } = {};
      const modalRef = Modal({
        title: "Recipe Saved",
        content: createElement("div", {}, [createElement("p", {}, ["Recipe saved successfully!"])]),
        size: "small",
        actions: () => Button({ title: "OK", type: "button", events: { click: () => { try { modalRefHolder.ref?.close(); } catch (_) {}; navigate(`/recipe/${result?.recipeid}`); } } })
      });
      modalRefHolder.ref = modalRef;
    } catch (err) {
      console.error("Upload failed:", err);
      const modalRef = Modal({
        title: "Save Failed",
        content: createElement("div", {}, [createElement("p", {}, ["Failed to save recipe. Please try again."])]),
        size: "small",
        actions: () => Button({ title: "Close", type: "button", events: { click: () => modalRef?.close() } })
      });
      // re-enable on failure
      submitBtnEl.disabled = false;
      form.removeAttribute("aria-busy");
      submitBtnEl.textContent = originalText;
    }
  });

  form.append(
    titleGroup,
    descriptionGroup,
    cuisineGroup,
    dietaryGroup,
    cookTimeGroup,
    servingsGroup,
    portionGroup,
    seasonGroup,
    tagsGroup,
    difficultyGroup,
    ingredientsGroup,
    stepsGroup,
    videoGroup,
    notesGroup,
    submitBtn
  );

  container.appendChild(form);
}