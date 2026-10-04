export interface IngredientAlternative {
  name: string;
  itemId?: string | number;
  itemid?: string | number;
  type?: string;
  [key: string]: unknown;
}

export interface Ingredient {
  name?: string;
  quantity?: number | string;
  unit?: string;
  itemId?: string | number;
  itemid?: string | number;
  productId?: string | number;
  productid?: string | number;
  type?: string;
  alternatives?: IngredientAlternative[];
  [key: string]: unknown;
}

export interface RecipeStepObject {
  text: string;
  [key: string]: unknown;
}

export type RecipeStep = string | RecipeStepObject;

export interface User {
  id?: string | number;
  userid?: string | number;
  userId?: string | number;
  username?: string;
  [key: string]: unknown;
}

export interface Recipe {
  recipeid: string | number;
  recipeId?: string | number;
  userid?: string | number;
  userId?: string | number;
  userID?: string | number;
  title?: string;
  name?: string;
  description?: string;
  ingredients?: Ingredient[];
  steps?: RecipeStep[];
  cookTime?: string;
  cooktime?: string;
  servings?: number | string;
  cuisine?: string;
  portionSize?: string;
  portionsize?: string;
  season?: string;
  dietary?: string[];
  tags?: string[];
  images?: string[];
  difficulty?: "Easy" | "Medium" | "Hard" | "";
  videoUrl?: string;
  videourl?: string;
  notes?: string;
  banner?: string;
  version?: string | number;
  lastUpdated?: string | number | Date;
  updatedAt?: string | number | Date;
  createdAt?: string | number | Date;
  views?: number;
  username?: string;
  [key: string]: unknown;
}

export interface RecipeListResponse {
  recipes?: Recipe[];
  data?: Recipe[];
  [key: string]: unknown;
}

export interface TabItem {
  title: string;
  id: string;
  render: (container: HTMLElement) => void;
}