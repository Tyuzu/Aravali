export function displayBazarBhav(contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const heading = document.createElement("h2");
  heading.textContent = "Bazaar Bhav";
  section.appendChild(heading);

  const table = document.createElement("table");
  table.className = "data-table";

  const body = document.createElement("tbody");
  [
    ["Rice", "₹2,440 / quintal", "Stable"],
    ["Wheat", "₹2,150 / quintal", "Rising"],
    ["Cotton", "₹6,500 / quintal", "Steady"],
    ["Maize", "₹2,200 / quintal", "Mixed"]
  ].forEach(([item, price, trend]) => {
    const row = document.createElement("tr");
    const name = document.createElement("td");
    name.textContent = String(item);
    const priceCell = document.createElement("td");
    priceCell.textContent = String(price);
    const trendCell = document.createElement("td");
    trendCell.textContent = String(trend);
    row.append(name, priceCell, trendCell);
    body.appendChild(row);
  });

  const head = document.createElement("thead");
  const headRow = document.createElement("tr");
  ["Crop", "Price", "Trend"].forEach((label) => {
    const th = document.createElement("th");
    th.textContent = label;
    headRow.appendChild(th);
  });
  head.appendChild(headRow);
  table.append(head, body);
  section.appendChild(table);
  contentContainer.appendChild(section);
}