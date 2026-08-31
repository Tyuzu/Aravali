-----------------------------------------------------------------------------------------
--
-- views/settings.lua
-- Settings View: Displays App Config & Actionable Logout Button
--
-----------------------------------------------------------------------------------------

local composer = require("composer")
local widget = require("widget")
local preferences = require("services.preferences")
local networkService = require("network.network")

local scene = composer.newScene()

-----------------------------------------------------------------------------------------
-- Logout Logic
-----------------------------------------------------------------------------------------
local function handleLogout()
    -- 1. Wipe local token session
    preferences.clearSession()

    -- 2. Wipe network memory cache
    networkService.clearCache()

    -- 3. Redirect to login via router
    local router = require("router.router")
    router.switchView("views.login", "Login", "/login", false, true)
end

-----------------------------------------------------------------------------------------
-- Composer Lifecycle
-----------------------------------------------------------------------------------------
function scene:create(event)
    local sceneGroup = self.view

    -- Background
    local bg = display.newRect(sceneGroup, display.contentCenterX, display.contentCenterY, display.actualContentWidth, display.actualContentHeight)
    bg:setFillColor(0.1, 0.1, 0.12)
end

function scene:show(event)
    local sceneGroup = self.view

    if event.phase == "will" then
        -- Remove stale dynamic content elements before re-rendering
        for i = sceneGroup.numChildren, 1, -1 do
            if sceneGroup[i].isDynamic then
                sceneGroup[i]:removeSelf()
            end
        end

        local apiData = event.params and event.params.apiData or {}
        local currentUser = preferences.getUserEmail() or "User"

        -- Current User Banner
        local userLabel = display.newText({
            parent = sceneGroup,
            text = "Logged in as: " .. currentUser,
            x = display.contentCenterX,
            y = display.screenOriginY + 90,
            font = native.systemFontBold,
            fontSize = 15
        })
        userLabel:setFillColor(0.4, 0.8, 1)
        userLabel.isDynamic = true

        -- Config Values Card Container
        local card = display.newRoundedRect(
            sceneGroup,
            display.contentCenterX,
            display.screenOriginY + 200,
            display.actualContentWidth - 40,
            160,
            8
        )
        card:setFillColor(0.16, 0.16, 0.2)
        card.isDynamic = true

        -- Render API Settings Parameters
        local settingsList = {
            "API Version: " .. tostring(apiData.apiVersion or "N/A"),
            "Theme: " .. tostring(apiData.theme or "N/A"),
            "Notifications: " .. (apiData.notificationsEnabled and "Enabled" or "Disabled"),
            "Dark Mode: " .. (apiData.darkMode and "Active" or "Inactive")
        }

        for i, textValue in ipairs(settingsList) do
            local itemText = display.newText({
                parent = sceneGroup,
                text = textValue,
                x = display.contentCenterX - 110,
                y = display.screenOriginY + 140 + (i * 24),
                font = native.systemFont,
                fontSize = 14,
                align = "left"
            })
            itemText.anchorX = 0
            itemText:setFillColor(0.9, 0.9, 0.9)
            itemText.isDynamic = true
        end

        -- Logout Button
        local logoutBtn = widget.newButton({
            label = "Log Out",
            x = display.contentCenterX,
            y = display.screenOriginY + 320,
            width = display.actualContentWidth - 40,
            height = 46,
            shape = "roundedRect",
            cornerRadius = 8,
            fillColor = { default = { 0.8, 0.2, 0.2 }, over = { 0.6, 0.15, 0.15 } },
            labelColor = { default = { 1, 1, 1 }, over = { 0.9, 0.9, 0.9 } },
            fontSize = 16,
            onRelease = handleLogout
        })
        logoutBtn.isDynamic = true
        sceneGroup:insert(logoutBtn)
    end
end

scene:addEventListener("create", scene)
scene:addEventListener("show", scene)

return scene