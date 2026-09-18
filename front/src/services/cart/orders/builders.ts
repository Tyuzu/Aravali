import {
  getFilteredOrders,
  toggleExpanded,
  getOrderProducts,
  getOrderSummaryMeta,
  formatDate,
  formatINR,
  capitalize,
  downloadReceipt,
} from "./orderutils.js";
import { createElement } from "../../../components/createElement.js";
import { Button } from "../../../components/base/Button.js";
import { Order, OrderItem, OrderPageState } from "./types.js";

const PAGE_SIZE = 5;
const STATUS_OPTIONS = [
  { value: "", label: "All" },
  { value: "pending", label: "Pending" },
  { value: "accepted", label: "Accepted" },
  { value: "paid", label: "Paid" },
  { value: "delivered", label: "Delivered" },
  { value: "rejected", label: "Rejected" },
  { value: "active", label: "Active" },
  { value: "closed", label: "Closed" },
];

type RerenderCallback = () => void;

/* ───────────────── Helper ───────────────── */
const valOrFallback = (val: string | undefined | null, fallback = "N/A") => val || fallback;

/* ───────────────── Page ───────────────── */

export function buildOrdersPage(state: OrderPageState, rerender: RerenderCallback): HTMLElement {
  const filteredOrders = getFilteredOrders(state.orders, state.filters);
  const totalPages = Math.max(1, Math.ceil(filteredOrders.length / PAGE_SIZE));
  
  if (state.currentPage > totalPages) {
    state.currentPage = totalPages;
  }

  const pagedOrders = filteredOrders.slice(
    (state.currentPage - 1) * PAGE_SIZE,
    state.currentPage * PAGE_SIZE
  );

  const isMobile = window.innerWidth <= 768;

  return createElement("section", { class: "user-orders-page" }, [
    createElement("h2", {}, ["My Orders"]),
    buildUserOrderFilters(state, rerender),
    buildOrdersSummary(filteredOrders.length, state.orders.length, state.currentPage, totalPages),
    isMobile
      ? buildMobileOrdersList(pagedOrders, state, rerender)
      : buildDesktopOrdersTable(pagedOrders, state, rerender),
    buildPaginationControls(state, filteredOrders.length, totalPages, rerender),
  ]);
}

/* ───────────────── Filters ───────────────── */

export function buildUserOrderFilters(state: OrderPageState, rerender: RerenderCallback): HTMLElement {
  const handleFilterChange = (key: keyof typeof state.filters, value: string) => {
    state.filters[key] = value;
    state.currentPage = 1;
    rerender();
  };

  return createElement("div", { class: "filters" }, [
    buildLabeledSelect("Status", STATUS_OPTIONS, state.filters.status, (val) =>
      handleFilterChange("status", val)
    ),
    createElement("label", {}, [
      "Date: ",
      createElement("input", {
        type: "date",
        value: state.filters.date || "",
        onchange: (e: Event) =>
          handleFilterChange("date", (e.target as HTMLInputElement).value),
      }),
    ]),
    Button({
      title: "Reset",
      events: {
        click: () => {
          state.filters.status = "";
          state.filters.date = "";
          state.currentPage = 1;
          rerender();
        },
      },
      class:"buttonx",
    }),
  ]);
}

export function buildOrdersSummary(
  filteredCount: number,
  totalCount: number,
  currentPage: number,
  totalPages: number
): HTMLElement {
  return createElement("p", { class: "orders-summary" }, [
    `Showing ${filteredCount} of ${totalCount} order(s) · Page ${currentPage} of ${totalPages}`,
  ]);
}

/* ───────────────── Desktop Table ───────────────── */

function buildDesktopOrdersTable(
  orders: Order[],
  state: OrderPageState,
  rerender: RerenderCallback
): HTMLElement {
  const headers = ["", "Order ID", "Date", "Type", "Total", "Status", "Payment", "Actions"];

  const tbodyChildren = orders.length
    ? orders.flatMap((order) => buildExpandableOrderRows(order, state, rerender))
    : [
        createElement("tr", {}, [
          createElement("td", { colspan: "8" }, ["No orders found."]),
        ]),
      ];

  return createElement("table", { class: "orders-table" }, [
    createElement("thead", {}, [
      createElement("tr", {}, headers.map((h) => createElement("th", {}, [h]))),
    ]),
    createElement("tbody", {}, tbodyChildren),
  ]);
}

function buildExpandableOrderRows(
  order: Order,
  state: OrderPageState,
  rerender: RerenderCallback
): HTMLElement[] {
  const expanded = state.expandedOrders.has(order.orderId);
  const products = getOrderProducts(order) || [];
  const meta = getOrderSummaryMeta(order);

  const addressInfo = valOrFallback(meta.address);
  const farmInfo = valOrFallback(meta.farmId);
  const approvedList =
    Array.isArray(meta.approvedBy) && meta.approvedBy.length ? meta.approvedBy.join(", ") : "N/A";

  const rows: HTMLElement[] = [
    createElement("tr", { class: "order-summary-row" }, [
      createElement("td", {}, [
        Button({
          title: expanded ? "−" : "+",
          classes: "toggle-btn",
          events: {
            click: () => {
              toggleExpanded(state, order.orderId);
              rerender();
            },
          },
        }),
      ]),
      createElement("td", {}, [String(meta.orderId || order.orderId)]),
      createElement("td", {}, [formatDate(order.createdAt)]),
      createElement("td", {}, [capitalize(valOrFallback(meta.orderType))]),
      createElement("td", {}, [formatINR(order.total || 0, true)]),
      createElement("td", {}, [capitalize(valOrFallback(meta.status))]),
      createElement("td", {}, [capitalize(valOrFallback(meta.payment))]),
      createElement("td", {}, [
        Button({
          title: "Receipt",
          events: { click: () => downloadReceipt(order) },
        }),
      ]),
    ]),
  ];

  // Don't inject hidden row nodes if not expanded
  if (expanded) {
    rows.push(
      createElement("tr", { class: "order-detail-row" }, [
        createElement("td", { colspan: "8" }, [
          createElement("div", { class: "order-detail-grid" }, [
            createElement("p", {}, [`Payment: ${capitalize(valOrFallback(meta.payment))}`]),
            createElement("p", {}, [`Address: ${addressInfo}`]),
            createElement("p", {}, [`Farm: ${farmInfo}`]),
            createElement("p", {}, [`Approved By: ${approvedList}`]),
            buildOrderItemsTable(products, farmInfo),
          ]),
        ]),
      ])
    );
  }

  return rows;
}

function buildOrderItemsTable(products: OrderItem[], farmFallback = "N/A"): HTMLElement {
  const rows = products.length
    ? products.map((item) =>
        createElement("tr", {}, [
          createElement("td", {}, [valOrFallback(item.entityName, farmFallback)]),
          createElement("td", {}, [valOrFallback(item.itemName)]),
          createElement("td", {}, [String(item.quantity || 0)]),
          createElement("td", {}, [formatINR(item.price || 0, true)]),
        ])
      )
    : [
        createElement("tr", {}, [
          createElement("td", { colspan: "4" }, ["No items found."]),
        ]),
      ];

  return createElement("table", { class: "order-items-table" }, [
    createElement("thead", {}, [
      createElement("tr", {}, ["Entity", "Item", "Qty", "Item Price"].map((h) => createElement("th", {}, [h]))),
    ]),
    createElement("tbody", {}, rows),
  ]);
}

/* ───────────────── Mobile Cards ───────────────── */

function buildMobileOrdersList(
  orders: Order[],
  state: OrderPageState,
  rerender: RerenderCallback
): HTMLElement {
  return createElement(
    "div",
    { class: "orders-cards" },
    orders.length
      ? orders.map((order) => buildExpandableOrderCard(order, state, rerender))
      : [createElement("p", {}, ["No orders found."])]
  );
}

function buildExpandableOrderCard(
  order: Order,
  state: OrderPageState,
  rerender: RerenderCallback
): HTMLElement {
  const expanded = state.expandedOrders.has(order.orderId);
  const products = getOrderProducts(order) || [];
  const meta = getOrderSummaryMeta(order);

  const cardChildren: HTMLElement[] = [
    createElement("div", { class: "order-card-header" }, [
      createElement("p", {}, [`Order ID: ${meta.orderId || order.orderId}`]),
      Button({
        title: expanded ? "Collapse" : "Expand",
        events: {
          click: () => {
            toggleExpanded(state, order.orderId);
            rerender();
          },
        },
        class: "buttonx"
      }),
    ]),
    createElement("p", {}, [`Date: ${formatDate(order.createdAt)}`]),
    createElement("p", {}, [`Type: ${capitalize(valOrFallback(meta.orderType))}`]),
    createElement("p", {}, [`Status: ${capitalize(valOrFallback(meta.status))}`]),
    createElement("p", {}, [`Payment: ${capitalize(valOrFallback(meta.payment))}`]),
    createElement("p", {}, [`Address: ${valOrFallback(meta.address)}`]),
    createElement("p", {}, [`Total: ${formatINR(order.total || 0, true)}`]),
  ];

  if (expanded) {
    const approvedList =
      Array.isArray(meta.approvedBy) && meta.approvedBy.length ? meta.approvedBy.join(", ") : "N/A";

    cardChildren.push(
      createElement("div", { class: "order-card-items" }, [
        createElement("p", {}, [`Farm ID: ${valOrFallback(meta.farmId)}`]),
        createElement("p", {}, [`Approved By: ${approvedList}`]),
        ...products.map((item) =>
          createElement("div", { class: "order-card-item" }, [
            createElement("p", {}, [`Farm: ${valOrFallback(item.entityName, valOrFallback(meta.farmId))}`]),
            createElement("p", {}, [`Item: ${valOrFallback(item.itemName)}`]),
            createElement("p", {}, [`Qty: ${item.quantity || 0}`]),
            createElement("p", {}, [`Item Price: ${formatINR(item.price || 0, true)}`]),
          ])
        ),
      ])
    );
  }

  cardChildren.push(
    Button({
      title: "Receipt",
      classes: "buttonx btn-receipt",
      events: { click: () => downloadReceipt(order) },
    })
  );

  return createElement("div", { class: "order-card" }, cardChildren);
}

/* ───────────────── Pagination ───────────────── */

function buildPaginationControls(
  state: OrderPageState,
  totalOrders: number,
  totalPages: number,
  rerender: RerenderCallback
): HTMLElement {
  const createNavButton = (title: string, isDisabled: boolean, delta: number) => {
    const btn = Button({
      title,
      events: {
        click: () => {
          if (!isDisabled) {
            state.currentPage += delta;
            rerender();
          }
        },
      },
    });
    if (isDisabled) btn.setAttribute("disabled", "true");
    return btn;
  };

  return createElement("div", { class: "pagination" }, [
    createNavButton("Prev", state.currentPage <= 1, -1),
    createElement("span", {}, [`Page ${state.currentPage} of ${totalPages} · ${totalOrders} order(s)`]),
    createNavButton("Next", state.currentPage >= totalPages, 1),
  ]);
}

/* ───────────────── Utilities ───────────────── */

interface SelectOption {
  value: string;
  label: string;
}

function buildLabeledSelect(
  labelText: string,
  options: SelectOption[],
  currentValue: string,
  onChange: (value: string) => void
): HTMLElement {
  return createElement("label", {}, [
    `${labelText}: `,
    createElement(
      "select",
      {
        onchange: (e: Event) => onChange((e.target as HTMLSelectElement).value),
      },
      options.map((o) =>
        createElement(
          "option",
          String(o.value) === String(currentValue) ? { value: o.value, selected: "selected" } : { value: o.value },
          [o.label]
        )
      )
    ),
  ]);
}