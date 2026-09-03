import fs from "fs";
import path from "path";

const BASE_URL = process.env.VITE_BASE_URL || "https://indium.netlify.app";

const routes = [
  "/", "baitos", "baitos/hire", "grocery", "recipes",
  "places", "itinerary", "events", "music", "artists",
  "social", "posts", "farms", "products", "tools"
];

const today = new Date().toISOString().split("T")[0]; // YYYY-MM-DD

const sitemap =
`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes.map(route => `
  <url>
    <loc>${BASE_URL}${route === "/" ? "/" : "/" + route.replace(/^\//, "")}</loc>
    <lastmod>${today}</lastmod>
    <changefreq>weekly</changefreq>
    <priority>${route === "/" ? "1.0" : "0.7"}</priority>
  </url>`).join("")}
</urlset>`;
// ensure public dir exists and write into it
const outDir = path.resolve("public");
if (!fs.existsSync(outDir)) fs.mkdirSync(outDir, { recursive: true });
const outPath = path.join(outDir, "sitemap.xml");
fs.writeFileSync(outPath, sitemap);
console.log("✅ Sitemap generated at:", outPath);
