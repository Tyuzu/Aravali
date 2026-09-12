import { createElement } from "../../components/createElement.js";
import { getState } from "../../state/state.js";
import { fetchRecipeById } from "./api.js";
import { displayMedia } from "../media/ui/mediaGallery.js";

import {
  getFavorites,
  renderAuthor,
  createRecipeBannerSection,
  renderInfoBox,
  renderTags
} from "./recipeRenderers.js";

import {
  renderIngredients,
  renderSteps,
  renderComments,
  renderActions
} from "./recipeSections.js";

import { Recipe, User } from "./types/recipe.js";

/* =========================
   MAIN DISPLAY
========================= */

export async function displayRecipe(
  content: HTMLElement,
  isLoggedIn: boolean,
  recipeid: string | number
): Promise<void> {
  content.replaceChildren();

  const container = createElement("div", { class: "recipe-page single-page-layout" });
  content.appendChild(container);

  const currentUser = (getState("user") as User | undefined)?.userid;

  let recipe: Recipe;

  try {
    recipe = await fetchRecipeById(recipeid);
  } catch {
    container.replaceChildren(
      createElement("p", { class: "error-message" }, ["Recipe not found or failed to load."])
    );
    return;
  }

  const isFavorite = getFavorites().map(String).includes(String(recipeid));

  /* HEADER & METADATA */
  const titleEl = createElement("h2", { class: "recipe-title" }, [
    recipe.title || recipe.name || "Untitled"
  ]);

  const metaInfo: HTMLElement[] = [];

  if (recipe.version) {
    metaInfo.push(
      createElement("p", { class: "version-info" }, [`Version ${recipe.version}`])
    );
  }

  if (recipe.lastUpdated) {
    metaInfo.push(
      createElement("p", { class: "version-info" }, [
        `Last updated: ${new Date(recipe.lastUpdated).toLocaleDateString()}`
      ])
    );
  }

  const authorEl = renderAuthor(recipe, currentUser);

  /* BANNER, INFO & TAGS */
  const bannerEl = createRecipeBannerSection(recipe, currentUser);
  const infoBox = renderInfoBox(recipe);
  const tagsEl = renderTags(recipe.tags);

  /* SECTIONS (Replacing Tabs) */

  // Top Sticky / Quick Actions Bar
  const actionsSection = createElement("section", { class: "recipe-section actions-section" }, [
    renderActions(recipe, getState("user") as User, content, isFavorite, recipeid)
  ]);

  // Ingredients Section
  const ingredientsSection = createElement("section", { class: "recipe-section ingredients-section" }, [
    createElement("h3", { class: "section-title" }, ["Ingredients"]),
    renderIngredients(recipe.ingredients, isLoggedIn, recipe)
  ]);

  // Preparation Steps Section
  const stepsSection = createElement("section", { class: "recipe-section steps-section" }, [
    createElement("h3", { class: "section-title" }, ["Instructions"]),
    renderSteps(recipeid, recipe.steps || [], recipe)
  ]);

  // Media Gallery Container
  const mediaContainer = createElement("div", { class: "media-gallery-wrapper" });
  displayMedia(mediaContainer, "recipe", recipeid, isLoggedIn);

  const mediaSection = createElement("section", { class: "recipe-section media-section" }, [
    createElement("h3", { class: "section-title" }, ["Photos & Media"]),
    mediaContainer
  ]);

  // Comments Section
  const commentsSection = createElement("section", { class: "recipe-section comments-section" }, [
    renderComments(recipe)
  ]);

  /* FINAL ASSEMBLY */
  container.replaceChildren(
    titleEl,
    ...metaInfo,
    authorEl,
    bannerEl,
    infoBox,
    tagsEl,
    actionsSection,
    ingredientsSection,
    stepsSection,
    mediaSection,
    commentsSection
  );
}