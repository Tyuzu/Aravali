import { apiFetch } from "../../api/api.js";

export interface FarmDashboardInventory {
  totalCrops?: number;
  totalQuantity?: number;
  inventoryValue?: number;
  featuredCrops?: number;
  lowStockCount?: number;
  outOfStockCount?: number;
  [key: string]: unknown;
}

export interface FarmDashboardOrders {
  total?: number;
  today?: number;
  pending?: number;
  delivered?: number;
  cancelled?: number;
  customers?: number;
  pendingOrders?: number;
  completedOrders?: number;
  cancelledOrders?: number;
  todayDeliveries?: number;
  [key: string]: unknown;
}

export interface FarmDashboardRevenue {
  monthly?: number;
  lifetime?: number;
  month?: number;
  totalRevenue?: number;
  [key: string]: unknown;
}

export interface FarmDashboardStats {
  healthScore?: number;
  [key: string]: unknown;
}

export interface FarmAlertItem {
  type?: string;
  severity?: string;
  message?: string;
  [key: string]: unknown;
}

export interface FarmDashboardResponse {
  success?: boolean;
  message?: string;
  farm?: Record<string, unknown>;
  dashboard?: {
    stats?: FarmDashboardStats;
    inventory?: FarmDashboardInventory;
    orders?: FarmDashboardOrders;
    revenue?: FarmDashboardRevenue;
    alerts?: FarmAlertItem[];
    recommendations?: string[];
    topCrops?: Array<Record<string, unknown>>;
    recentOrders?: Array<Record<string, unknown>>;
    [key: string]: unknown;
  };
  [key: string]: unknown;
}

export async function getFarmDashboard(): Promise<FarmDashboardResponse> {
  return await apiFetch<FarmDashboardResponse>("/dash/farms", "GET");
}

export default {
  getFarmDashboard
};
