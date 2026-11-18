import { AppRoute } from "./routeTypes.js";

export const mediaRoutes: AppRoute[] = [
  /* =======================================================
     MEDIA EDITORS HUB
  ================================================ ======= */
  {
    path: "/media-tools",
    component: () =>
      import(
        "../../pages/editors/mediaEditors.js"
      ),
    functionName: "MediaEditor",
    meta: {
      title: "Online Media Editing Suite - Free Image, Video & Document Tools",
      description: "Discover our comprehensive suite of online media tools. Edit images, edit videos, and scan documents directly in your browser with zero downloads.",
      keywords: "online media editor, free image editor, video editor online, document scanner web"
    }
  },

  /* =======================================================
     IMAGE EDITOR
  ================================================       * */
  {
    path: "/tools/online-image-editor",
    component: () =>
      import(
        "../../pages/editors/imageEditor.js"
      ),
    functionName: "ImageEditor",
    meta: {
      title: "Free Online Image Editor - Crop, Filter & Enhance Photos",
      description: "Edit your photos instantly with our free online image editor. Crop, resize, apply filters, and enhance pictures right in your browser.",
      keywords: "online photo editor, free image enhancer, crop image online, photo filters"
    }
  },

  /* =======================================================
     VIDEO EDITOR
  ================================================       * */
  {
    path: "/tools/online-video-editor",
    component: () =>
      import(
        "../../pages/editors/videoEditor.js"
      ),
    functionName: "VideoEditor",
    meta: {
      title: "Free Online Video Editor - Trim, Cut & Edit Videos Fast",
      description: "Trim, cut, and edit your videos online with ease. No watermarks, no complicated software installation required.",
      keywords: "online video editor, free video trimmer, cut video online, web video maker"
    }
  },

  /* =======================================================
     SCAN EDITOR / DOCUMENT SCANNER
  ================================================       * */
  {
    path: "/tools/online-document-scanner",
    component: () =>
      import(
        "../../pages/editors/scanEditor.js"
      ),
    functionName: "ScanEditor",
    meta: {
      title: "Online Document Scanner - Scan & Convert Papers to PDF",
      description: "Transform your physical documents into clean digital PDFs instantly using our online document scanner tool.",
      keywords: "online document scanner, scan to pdf, paper scanner web app, OCR scanner"
    }
  },

/* =======================================================
     PDF Editor
  ======================================================= */
  {
    path: "/tools/online-pdf-editor",
    component: () =>
      import(
        "../../pages/editors/pdfEditor.js"
      ),
    functionName: "PDFEditor",
    meta: {
      title: "Online PDF Editor - Edit, Annotate & Modify PDFs Free",
      description: "Edit, annotate, add text, sign, and rearrange pages in your PDF files directly in your browser—fast, free, and without installing software.",
      keywords: "online pdf editor, edit pdf online, free pdf editor, annotate pdf, modify pdf, sign pdf online"
    }
  },

  /* =======================================================
     DYNAMIC ROUTES (SEO Friendly Slug Support)
  ================================================       * */
  {
    path: "/post/:id/:slug",
    component: () =>
      import(
        "../../pages/editors/postViewPage.js"
      ),
    functionName: "PostViewer",
    meta: {
      title: "Media Post Viewer",
      description: "View and share edited media creations and community posts.",
      keywords: "media post, user creation, shared design"
    }
  },
];