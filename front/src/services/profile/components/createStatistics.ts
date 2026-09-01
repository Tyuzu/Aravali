import type { StatItem, UserProfile } from "../profileGenHelpers.js";

export function createStatistics(profile: UserProfile): HTMLDivElement {
    const statistics = document.createElement("div");
    statistics.className = "statistics";

    const stats: StatItem[] = [
        { label: "Rupees", value: profile.wallet_balance || 0 },
        { label: "Followers", value: profile.followerscount || 0 },
        { label: "Following", value: profile.followscount || 0 },
    ];

    stats.forEach(({ label, value }) => {
        const statItem = document.createElement("p");
        statItem.className = "hflex";

        const strong = document.createElement("strong");
        strong.textContent = String(value);

        const labelSpan = document.createTextNode(` ${label}`);

        statItem.appendChild(strong);
        statItem.appendChild(labelSpan);
        statistics.appendChild(statItem);
    });

    return statistics;
}
