// src/ui/cart/printInvoice.ts
import Notify from "../../components/ui/Notify.js";
import { getState } from "../../state/state.js";

// --- INTERFACES & TYPES ---
export interface InvoiceItem {
  price?: number;
  quantity?: number;
  unit?: string;
  itemUnit?: string;
  measurementUnit?: string;
  unitName?: string;
  uom?: string;
  itemName?: string;
  name?: string;
  [key: string]: unknown;
}

export interface InvoiceOrder {
  orderid?: string | number;
  orderId?: string | number;
  OrderID?: string | number;
  id?: string | number;
  _id?: string | number;
  createdAt?: string | number;
  created_at?: string | number;
  status?: string;
  paymentStatus?: string;
  payment_status?: string;
  customerName?: string;
  customer_name?: string;
  name?: string;
  fullName?: string;
  phone?: string;
  customerPhone?: string;
  customer_phone?: string;
  phone_number?: string;
  phoneNumber?: string;
  address?: string;
  deliveryAddress?: string;
  shippingAddress?: string;
  subtotal?: number;
  discount?: number;
  tax?: number;
  delivery?: number;
  total?: number;
  user?: Record<string, unknown>;
  [key: string]: unknown;
}

/* ────────────────────── Core Functional Helpers ────────────────────── */

function toRupees(amount?: number): number {
  return Number(amount || 0) / 100;
}

function formatINR(amount: number): string {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    minimumFractionDigits: 2
  }).format(amount);
}

/**
 * Defensive XSS Sanitizer to safeguard HTML string templates
 */
function escapeHTML(str: unknown): string {
  if (str === undefined || str === null || str === "") {
    return "";
  }

  if (typeof str !== "string" && typeof str !== "number" && typeof str !== "boolean") {
    return "";
  }

  const text = String(str);

  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

function firstNonEmpty(...values: unknown[]): string {
  for (const value of values) {
    if (value !== undefined && value !== null) {
      if (typeof value === "string" || typeof value === "number" || typeof value === "boolean") {
        const text = String(value);
        if (text.trim() !== "") {
          return text;
        }
      }
    }
  }
  return "-";
}

function readNestedString(record: Record<string, unknown> | undefined, keys: string[]): string {
  if (!record) {
    return "";
  }

  for (const key of keys) {
    const value = record[key];
    if (typeof value === "string" && value.trim() !== "") {
      return value;
    }
  }

  return "";
}

/* ────────────────────── Main Print Processing ────────────────────── */

export function printInvoice(order: InvoiceOrder, items: InvoiceItem[] | Record<string, InvoiceItem[]> = []): void {
  if (!order || typeof order !== "object") {
    Notify("Invalid order data payload", { type: "error", duration: 3000 });
    return;
  }

  const invoiceWindow = window.open("", "_blank");

  if (!invoiceWindow) {
    Notify("Pop-up blocked! Please allow windows to view invoice.", {
      type: "warning",
      duration: 5000
    });
    return;
  }

  const flattenedItems = Array.isArray(items) ? items : Object.values(items || {}).flat();

  // FIXED: Corrected line-item total calculation bug and escaped properties securely
  const itemRows = flattenedItems
    .map(item => {
      const price = toRupees(item.price || 0);
      const quantity = Number(item.quantity) || 0;
      const unit = firstNonEmpty(
        item.unit,
        item.itemUnit,
        item.measurementUnit,
        item.unitName,
        item.uom,
        ""
      );
      const quantityLabel = unit ? `${quantity} ${unit}` : String(quantity);
      const total = price * quantity;

      return `
        <tr>
          <td>${escapeHTML(item.itemName || item.name || "Item")}</td>
          <td>${escapeHTML(quantityLabel)}</td>
          <td>${escapeHTML(formatINR(price))}</td>
          <td>${escapeHTML(formatINR(total))}</td>
        </tr>
      `;
    })
    .join("");

  const rawDate = order.createdAt || order.created_at || Date.now();
  const invoiceDate = new Date(rawDate).toLocaleString();
  const invoiceOrderId = firstNonEmpty(
    order.orderid,
    order.orderId,
    order.OrderID,
    order.id,
    order._id,
    "Draft"
  );
  const currentProfile = (getState("userProfile") || {}) as Record<string, unknown>;
  const userProfile = (order.user && typeof order.user === "object" ? order.user : {}) as Record<string, unknown>;
  const fallbackProfile = {
    ...currentProfile,
    ...userProfile,
  } as Record<string, unknown>;

  const customerName = firstNonEmpty(
    order.customerName,
    order.customer_name,
    order.name,
    order.fullName,
    readNestedString(fallbackProfile, ["name", "fullName"]),
    readNestedString(userProfile, ["name", "fullName"]),
    "-"
  );
  const customerPhone = firstNonEmpty(
    order.phone,
    order.customerPhone,
    order.customer_phone,
    order.phone_number,
    order.phoneNumber,
    readNestedString(fallbackProfile, ["phone", "phone_number", "mobile"]),
    readNestedString(userProfile, ["phone", "phone_number", "mobile"]),
    "-"
  );
  const customerAddress = firstNonEmpty(
    order.address,
    order.deliveryAddress,
    order.shippingAddress,
    readNestedString(fallbackProfile, ["address"]),
    readNestedString(userProfile, ["address"]),
    "-"
  );

  invoiceWindow.document.write(`
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Invoice - ${escapeHTML(invoiceOrderId)}</title>
        <style>
          * { box-sizing: border-box; }
          body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Arial, sans-serif;
            padding: 32px;
            color: #222;
            background: #fff;
            line-height: 1.5;
          }
          .header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 24px;
            border-bottom: 2px solid #ddd;
            padding-bottom: 16px;
          }
          .title { margin: 0; font-size: 28px; color: #111; }
          .meta p { margin: 4px 0; }
          .section { margin-top: 24px; }
          .section h3 { margin-bottom: 8px; border-bottom: 1px solid #eee; padding-bottom: 4px; }
          table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 16px;
          }
          th { background: #f8f9fa; font-weight: 600; }
          th, td { border: 1px solid #dee2e6; padding: 12px; text-align: left; }
          td:nth-child(2), td:nth-child(3), td:nth-child(4) { text-align: right; }
          .totals { margin-top: 24px; margin-left: auto; width: 320px; }
          .totals-row { display: flex; justify-content: space-between; margin: 8px 0; }
          .grand-total {
            font-size: 20px;
            font-weight: bold;
            border-top: 2px solid #222;
            padding-top: 10px;
            margin-top: 10px;
          }
          .footer {
            margin-top: 60px;
            text-align: center;
            color: #6c757d;
            font-size: 13px;
          }
          @media print {
            body { padding: 0; }
            .no-print { display: none; }
          }
        </style>
      </head>
      <body>
        <div class="header">
          <div>
            <h1 class="title">Invoice</h1>
            <p style="margin-top: 8px;">Order ID: <strong>${escapeHTML(invoiceOrderId)}</strong></p>
          </div>
          <div class="meta">
            <p><strong>Date:</strong> ${escapeHTML(invoiceDate)}</p>
            <p><strong>Status:</strong> ${escapeHTML(order.status || "-")}</p>
            <p><strong>Payment:</strong> ${escapeHTML(order.paymentStatus || order.payment_status || "-")}</p>
          </div>
        </div>

        <div class="section">
          <h3>Customer Details</h3>
          <p><strong>Name:</strong> ${escapeHTML(customerName)}</p>
          <p><strong>Phone:</strong> ${escapeHTML(customerPhone)}</p>
          <p><strong>Address:</strong> ${escapeHTML(customerAddress)}</p>
        </div>

        <div class="section">
          <h3>Items Summary</h3>
          <table>
            <thead>
              <tr>
                <th>Item</th>
                <th>Qty / Unit</th>
                <th>Price</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              ${itemRows}
            </tbody>
          </table>
        </div>

        <div class="totals">
          <div class="totals-row">
            <span>Subtotal</span>
            <span>${escapeHTML(formatINR(toRupees(order.subtotal || 0)))}</span>
          </div>
          <div class="totals-row">
            <span>Discount</span>
            <span>${escapeHTML(formatINR(toRupees(order.discount || 0)))}</span>
          </div>
          <div class="totals-row">
            <span>Tax</span>
            <span>${escapeHTML(formatINR(toRupees(order.tax || 0)))}</span>
          </div>
          <div class="totals-row">
            <span>Delivery</span>
            <span>${escapeHTML(formatINR(toRupees(order.delivery || 0)))}</span>
          </div>
          <div class="totals-row grand-total">
            <span>Total Amount</span>
            <span>${escapeHTML(formatINR(toRupees(order.total || 0)))}</span>
          </div>
        </div>

        <div class="footer">
          Thank you for choosing us for your purchase.
        </div>

        <script>
          document.addEventListener("DOMContentLoaded", () => {
            setTimeout(() => {
              window.print();
              window.onafterprint = () => { window.close(); };
            }, 500);
          });
        </script>
      </body>
    </html>
  `);

  invoiceWindow.document.close();
}

export default printInvoice;