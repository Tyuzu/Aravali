// generate-manifest.js
import sharp from "sharp";
import { writeFileSync, mkdirSync, existsSync } from "fs";
import path from "path";

type IconEntry = {
  src: string;
  sizes: string;
  type: string;
};

const domain = process.env.VITE_DOMAIN || "https://indium.netlify.app";
const baseIconCandidates = [
  "assets/icon-192.png",
  "assets/logo.png",
  "public/assets/icon-192.png",
  "public/assets/logo.png",
];
const outputDir = path.resolve("public", "assets");
const manifestPath = path.join(path.resolve("public"), "manifest.json");
const sizes = [72, 96, 128, 192, 512];

if (!existsSync(outputDir)) mkdirSync(outputDir, { recursive: true });
const manifestDir = path.dirname(manifestPath);
if (!existsSync(manifestDir)) mkdirSync(manifestDir, { recursive: true });

const baseIcon = baseIconCandidates.find((candidate) => existsSync(candidate));
let icons: IconEntry[] = [];
if (!baseIcon) {
  console.warn("⚠️ no app icon found; manifest will not include generated icons.");
} else {
  for (const size of sizes) {
    const output = path.join(outputDir, `icon-${size}.png`);
    await sharp(baseIcon)
      .resize(size, size, { fit: "cover", position: "centre" })
      .png({ compressionLevel: 9, quality: 90 })
      .toFile(output);
    console.log(`✅ generated icon-${size}.png`);
  }

  icons = sizes.map((size) => ({
    src: `/assets/icon-${size}.png`,
    sizes: `${size}x${size}`,
    type: "image/png",
  }));
}

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
  related_applications: [{ platform: "webapp", url: `${domain}/manifest.json` }],
  prefer_related_applications: false,
  url_handlers: [{ origin: new URL(domain).origin }],
};

writeFileSync(manifestPath, JSON.stringify(manifest, null, 2));
console.log("✅ manifest generated at:", manifestPath);
