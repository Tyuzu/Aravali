import "../../../css/layout/sticky.css";
import { createElement } from "../createElement.js";
import { notifSVG, cartSVG, chatCircleSVG, menuSVG, aSVG } from "../svgs/featherSVGs";
import { navigate } from "../../routes/navigate.js";
import { getState, subscribe } from "../../state/state.js";
import { openNotificationsModal } from "../../services/notifications/notifModal.js";
import { goHome, toggleSidebar } from "./sidebar.js";
import { createIconButton } from "../../utils/svgIconButton.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */

export type ImgLinkOption = Node | (() => Node) | null;

export interface StickyExtraOptions {
    imglink?: ImgLinkOption;
}

type UnsubscribeFn = () => void;

/* =========================================================
   BADGE HELPER
========================================================= */

function createBadge(count: number): HTMLElement {
    const displayCount: string = count > 99 ? "99+" : String(count);

    return createElement(
        "span",
        {
            class: "nav-badge",
            "aria-label": `${count} unread`
        },
        [displayCount]
    );
}

/* =========================================================
   NAV UPDATE LOGIC
========================================================= */

function updateNav(container: HTMLElement, extraOptions: StickyExtraOptions = {}): void {
    // Aligned with system-wide state auth resolution
    const isLoggedIn: boolean = Boolean(getState("isLoggedIn"));
    const unreadMessages: number = Number(getState("unreadMessages") || 0);
    const unreadNotifications: number = Number(getState("unreadNotifications") || 0);

    const imglink: ImgLinkOption = extraOptions?.imglink || null;
    
    // State key snapshot to prevent redundant DOM re-renders
    const nextStateKey = `${isLoggedIn}-${unreadMessages}-${unreadNotifications}-${Boolean(imglink)}`;

    if (container.dataset.stateKey === nextStateKey) {
        return;
    }

    container.dataset.stateKey = nextStateKey;
    const fragment: DocumentFragment = document.createDocumentFragment();

    // 0. Home Button
    fragment.appendChild(
        createIconButton({
            classSuffix: "menu",
            svgMarkup: aSVG,
            onClick: goHome,
            label: "Go Home"
        })
    );

    // 1. Sidebar Toggle Button
    fragment.appendChild(
        createIconButton({
            classSuffix: "menu",
            svgMarkup: menuSVG,
            onClick: toggleSidebar,
            label: "Open menu"
        })
    );

    // // 2. Profile / Custom Image Link
    // if (imglink) {
    //     if (typeof imglink === "function") {
    //         fragment.appendChild(imglink());
    //     } else if (imglink instanceof Node) {
    //         fragment.appendChild(imglink);
    //     }
    // }

    // 3. Authenticated Navigation Action Buttons
    if (isLoggedIn) {
        // Chat / Messages Button
        const chatBtn: HTMLElement = createIconButton({
            classSuffix: "stickychat",
            svgMarkup: chatCircleSVG,
            onClick: () => navigate("/newchats"),
            label: "Chats"
        });

        if (unreadMessages > 0) {
            chatBtn.appendChild(createBadge(unreadMessages));
        }
        fragment.appendChild(chatBtn);

        // Shopping Cart Button
        fragment.appendChild(
            createIconButton({
                classSuffix: "cart",
                svgMarkup: cartSVG,
                onClick: () => navigate("/cart"),
                label: "Shopping cart"
            })
        );

        // Notifications Button
        const notifBtn: HTMLElement = createIconButton({
            classSuffix: "notif",
            svgMarkup: notifSVG,
            onClick: openNotificationsModal,
            label: "Notifications"
        });

        if (unreadNotifications > 0) {
            notifBtn.appendChild(createBadge(unreadNotifications));
        }
        fragment.appendChild(notifBtn);
    }

    container.replaceChildren(fragment);
}

/* =========================================================
   STICKY COMPONENT
========================================================= */

export function Sticky(extraOptions: StickyExtraOptions = {}): HTMLDivElement {
    const container = createElement("div", {
        class: "plypzstp"
    }) as HTMLDivElement;

    // Initial render
    updateNav(container, extraOptions);

    let renderAnimationFrame: number | null = null;

    const scheduleUpdate = (): void => {
        if (renderAnimationFrame !== null) {
            cancelAnimationFrame(renderAnimationFrame);
        }
        renderAnimationFrame = requestAnimationFrame(() => {
            updateNav(container, extraOptions);
        });
    };

    // Subscriptions
    const unsubscribers: UnsubscribeFn[] = [
        subscribe("isLoggedIn", scheduleUpdate),
        subscribe("user", scheduleUpdate),
        subscribe("unreadMessages", scheduleUpdate),
        subscribe("unreadNotifications", scheduleUpdate)
    ];

    // Safe Observer pattern: Wait until microtask queue runs so parent can attach container to DOM
    queueMicrotask(() => {
        const observer = new MutationObserver(() => {
            if (!document.body.contains(container)) {
                if (renderAnimationFrame !== null) {
                    cancelAnimationFrame(renderAnimationFrame);
                }
                unsubscribers.forEach((unsub) => unsub?.());
                observer.disconnect();
            }
        });

        if (document.body) {
            observer.observe(document.body, {
                childList: true,
                subtree: true
            });
        }
    });

    return container;
}

export { Sticky as sticky };