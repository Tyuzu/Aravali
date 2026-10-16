import { AppRoute } from "./routeTypes.js";

export const farmsRoutes: AppRoute[] = [
  /* =======================================================
     DASHBOARD
  ======================================================= */

  {
    path: "/dash",
    component: () =>
      import(
        "../../pages/dash/dash.js"
      ),
    functionName: "Dash",
    meta: {
      requiresAuth: true,
      roles: ["farmer", "admin"],
      roleMatchMode: "ANY",
      title: "Dashboard"
    }
  },

  /* =======================================================
     FARMS
  ======================================================= */

  {
    path: "/farms",
    component: () =>
      import(
        "../../pages/farm/farms.js"
      ),
    functionName: "Farms",
    meta: {
      title: "Farms"
    }
  },

  {
    path: "/create-farm",
    component: () =>
      import(
        "../../pages/farm/createNewFarm.js"
      ),
    functionName: "CreateFarm",
    meta: {
      requiresAuth: true,
      title: "Create Farm"
    }
  },

  /* =======================================================
     FARM TOOLS
  ======================================================= */

  {
    path: "/tools",
    component: () =>
      import(
        "../../pages/farm/tools.js"
      ),
    functionName: "Tools",
    meta: {
      title: "Tools"
    }
  },

  {
    path: "/products",
    component: () =>
      import(
        "../../pages/farm/products.js"
      ),
    functionName: "Products",
    meta: {
      title: "Products"
    }
  },

  {
    path: "/recipes",
    component: () =>
      import(
        "../../pages/recipe/recipes.js"
      ),
    functionName: "Recipes",
    meta: {
      title: "Recipes"
    }
  },

  {
    path: "/create-recipe",
    component: () =>
      import(
        "../../pages/recipe/createNewRecipe.js"
      ),
    functionName: "CreateRecipe",
    meta: {
      requiresAuth: true,
      title: "Create Recipe"
    }
  },

  /* =======================================================
     DYNAMIC ROUTES
  ======================================================= */

  {
    path: "/products/:type/:id",
    component: () =>
      import(
        "../../pages/product/product.js"
      ),
    functionName: "Product",
    meta: {
      title: "Product"
    }
  },

  {
    path: "/farm/:id",
    component: () =>
      import(
        "../../pages/farm/displayFarm.js"
      ),
    functionName: "Farm",
    meta: {
      title: "Farm"
    }
  },

  {
    path: "/recipe/:id",
    component: () =>
      import(
        "../../pages/recipe/recipePage.js"
      ),
    functionName: "Recipe",
    meta: {
      title: "Recipe"
    }
  }
];