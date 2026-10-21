import { AppRoute } from "./routeTypes.js";

export const userRoutes: AppRoute[] = [


  {
    path: "/login",
    component: () => import("../../pages/auth/auth.js"),
    functionName: "Auth",
    meta: {
      guestOnly: true,
      title: "Login"
    }
  },

  {
    path: "/profile",
    component: () => import("../../pages/profile/userProfile.js"),
    functionName: "MyProfile",
    meta: {
      requiresAuth: true,
      title: "My Profile"
    }
  },



  /* =======================================================
     DYNAMIC ROUTES
  ======================================================= */

  {
    path: "/user/:id",
    component: () => import("../../pages/profile/userProfile.js"),
    functionName: "UserProfile",
    meta: {
      title: "User Profile"
    }
  },

];