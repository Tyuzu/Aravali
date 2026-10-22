export function displayProduct(container: HTMLElement, isLoggedIn: boolean, productid: string | number, productType: string): void {
  container.replaceChildren();

  const section = document.createElement("section");
  section.className = "content-page";

  const badge = document.createElement("span");
  badge.className = "page-badge";
  badge.textContent = isLoggedIn ? "Member access" : "Public preview";

  const type = (productType || "product").toString().replace(/[-_]+/g, " ");
  const title = document.createElement("h2");
  title.textContent = `${type.charAt(0).toUpperCase()}${type.slice(1)} ${productid ? `#${productid}` : ""}`.trim();

  const lead = document.createElement("p");
  lead.textContent = "This page is the product detail view, ready to show pricing information, specifications, and marketplace actions.";

  const cards = document.createElement("div");
  cards.className = "info-grid";

  [
    ["Category", type],
    ["Status", "In stock"],
    ["Price", "₹1,499"],
    ["Availability", "2-3 days"]
  ].forEach(([label, value]) => {
    const card = document.createElement("article");
    card.className = "info-card";
    const small = document.createElement("small");
    small.textContent = label;
    const strong = document.createElement("strong");
    strong.textContent = String(value);
    card.append(small, strong);
    cards.appendChild(card);
  });

  section.append(badge, title, lead, cards);
  container.appendChild(section);
}