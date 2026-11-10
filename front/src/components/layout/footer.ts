import "../../../css/layout/footer.css";
import { setLanguage, t } from "../../i18n/i18n.js";
import { navigate } from "../../routes/navigate.js";
import { webSiteName } from "../../config/envConfig.js";
import { createElement } from "../createElement.js";

interface NavPage {
  href: string;
  label: string;
}

const footerLocaleLabels: Record<string, Record<string, string>> = {
  en: {
    aboutUs: "About Us",
    contactUs: "Contact Us",
    faq: "FAQ",
    terms: "Terms & Conditions",
    privacy: "Privacy Policy",
    refund: "Refund Policy",
    shipping: "Shipping Policy",
    returns: "Return Policy",
    disclaimer: "Disclaimer",
    blog: "Blog",
    languageLabel: "Select Page Language",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  },
  es: {
    aboutUs: "Sobre nosotros",
    contactUs: "Contáctanos",
    faq: "Preguntas frecuentes",
    terms: "Términos y condiciones",
    privacy: "Política de privacidad",
    refund: "Política de reembolso",
    shipping: "Política de envío",
    returns: "Política de devoluciones",
    disclaimer: "Descargo de responsabilidad",
    blog: "Blog",
    languageLabel: "Selecciona el idioma de la página",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  },
  fr: {
    aboutUs: "À propos",
    contactUs: "Contactez-nous",
    faq: "FAQ",
    terms: "Conditions générales",
    privacy: "Politique de confidentialité",
    refund: "Politique de remboursement",
    shipping: "Politique d'expédition",
    returns: "Politique de retour",
    disclaimer: "Avertissement",
    blog: "Blog",
    languageLabel: "Choisir la langue de la page",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  },
  hi: {
    aboutUs: "हमारे बारे में",
    contactUs: "संपर्क करें",
    faq: "सामान्य प्रश्न",
    terms: "नियम और शर्तें",
    privacy: "गोपनीयता नीति",
    refund: "रिफंड नीति",
    shipping: "शिपिंग नीति",
    returns: "रिटर्न नीति",
    disclaimer: "अस्वीकरण",
    blog: "ब्लॉग",
    languageLabel: "पृष्ठ की भाषा चुनें",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  },
  ar: {
    aboutUs: "معلومات عنا",
    contactUs: "تواصل معنا",
    faq: "الأسئلة الشائعة",
    terms: "الشروط والأحكام",
    privacy: "سياسة الخصوصية",
    refund: "سياسة الاسترداد",
    shipping: "سياسة الشحن",
    returns: "سياسة الإرجاع",
    disclaimer: "إخلاء المسؤولية",
    blog: "المدونة",
    languageLabel: "اختر لغة الصفحة",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  },
  jp: {
    aboutUs: "会社概要",
    contactUs: "お問い合わせ",
    faq: "よくある質問",
    terms: "利用規約",
    privacy: "プライバシーポリシー",
    refund: "返金ポリシー",
    shipping: "配送ポリシー",
    returns: "返品ポリシー",
    disclaimer: "免責事項",
    blog: "ブログ",
    languageLabel: "ページの言語を選択",
    languageEn: "English",
    languageEs: "Español",
    languageFr: "Français",
    languageHi: "हिन्दी",
    languageAr: "العربية",
    languageJp: "日本語"
  }
};

function getFooterLocale(lang: string): Record<string, string> {
  const normalizedLang = (lang || "en").toLowerCase();
  const aliasMap: Record<string, string> = {
    ja: "jp",
    jp: "jp"
  };
  const safeLang = aliasMap[normalizedLang] || normalizedLang;
  return footerLocaleLabels[safeLang] || footerLocaleLabels.en;
}

const handleNavigation = (event: MouseEvent, href: string): void => {
  event.preventDefault();
  if (!href) {
    console.error("handleNavigation received null href");
    return;
  }
  navigate(href);
};

const Footer = (): HTMLElement => {
  const savedLang = (localStorage.getItem("lang") || "en").toLowerCase();
  const locale = getFooterLocale(savedLang);
  const supportedLangs = ["en", "es", "fr", "hi", "ar", "jp"];
  const normalizedSavedLang = savedLang === "ja" ? "jp" : savedLang;
  const activeLang = supportedLangs.includes(normalizedSavedLang) ? normalizedSavedLang : "en";

  const pages: NavPage[] = [
    { href: "/about", label: locale.aboutUs || t("footer.aboutUs", {}, "About Us") },
    { href: "/contact", label: locale.contactUs || t("footer.contactUs", {}, "Contact Us") },
    { href: "/faq", label: locale.faq || t("footer.faq", {}, "FAQ") },
    { href: "/terms", label: locale.terms || t("footer.terms", {}, "Terms & Conditions") },
    { href: "/privacy", label: locale.privacy || t("footer.privacy", {}, "Privacy Policy") },
    { href: "/refund", label: locale.refund || t("footer.refund", {}, "Refund Policy") },
    { href: "/shipping", label: locale.shipping || t("footer.shipping", {}, "Shipping Policy") },
    { href: "/returns", label: locale.returns || t("footer.returns", {}, "Return Policy") },
    { href: "/disclaimer", label: locale.disclaimer || t("footer.disclaimer", {}, "Disclaimer") },
    { href: "/blog", label: locale.blog || t("footer.blog", {}, "Blog") }
  ];

  const navLinks = pages.map(({ href, label }) => {
    return createElement(
      "a",
      {
        href,
        class: "footer-link",
        events: {
          click: ((e: MouseEvent) => handleNavigation(e, href)) as EventListener
        }
      },
      [label]
    );
  });

  const nav = createElement("nav", { class: "footer-nav" }, navLinks);

  const langSelect = createElement(
    "select",
    {
      name: "lang-select",
      class: "lang-select",
      "aria-label": locale.languageLabel || t("footer.languageLabel", {}, "Select Page Language"),
      events: {
        change: (async (e: Event) => {
          const target = e.target as HTMLSelectElement | null;
          const lang = target?.value;
          if (lang) {
            localStorage.setItem("lang", lang);
            await setLanguage(lang);
          }
        }) as EventListener
      }
    },
    [
      createElement("option", { value: "en" }, [locale.languageEn || t("footer.lang.en", {}, "English")]),
      createElement("option", { value: "es" }, [locale.languageEs || t("footer.lang.es", {}, "Español")]),
      createElement("option", { value: "fr" }, [locale.languageFr || t("footer.lang.fr", {}, "Français")]),
      createElement("option", { value: "hi" }, [locale.languageHi || t("footer.lang.hi", {}, "हिन्दी")]),
      createElement("option", { value: "ar" }, [locale.languageAr || t("footer.lang.ar", {}, "العربية")]),
      createElement("option", { value: "jp" }, [locale.languageJp || t("footer.lang.jp", {}, "日本語")])
    ]
  ) as HTMLSelectElement;

  langSelect.value = activeLang;

  const footerBottom = createElement("div", { class: "footer-bottom" }, [
    langSelect,
    createElement("p", {}, [
      t("footer.copyright", { year: new Date().getFullYear(), site: webSiteName }, `© ${new Date().getFullYear()} ${webSiteName}. All rights reserved.`)
    ])
  ]);

  return createElement("div", { class: "footer-container" }, [
    nav,
    footerBottom
  ]);
};

export { Footer };