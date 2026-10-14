import { AppRoute } from "./routeTypes.js";

export const busRoutes: AppRoute[] = [
  /* =======================================================
     STATIC & OVERVIEW PAGES
  ======================================================= */
  {
    path: "/bus",
    component: () => import("../../pages/bustime/bus.js"),
    functionName: "BusHome",
    meta: {
      title: "Bus Home"
    }
  },
  {
    path: "/bus/time",
    component: () => import("../../pages/bustime/bustime.js"),
    functionName: "BusTime",
    meta: {
      title: "Bus Timetable"
    }
  },
  {
    path: "/bus/map",
    component: () => import("../../pages/bustime/busmap.js"),
    functionName: "BusMap",
    meta: {
      title: "Bus Map"
    }
  },

  /* =======================================================
     DYNAMIC & DETAIL PAGES
  ======================================================= */
  {
    path: "/bus/time/:busid",
    component: () => import("../../pages/bustime/busRouteDetail.js"),
    functionName: "BusRouteDetail",
    meta: {
      title: "Bus Route Details"
    }
  },
  {
    path: "/bus/live/:busid",
    component: () => import("../../pages/bustime/busLiveTracking.js"),
    functionName: "BusLiveTracking",
    meta: {
      title: "Live Bus Tracking"
    }
  },
  {
    path: "/bus/depot/:depotid",
    component: () => import("../../pages/bustime/busDepot.js"),
    functionName: "BusDepot",
    meta: {
      title: "Bus Depot"
    }
  }
];