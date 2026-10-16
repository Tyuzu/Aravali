import { AppRoute } from "./routeTypes.js";

export const cropRoutes: AppRoute[] = [
  {
    path: "/crops",
    component: () =>
      import(
        "../../pages/crop/crops.js"
      ),
    functionName: "Crops",
    meta: {
      title: "Crops"
    }
  },

  {
    path: "/bazarbhav",
    component: () =>
      import(
        "../../pages/crop/bazaar.js"
      ),
    functionName: "BazaarBhav",
    meta: {
      title: "Bazaar Bhav"
    }
  },

  {
    path: "/grocery",
    component: () =>
      import(
        "../../pages/crop/crops.js"
      ),
    functionName: "Crops",
    meta: {
      title: "Grocery"
    }
  },

  {
    path: "/crop/:id",
    component: () =>
      import(
        "../../pages/crop/cropPage.js"
      ),
    functionName: "Crop",
    meta: {
      title: "Crop"
    }
  },

  {
    path: "/aboutcrop/:id",
    component: () =>
      import(
        "../../pages/crop/aboutCropPage.js"
      ),
    functionName: "AboutCrop",
    meta: {
      title: "About Crop"
    }
  }
];
