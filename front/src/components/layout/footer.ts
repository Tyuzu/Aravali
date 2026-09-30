import "../../../css/layout/footer.css";
import { getLanguage, setLanguage, t } from "../../i18n/i18n.js";
import { navigate } from "../../routes/navigate.js";
import { webSiteName } from "../../config/env.js";
import { createElement } from "../createElement.js";

interface NavPage {
  href: string;
  label: string;
}

/** Handle navigation cleanly */
const handleNavigation = (event: MouseEvent, href: string): void => {
  event.preventDefault();
  if (!href) {
    console.error("🚨 handleNavigation received null href");
    return;
  }
  navigate(href);
};

const Footer = (): HTMLElement => {
  const pages: NavPage[] = [
    { href: "/about", label: t("footer.aboutUs", {}, "About Us") },
    { href: "/contact", label: t("footer.contactUs", {}, "Contact Us") },
    { href: "/faq", label: t("footer.faq", {}, "FAQ") },
    { href: "/terms", label: t("footer.terms", {}, "Terms & Conditions") },
    { href: "/privacy", label: t("footer.privacy", {}, "Privacy Policy") },
    { href: "/refund", label: t("footer.refund", {}, "Refund Policy") },
    { href: "/shipping", label: t("footer.shipping", {}, "Shipping Policy") },
    { href: "/returns", label: t("footer.returns", {}, "Return Policy") },
    { href: "/disclaimer", label: t("footer.disclaimer", {}, "Disclaimer") },
    { href: "/blog", label: t("footer.blog", {}, "Blog") }
  ];

  const navLinks = pages.map(({ href, label }) =>
    createElement(
      "a",
      {
        href,
        class: "footer-link",
        events: {
          click: (e: Event) => handleNavigation(e as MouseEvent, href)
        }
      },
      [label]
    )
  );

  const nav = createElement("nav", { class: "footer-nav", "aria-label": "Footer Navigation" }, navLinks);

  const currentYear = new Date().getFullYear();
  const currentLang = typeof getLanguage === "function" ? getLanguage() : (localStorage.getItem("lang") || "en");

  const langSelect = createElement(
    "select",
    {
      name: "lang-select",
      class: "lang-select",
      "aria-label": t("footer.languageLabel", {}, "Select Page Language"),
      events: {
        change: async (e: Event) => {
          const target = e.target as HTMLSelectElement | null;
          const lang = target?.value;
          if (lang) {
            await setLanguage(lang);
          }
        }
      }
    },
    [
      createElement("option", { value: "en" }, [t("footer.lang.en", {}, "English")]),
      createElement("option", { value: "es" }, [t("footer.lang.es", {}, "Español")]),
      createElement("option", { value: "fr" }, [t("footer.lang.fr", {}, "Français")]),
      createElement("option", { value: "hi" }, [t("footer.lang.hi", {}, "हिन्दी")]),
      createElement("option", { value: "ar" }, [t("footer.lang.ar", {}, "العربية")]),
      createElement("option", { value: "jp" }, [t("footer.lang.jp", {}, "日本語")])
    ]
  ) as HTMLSelectElement;

  langSelect.value = currentLang;

  const copyrightText = t(
    "footer.copyright",
    { year: currentYear, site: webSiteName },
    `© ${currentYear} ${webSiteName}. All rights reserved.`
  );

  const footerBottom = createElement("div", { class: "footer-bottom" }, [
    langSelect,
    createElement("p", { class: "footer-copyright" }, [copyrightText])
  ]);

  return createElement("footer", { class: "footer-container" }, [
    nav,
    footerBottom
  ]);
};

export { Footer };