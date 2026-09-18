import { createElement } from "../../components/createElement.js";
import { buildOrdersPage } from "./orders/builders.js";
import { normalizeOrders } from "./orders/orderutils.js";
import { OrderPageState } from "./orders/types.js";
import { getMyOrders } from "./api.js";

/**
 * Renders and coordinates the User Orders page.
 * @param container - Target parent node element wrapper.
 * @param isLoggedIn - Authentication state.
 */
export async function displayMyOrders(
  container: HTMLElement | null,
  isLoggedIn?: boolean
): Promise<void> {
  if (!container || !(container instanceof HTMLElement)) {
    console.error("displayMyOrders: Missing or invalid DOM container element.");
    return;
  }

  container.replaceChildren();

  if (!isLoggedIn) {
    container.append(
      createElement("p", { class: "auth-warning" }, ["You must be logged in to view your orders."])
    );
    return;
  }

  // Define initial state typed explicitly
  const state: OrderPageState = {
    orders: [],
    filters: {
      status: "",
      date: "",
    },
    currentPage: 1,
    expandedOrders: new Set<string>(),
    loading: true,
  };

  const render = () => {
    container.replaceChildren(buildOrdersPage(state, render));
  };

  // Show a clear loading state first instead of empty orders summary
  container.replaceChildren(
    createElement("section", { class: "user-orders-page" }, [
      createElement("h2", {}, ["My Orders"]),
      createElement("p", { class: "loading-msg" }, ["Loading your orders..."]),
    ])
  );

  try {
    const res = await getMyOrders();

    // Safely extract orders array from response variations
    const ordersData = Array.isArray(res) ? res : res?.orders;
    if (!Array.isArray(ordersData)) {
      throw new Error("Invalid format received from orders API.");
    }

    state.orders = normalizeOrders(ordersData);
    state.loading = false;

    // First interactive render after successful data load
    render();
  } catch (err) {
    console.error("Failed to fetch user orders:", err);
    state.loading = false;

    container.replaceChildren(
      createElement("section", { class: "user-orders-page" }, [
        createElement("h2", {}, ["My Orders"]),
        createElement("p", { class: "error-msg" }, ["Failed to load orders. Please try again later."]),
      ])
    );
  }
}

export default displayMyOrders;