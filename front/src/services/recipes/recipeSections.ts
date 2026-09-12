import { createElement } from "../../components/createElement.js";
import Button from "../../components/base/Button.js";
import { addToCart, isValidCartQuantity } from "../cart/addToCart.js";
import { getState } from "../../state/state.js";
import { createCommentsSection } from "../comments/comments.js";
import { editRecipe } from "./createOrEditRecipe.js";
import { makeInlineEditable, getStepKey } from "./recipeRenderers.js";
import Notify from "../../components/ui/Notify.js";
import {
  isRecipeOwner,
  normalizeCartQuantity,
  getIngredientItemId
} from "./recipeHelpers.js";
import { Recipe, Ingredient, RecipeStep, User } from "./types/recipe.js";

/* ============================================================
   HELPERS & STORAGE
============================================================ */
const Storage = {
  get<T>(key: string, defaultValue: T): T {
    try {
      const item = localStorage.getItem(key);
      return item ? (JSON.parse(item) as T) : defaultValue;
    } catch (err) {
      console.warn(`[Storage] Failed to read ${key}:`, err);
      return defaultValue;
    }
  },
  set(key: string, value: unknown): boolean {
    try {
      localStorage.setItem(key, JSON.stringify(value));
      return true;
    } catch (err) {
      console.warn(`[Storage] Failed to write ${key}:`, err);
      return false;
    }
  }
};

/* ============================================================
   INGREDIENTS
============================================================ */
export function renderIngredients(
  ingredients?: Ingredient[],
  isLoggedIn?: boolean,
  recipe?: Recipe
): HTMLElement {
  const container = createElement("ul", { class: "ingredients-list" }) as HTMLUListElement;

  if (!ingredients?.length) {
    return createElement("ul", { class: "ingredients-list empty" }, [
      createElement("li", { class: "empty-state" }, ["No ingredients listed for this recipe."]),
    ]);
  }

  const canEdit = isRecipeOwner(recipe);

  ingredients.forEach((ingredient, index) => {
    const itemId = getIngredientItemId(ingredient);
    const quantity = ingredient.quantity ?? "";
    const unit = ingredient.unit ?? "";
    const name = ingredient.name ?? "Unnamed ingredient";

    const labelText = `${quantity} ${unit} ${name}`.replace(/\s+/g, " ").trim();
    const textContainer = createElement("span", { class: "ingredient-text" }, [labelText]);

    const li = createElement("li", { class: "ingredient-item", "data-index": String(index) }, [
      textContainer
    ]);

    // Availability Tag & Quick Add
    if (!itemId) {
      li.appendChild(createElement("span", { class: "badge badge-warning" }, ["Unavailable in store"]));
    } else if (isLoggedIn) {
      const addBtn = createAddToCartButton(ingredient, quantity);
      li.appendChild(addBtn);
    }

    // Recipe Owner Controls
    if (canEdit) {
      const editBtn = Button({
        title: "Edit",
        classes: "tiny-button secondary",
        events: {
          click: (e?: Event) => {
            e?.stopPropagation();
            makeInlineEditable(textContainer, name, (newValue: string) => {
              const cleanValue = newValue.trim();
              ingredient.name = cleanValue;
              textContainer.textContent = `${ingredient.quantity ?? ""} ${ingredient.unit ?? ""} ${cleanValue}`.trim();
            });
          }
        }
      });

      const delBtn = Button({
        title: "Delete",
        classes: "tiny-button danger",
        events: {
          click: (e?: Event) => {
            e?.stopPropagation();
            if (!confirm(`Remove "${name}" from ingredients?`)) return;
            li.remove();
            ingredients.splice(index, 1);
          }
        }
      });

      const actionsGroup = createElement("div", { class: "inline-actions" }, [editBtn, delBtn]);
      li.appendChild(actionsGroup);
    }

    container.appendChild(li);
  });

  return container;
}

function createAddToCartButton(ingredient: Ingredient, recipeQuantity?: number | string): HTMLElement {
  let isAdding = false;

  const btn = Button({
    title: "Add to Cart",
    classes: "small-button primary",
    events: {
      click: async (e?: Event) => {
        e?.stopPropagation();
        if (isAdding) return;

        const itemId = getIngredientItemId(ingredient);
        if (!itemId) {
          Notify("This ingredient is currently unavailable.", { type: "warning" });
          return;
        }

        const quantity = normalizeCartQuantity(recipeQuantity);
        if (!Number.isInteger(quantity) || quantity < 1 || !isValidCartQuantity(quantity)) {
          Notify("Invalid ingredient quantity.", { type: "warning" });
          return;
        }

        const isLoggedIn = Boolean(getState("token"));

        isAdding = true;
        btn.disabled = true;
        btn.textContent = "Adding...";

        try {
          await addToCart({
            itemId,
            itemType: "product",
            quantity,
            isLoggedIn,
            onCartUpdated: (res) => console.debug("Cart updated:", res)
          });
          Notify(`Added ${ingredient.name || "item"} to cart.`, { type: "success", duration: 2000 });
        } catch (error) {
          console.error("Failed to add to cart:", error);
          Notify("Could not add item to cart.", { type: "error" });
        } finally {
          isAdding = false;
          btn.disabled = false;
          btn.textContent = "Add to Cart";
        }
      }
    }
  }) as HTMLButtonElement;

  return btn;
}

/* ============================================================
   STEPS
============================================================ */
export function renderSteps(
  recipeid: string | number,
  steps?: RecipeStep[],
  recipe?: Recipe
): HTMLElement {
  const safeSteps: RecipeStep[] = Array.isArray(steps) ? steps : [];
  const storageKey = getStepKey(recipeid);
  
  const savedIndices = Storage.get<number[]>(storageKey, []);
  let completedSteps = new Set<number>(savedIndices.filter((idx) => idx >= 0 && idx < safeSteps.length));

  // Progress Bar Elements
  const progressFill = createElement("div", { class: "progress-fill" });
  const progressText = createElement("span", { class: "progress-text" });
  const progressBar = createElement(
    "div",
    {
      class: "progress-bar",
      role: "progressbar",
      "aria-valuemin": "0",
      "aria-valuemax": "100",
      "aria-valuenow": "0"
    },
    [progressFill, progressText]
  );

  function syncProgress(): void {
    const total = safeSteps.length;
    const count = completedSteps.size;
    const percentage = total ? Math.round((count / total) * 100) : 0;

    progressFill.style.width = `${percentage}%`;
    progressText.textContent = `${percentage}% completed (${count}/${total})`;
    progressBar.setAttribute("aria-valuenow", String(percentage));
  }

  syncProgress();

  const stepsList = createElement("ol", { class: "recipe-steps-list" });
  const canEdit = isRecipeOwner(recipe);

  safeSteps.forEach((step, index) => {
    const stepText = typeof step === "object" ? step?.text ?? "" : String(step ?? "");
    const isChecked = completedSteps.has(index);

    const checkbox = createElement("input", {
      type: "checkbox",
      id: `step-${recipeid}-${index}`,
      "aria-label": `Mark step ${index + 1} as completed`
    }) as HTMLInputElement;

    checkbox.checked = isChecked;

    const textSpan = createElement(
      "span",
      { class: `step-text ${isChecked ? "completed" : ""}` },
      [stepText]
    );

    checkbox.addEventListener("change", (e: Event) => {
      const checked = (e.target as HTMLInputElement).checked;
      
      if (checked) {
        completedSteps.add(index);
        textSpan.classList.add("completed");
      } else {
        completedSteps.delete(index);
        textSpan.classList.remove("completed");
      }

      Storage.set(storageKey, Array.from(completedSteps));
      syncProgress();
    });

    const li = createElement("li", { class: "step-item" }, [checkbox, textSpan]);

    if (canEdit) {
      const editBtn = Button({
        title: "Edit",
        classes: "tiny-button secondary",
        events: {
          click: (e?: Event) => {
            e?.stopPropagation();
            makeInlineEditable(textSpan, stepText, (newValue: string) => {
              const clean = newValue.trim();
              if (typeof safeSteps[index] === "object") {
                (safeSteps[index] as { text: string }).text = clean;
              } else {
                safeSteps[index] = { text: clean };
              }
              textSpan.textContent = clean;
            });
          }
        }
      });

      const delBtn = Button({
        title: "Delete",
        classes: "tiny-button danger",
        events: {
          click: (e?: Event) => {
            e?.stopPropagation();
            if (!confirm(`Delete step ${index + 1}?`)) return;

            li.remove();
            safeSteps.splice(index, 1);

            // Reindex completed steps
            const reindexed = new Set<number>();
            completedSteps.forEach((idx) => {
              if (idx < index) reindexed.add(idx);
              if (idx > index) reindexed.add(idx - 1);
            });
            completedSteps = reindexed;

            Storage.set(storageKey, Array.from(completedSteps));
            syncProgress();
          }
        }
      });

      li.appendChild(createElement("div", { class: "inline-actions" }, [editBtn, delBtn]));
    }

    stepsList.appendChild(li);
  });

  return createElement("div", { class: "steps-section" }, [progressBar, stepsList]);
}

/* ============================================================
   COMMENTS
============================================================ */
export function renderComments(recipe: Recipe): HTMLElement {
  const wrapper = createElement("div", { class: "recipe-comments" });
  const heading = createElement("h4", {}, ["Comments"]);
  
  const toggle = createElement(
    "button",
    {
      type: "button",
      class: "toggle-comments btn btn-link",
      "aria-expanded": "false"
    },
    ["💬 Show Comments"]
  ) as HTMLButtonElement;

  let commentsEl: HTMLElement | null = null;
  let isLoading = false;
  let isVisible = false;

  toggle.addEventListener("click", async () => {
    if (commentsEl) {
      isVisible = !isVisible;
      commentsEl.style.display = isVisible ? "block" : "none";
      toggle.textContent = isVisible ? "💬 Hide Comments" : "💬 Show Comments";
      toggle.setAttribute("aria-expanded", String(isVisible));
      return;
    }

    if (isLoading) return;
    isLoading = true;
    toggle.disabled = true;
    toggle.textContent = "Loading comments...";

    try {
      const user = getState("user") as User | undefined;
      commentsEl = await createCommentsSection("recipe", recipe.recipeid, user?.userid);
      
      if (!commentsEl) throw new Error("Comments container unavailable.");

      wrapper.appendChild(commentsEl);
      isVisible = true;
      toggle.textContent = "💬 Hide Comments";
      toggle.setAttribute("aria-expanded", "true");
    } catch (error) {
      console.error("Failed to load comments:", error);
      Notify("Failed to load comments.", { type: "error" });
      toggle.textContent = "💬 Show Comments";
    } finally {
      isLoading = false;
      toggle.disabled = false;
    }
  });

  wrapper.append(heading, toggle);
  return wrapper;
}

/* ============================================================
   ACTIONS
============================================================ */
export function renderActions(
  recipe: Recipe,
  currentUser: User | string | number | null,
  contentContainer: HTMLElement,
  isFavorite: boolean,
  recipeid: string | number
): HTMLElement {
  let favoriteState = isFavorite;

  const favBtn = Button({
    title: favoriteState ? "❤️ Saved" : "🤍 Save Recipe",
    classes: "buttonx secondary",
    events: {
      click: () => {
        const key = "favoriteRecipes";
        const favorites = Storage.get<string[]>(key, []);
        const targetId = String(recipeid);

        let updated: string[];
        if (favoriteState) {
          updated = favorites.filter((id) => id !== targetId);
          favoriteState = false;
        } else {
          updated = Array.from(new Set([...favorites, targetId]));
          favoriteState = true;
        }

        if (Storage.set(key, updated)) {
          favBtn.textContent = favoriteState ? "❤️ Saved" : "🤍 Save Recipe";
          Notify(favoriteState ? "Recipe saved to favorites." : "Recipe removed from favorites.", {
            type: "info",
            duration: 2000
          });
        } else {
          Notify("Unable to update favorites.", { type: "warning" });
        }
      }
    }
  });

  const shareBtn = Button({
    title: "🔗 Copy Link",
    classes: "buttonx secondary",
    events: {
      click: async () => {
        try {
          if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(window.location.href);
          } else {
            // Fallback for older browsers / non-secure contexts
            const textarea = document.createElement("textarea");
            textarea.value = window.location.href;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand("copy");
            document.body.removeChild(textarea);
          }
          shareBtn.textContent = "✓ Copied!";
          setTimeout(() => (shareBtn.textContent = "🔗 Copy Link"), 2000);
          Notify("Recipe link copied to clipboard.", { type: "success", duration: 2000 });
        } catch (error) {
          console.error("Failed to copy link:", error);
          Notify("Unable to copy link.", { type: "warning" });
        }
      }
    }
  });

  const printBtn = Button({
    title: "🖨️ Print",
    classes: "buttonx secondary",
    events: { click: () => window.print() }
  });

  const actions: HTMLElement[] = [favBtn, shareBtn, printBtn];

  // Resolve ownership safely
  const currentUserId = typeof currentUser === "object" && currentUser ? currentUser.id ?? currentUser.userid : currentUser;
  const isOwner = Boolean(currentUserId && recipe?.userid && String(currentUserId) === String(recipe.userid));

  if (isOwner) {
    const editBtn = Button({
      title: "✏️ Edit Recipe",
      classes: "buttonx secondary",
      events: { click: () => editRecipe(contentContainer, recipe) }
    });
    actions.push(editBtn);
  }

  const backBtn = Button({
    title: "← Back",
    classes: "buttonx primary",
    events: { click: () => history.back() }
  });
  actions.push(backBtn);

  return createElement("div", { class: "recipe-actions" }, actions);
}