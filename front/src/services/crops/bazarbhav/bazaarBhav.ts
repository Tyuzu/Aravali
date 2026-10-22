// displayBazarBhav.ts
import { createElement } from "../../../components/createElement.js";

export function displayBazarBhav(contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  const wrapper = createElement("div", { class: "bazar-bhav-container" });

  const header = createElement("div", { class: "catalog-header" });
  header.appendChild(createElement("h2", {}, ["APMC Bazar Bhav (Market Rates)"]));
  header.appendChild(createElement("p", { class: "text-muted" }, ["Daily wholesale commodity prices from major agricultural markets."]));
  wrapper.appendChild(header);

  // Market Rates Table
  const tableWrapper = createElement("div", { class: "table-responsive" });
  const table = createElement("table", { class: "bazar-table" });
  
  table.innerHTML = `
    <thead>
      <tr>
        <th>Commodity</th>
        <th>Market / APMC</th>
        <th>Min Price (₹/Quintal)</th>
        <th>Max Price (₹/Quintal)</th>
        <th>Modal Price (₹/Quintal)</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><strong>Wheat (Lokwan)</strong></td>
        <td>Delhi APMC</td>
        <td>₹2,450</td>
        <td>₹2,680</td>
        <td>₹2,560</td>
      </tr>
      <tr>
        <td><strong>Soybean</strong></td>
        <td>Indore APMC</td>
        <td>₹4,300</td>
        <td>₹4,750</td>
        <td>₹4,500</td>
      </tr>
      <tr>
        <td><strong>Basmati Rice</strong></td>
        <td>Karnal APMC</td>
        <td>₹3,500</td>
        <td>₹4,100</td>
        <td>₹3,800</td>
      </tr>
      <tr>
        <td><strong>Cotton</strong></td>
        <td>Nagpur APMC</td>
        <td>₹6,800</td>
        <td>₹7,200</td>
        <td>₹7,000</td>
      </tr>
    </tbody>
  `;
  
  tableWrapper.appendChild(table);
  wrapper.appendChild(tableWrapper);
  contentContainer.appendChild(wrapper);
}