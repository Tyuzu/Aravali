import { formatDate } from "../profileHelpers.js";
import type { UserProfile } from "../profileGenHelpers.js";

export interface InfoItem {
    label: string;
    value: string | Node;
}

export function createProfileInfo(profile: UserProfile): HTMLDivElement {
    const profileInfo = document.createElement("div");
    profileInfo.className = "profile-info";

    const lastLoginRaw = profile.last_login;
    const lastLoginStr = lastLoginRaw instanceof Date ? lastLoginRaw.toISOString() : (lastLoginRaw as string | undefined);
    const lastLoginValue = formatDate(lastLoginStr) || "Never logged in";

    const infoItems: InfoItem[] = [
        { label: "Last Login", value: lastLoginValue },
        { label: "Verification Status", value: profile.is_verified ? "Verified" : "Not Verified" },
    ];

    infoItems.forEach(({ label, value }) => {
        const infoItem = document.createElement("div");
        infoItem.className = "info-item";

        const strongLabel = document.createElement("strong");
        strongLabel.textContent = `${label}:`;
        infoItem.appendChild(strongLabel);

        if (value instanceof Node) {
            infoItem.appendChild(document.createTextNode(" "));
            infoItem.appendChild(value);
        } else {
            infoItem.appendChild(document.createTextNode(` ${value}`));
        }

        profileInfo.appendChild(infoItem);
    });

    return profileInfo;
}
