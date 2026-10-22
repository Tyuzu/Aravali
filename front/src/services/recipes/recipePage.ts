// displayRecipe.ts
import { createElement } from "../../components/createElement.js";

// Mock data fetcher for demonstration
interface Recipe {
  id: string | number;
  title: string;
  category: string;
  prepTime: number;
  cookTime: number;
  servings: number;
  description: string;
  ingredients: string[];
  instructions: string[];
  author: string;
}

const mockRecipesDatabase: Record<string | number, Recipe> = {
  1: {
    id: 1,
    title: "Classic Margherita Pizza",
    category: "Main Course",
    prepTime: 20,
    cookTime: 15,
    servings: 2,
    description: "A traditional Neapolitan pizza topped with fresh mozzarella, tomatoes, and basil.",
    ingredients: ["1 pizza dough ball", "1/2 cup San Marzano tomato sauce", "8 oz fresh mozzarella", "Fresh basil leaves", "2 tbsp extra virgin olive oil"],
    instructions: ["Preheat oven to 500°F (260°C).", "Roll out the pizza dough on a floured surface.", "Spread tomato sauce evenly and top with sliced mozzarella.", "Bake for 12-15 minutes until crust is golden.", "Garnish with fresh basil and olive oil before serving."],
    author: "Chef Mario"
  }
  // Additional mock recipes can be added here...
};

export function displayRecipe(container: HTMLElement, isLoggedIn: boolean, recipeid: string | number): void {
  container.replaceChildren();

  const recipe = mockRecipesDatabase[recipeid] || {
    id: recipeid,
    title: `Recipe #${recipeid}`,
    category: "Gourmet",
    prepTime: 15,
    cookTime: 30,
    servings: 4,
    description: "Delicious homemade dish prepared with fresh ingredients.",
    ingredients: ["Ingredient A", "Ingredient B", "Ingredient C"],
    instructions: ["Prepare components.", "Cook according to directions.", "Serve warm."],
    author: "Home Chef"
  };

  const wrapper = createElement("div", { class: "recipe-detail-container" });

  // Header Section
  const header = createElement("div", { class: "recipe-detail-header" });
  const title = createElement("h1", { class: "recipe-title" }, [recipe.title]);
  const meta = createElement("div", { class: "recipe-meta-badges" }, [
    createElement("span", { class: "badge" }, [recipe.category]),
    createElement("span", { class: "badge" }, [`Prep: ${recipe.prepTime}m`]),
    createElement("span", { class: "badge" }, [`Cook: ${recipe.cookTime}m`]),
    createElement("span", { class: "badge" }, [`Servings: ${recipe.servings}`])
  ]);
  header.append(title, meta);

  // Description & Author
  const desc = createElement("p", { class: "recipe-description" }, [recipe.description]);
  const authorInfo = createElement("div", { class: "recipe-author" }, [`Recipe by ${recipe.author}`]);

  // Ingredients List
  const ingredientsSection = createElement("div", { class: "recipe-section" });
  ingredientsSection.appendChild(createElement("h3", {}, ["Ingredients"]));
  const ingList = createElement("ul", { class: "ingredient-bullet-list" });
  recipe.ingredients.forEach(ing => {
    ingList.appendChild(createElement("li", {}, [ing]));
  });
  ingredientsSection.appendChild(ingList);

  // Instructions Steps
  const instructionsSection = createElement("div", { class: "recipe-section" });
  instructionsSection.appendChild(createElement("h3", {}, ["Instructions"]));
  const instList = createElement("ol", { class: "instruction-numbered-list" });
  recipe.instructions.forEach(step => {
    instList.appendChild(createElement("li", {}, [step]));
  });
  instructionsSection.appendChild(instList);

  // Action Buttons (Conditional on login)
  const actions = createElement("div", { class: "recipe-actions" });
  const backBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["← Back to Recipes"]);
  backBtn.addEventListener("click", () => history.back());
  actions.appendChild(backBtn);

  if (isLoggedIn) {
    const saveBtn = createElement("button", { type: "button", class: "btn-submit" }, ["Save to Favorites"]);
    saveBtn.addEventListener("click", () => alert("Recipe saved to your favorites!"));
    actions.appendChild(saveBtn);
  }

  wrapper.append(header, desc, authorInfo, ingredientsSection, instructionsSection, actions);
  container.appendChild(wrapper);
}