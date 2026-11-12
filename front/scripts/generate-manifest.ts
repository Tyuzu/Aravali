// generate-manifest.mjs
import { writeFileSync, mkdirSync, existsSync } from "fs";
import path from "path";
import sharp from "sharp";

const domain = process.env.VITE_DOMAIN || "https://indium.netlify.app";
const outputDir = path.resolve("public");
const assetDir = path.join(outputDir, "assets");
const manifestPath = path.join(outputDir, "manifest.json");

if (!existsSync(outputDir)) mkdirSync(outputDir, { recursive: true });
if (!existsSync(assetDir)) mkdirSync(assetDir, { recursive: true });

const sizeSet = [72, 96, 128, 192, 512];
const sourceCandidates = [
  path.resolve("assets", "icon-192.png"),
  path.resolve("assets", "logo.png"),
  path.resolve("public", "assets", "icon-192.png"),
  path.resolve("public", "assets", "logo.png"),
];

const sourceIcon = sourceCandidates.find((candidate) => existsSync(candidate));

if (sourceIcon) {
  for (const size of sizeSet) {
    const outputPath = path.join(assetDir, `icon-${size}.png`);
    await sharp(sourceIcon)
      .resize(size, size, { fit: "cover", position: "centre" })
      .png({ compressionLevel: 9, quality: 90 })
      .toFile(outputPath);
  }
  console.log("✅ generated PWA icons from:", sourceIcon);
} else {
  console.warn("⚠️ no source icon found for PWA generation; manifest will be written without generated icons");
}

const icons = sizeSet
  .filter((size) => existsSync(path.join(assetDir, `icon-${size}.png`)))
  .map((size) => ({
    src: `/assets/icon-${size}.png`,
    sizes: `${size}x${size}`,
    type: "image/png",
  }));

const manifest = {
  id: "/",
  name: "Scav",
  short_name: "Scav",
  description: "Scav is a mobile-friendly community app for discovery, buying, selling, and managing everyday experiences.",
  start_url: "/?source=homescreen",
  scope: "/",
  display: "standalone",
  display_override: ["standalone", "browser"],
  orientation: "portrait",
  lang: "en",
  theme_color: "#222222",
  background_color: "#ffffff",
  categories: ["shopping", "productivity", "lifestyle"],
  icons,
  shortcuts: [
    {
      name: "Dashboard",
      short_name: "Dashboard",
      url: "/dashboard",
      icons: [{ src: "/assets/icon-96.png", sizes: "96x96", type: "image/png" }],
    },
    {
      name: "Profile",
      short_name: "Profile",
      url: "/profile",
      icons: [{ src: "/assets/icon-96.png", sizes: "96x96", type: "image/png" }],
    },
  ],
  related_applications: [
    { platform: "webapp", url: `${domain}/manifest.json` },
  ],
  prefer_related_applications: false,
  url_handlers: [{ origin: new URL(domain).origin }],
};

writeFileSync(manifestPath, JSON.stringify(manifest, null, 2));
console.log("✅ manifest generated at:", manifestPath);
