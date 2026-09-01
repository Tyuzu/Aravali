import { getState } from "../../../state/state.js";
import Button, { ButtonOptions } from "../../../components/base/Button.js";
import { toggleAction } from "../../beats/toggleFollows.js";
import { meChat } from "../../mechat/plugnplay.js";
import { reportEntity } from "../../reporting/reporting.js";
import { logout } from "../../auth/authService.js";

import type { UserProfile } from "../profileGenHelpers.js";

export function FollowUser(followBtn: HTMLButtonElement, userId: string | number): void {
    toggleAction({
        entityId: userId,
        entityType: "user",
        button: followBtn,
        apiPath: "/subscribes/",
        labels: { on: "Unfollow", off: "Follow" },
        actionName: "followed"
    });
}

export function createProfileActions(profile: UserProfile, isLoggedIn: boolean): HTMLDivElement {
    const profileActions = document.createElement("div");
    profileActions.className = "profile-actions";

    // Prefer the auth alias `userid` which normalizes user id whether `user` is a string or object
    const currentUser = getState("userid") || (getState("user") && typeof getState("user") === "string" ? getState("user") : getState("user")?.userid);

    if (String(profile.userid) === String(currentUser)) {
        const logoutOptions: ButtonOptions = {
            title: "Logout",
            id: "logout-btn",
            events: { click: async () => await logout() },
            classes: "dropdown-item logout-btn"
        };
        const logoutButton = Button(logoutOptions);
        profileActions.appendChild(logoutButton);

        const editOptions: ButtonOptions = {
            title: "Edit Profile",
            id: "edit-profile-btn",
            classes: "btn edit-btn",
            "data-action": "edit-profile"
        };
        const editButton = Button(editOptions);
        profileActions.appendChild(editButton);
    }

    if (isLoggedIn && profile.userid !== undefined && String(profile.userid) !== String(currentUser)) {
        const followButton = Button({
            title: profile.is_following ? "Unfollow" : "Follow",
            id: "follow-btn",
            classes: "btn follow-button",
            styles: { backgroundColor: "green" },
            "data-action": "toggle-follow",
            "data-userid": String(profile.userid)
        });
        followButton.addEventListener("click", () => FollowUser(followButton, profile.userid as string | number));
        profileActions.appendChild(followButton);

        const sendMessageOptions: ButtonOptions = {
            title: "Send Message",
            id: "send-msg",
            events: {
                click: () => meChat(profile.userid as string | number, "user", currentUser)
            },
            classes: "buttonx"
        };
        const sendMessagebtn = Button(sendMessageOptions);
        profileActions.appendChild(sendMessagebtn);

        const reportOptions: ButtonOptions = {
            title: "Report",
            id: "report-btn",
            events: {
                click: () => reportEntity(String(profile.userid), "user")
            },
            classes: "report-btn"
        };
        const reportButton = Button(reportOptions);
        profileActions.appendChild(reportButton);
    }

    return profileActions;
}
