import { createElement } from "../../components/createElement.js";
import { secnav, type NavItemConfig } from "../../components/layout/secNav.js";

export function displayMediaEditor(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): void {
  // Main container
  const mediaEditorContainer = createElement(
    "div",
    { class: "media-editor-container" },
    []
  );

  // Header Section (Wrapped for clean layout)
  const headerWrapper = createElement("div", { class: "media-editor-header-wrapper" }, []);
  
  const mediaEditorHeader = createElement(
    "h2",
    { class: "media-editor-header" },
    ["Online Media Editing Suite"]
  );

  const mediaEditorDescription = createElement(
    "p",
    { class: "media-editor-description" },
    ["Edit, manage, scan, and convert your media files instantly in your browser with zero downloads."]
  );

  headerWrapper.appendChild(mediaEditorHeader);
  headerWrapper.appendChild(mediaEditorDescription);

  // Navigation items updated to match new SEO-friendly route paths
  const navItems: NavItemConfig[] = [
    {
      label: "Image Editor",
      href: "/tools/online-image-editor",
      active: window.location.pathname === "/tools/online-image-editor"
    },
    {
      label: "Video Editor",
      href: "/tools/online-video-editor",
      active: window.location.pathname === "/tools/online-video-editor"
    },
    {
      label: "Document Scanner",
      href: "/tools/online-document-scanner",
      active: window.location.pathname === "/tools/online-document-scanner"
    },
    {
      label: "PDF Editor",
      href: "/tools/online-pdf-editor",
      active: window.location.pathname === "/tools/online-pdf-editor"
    },
    {
      label: "Media Converter",
      href: "/tools/online-media-converter",
      active: window.location.pathname === "/tools/online-media-converter"
    }
  ];

  // Secondary navigation
  const secondaryNav = secnav(navItems);

  // Main content section
  const editorContent = createElement(
    "div",
    { class: "media-editor-content" },
    []
  );

  // Welcome section card
  const welcomeSection = createElement(
    "div",
    { class: "media-editor-welcome" },
    [
      createElement("h3", { class: "welcome-title" }, ["Welcome to Your Media Workspace"]),
      createElement("p", { class: "welcome-text" }, [
        "Choose a tool from the navigation above to start editing photos, clipping videos, or scanning documents securely."
      ])
    ]
  );

  editorContent.appendChild(welcomeSection);

  // Login notice card (conditionally rendered)
  if (!isLoggedIn) {
    const loginNotice = createElement(
      "div",
      { class: "media-editor-login-notice" },
      [
        createElement("div", { class: "notice-icon" }, ["ℹ️"]),
        createElement("div", { class: "notice-content" }, [
          createElement("strong", {}, ["Guest Mode Active"]),
          createElement("p", {}, [
            "Some advanced media saving and cloud history features require you to log in."
          ])
        ])
      ]
    );

    editorContent.appendChild(loginNotice);
  }

  // Assemble layout cleanly
  mediaEditorContainer.appendChild(headerWrapper);

  if (secondaryNav) {
    mediaEditorContainer.appendChild(secondaryNav);
  }

  mediaEditorContainer.appendChild(editorContent);

  // Render inside provided container
  contentContainer.replaceChildren(mediaEditorContainer);
}