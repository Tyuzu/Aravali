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
const baseIcon = "assets/logo.png";      // your source logo
const outputDir = path.resolve("public", "assets");
const manifestPath = path.join(path.resolve("public"), "manifest.json");

const sizes = [72, 96, 128, 192, 512];

// ensure output directory exists
if (!existsSync(outputDir)) mkdirSync(outputDir, { recursive: true });

// ensure manifest directory exists
const manifestDir = path.dirname(manifestPath);
if (!existsSync(manifestDir)) mkdirSync(manifestDir, { recursive: true });

// generate icons (skip if base icon missing)
let icons: IconEntry[] = [];
if (!existsSync(baseIcon)) {
  console.warn(`⚠️ base icon not found at ${baseIcon}, skipping icon generation.`);
} else {
  for (const size of sizes) {
    const output = path.join(outputDir, `icon-${size}.png`);
    await sharp(baseIcon)
      .resize(size, size)
      .png({ compressionLevel: 9 })
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
  name: "Scav",
  short_name: "SPA",
  start_url: "/?source=homescreen",
  scope: "/",
  display: "standalone",
  orientation: "portrait",
  theme_color: "#222222",
  background_color: "#ffffff",
  icons,
  shortcuts: [
    {
      name: "Dashboard",
      short_name: "Dashboard",
      url: "/dashboard",
      icons: [{ src: "/assets/icon-96.png", sizes: "96x96" }],
    },
    {
      name: "Profile",
      short_name: "Profile",
      url: "/profile",
      icons: [{ src: "/assets/icon-96.png", sizes: "96x96" }],
    },
  ],
  related_applications: [
    { platform: "webapp", url: `${domain}/manifest.json` },
    // { platform: "play", id: "com.yourcompany.app" },
  ],
  prefer_related_applications: false,
  url_handlers: [
    { origin: new URL(domain).origin },
  ],
};

// write manifest
writeFileSync(manifestPath, JSON.stringify(manifest, null, 2));
console.log("✅ manifest generated at:", manifestPath);
