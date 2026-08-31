-----------------------------------------------------------------------------------------
--
-- router/router.lua
-- Navigation Manager: Header, TabBar, Route Guarding & Navigation Stack
--
-----------------------------------------------------------------------------------------

local widget = require("widget")
local composer = require("composer")
local networkService = require("network.network")
local preferences = require("services.preferences")

local Router = {}

-- Safe Area Geometry (Notch & Home Indicator Handling)
local topInset = display.safeInsetTop or 0
local bottomInset = display.safeInsetBottom or 0

local topEdge = display.screenOriginY
local bottomEdge = display.screenOriginY + display.actualContentHeight
local leftEdge = display.screenOriginX
local screenWidth = display.actualContentWidth

-- Adjust heights for safe area insets
local headerHeight = 60 + topInset
local tabBarHeight = 50 + bottomInset

-- State Tracking & Guarding
local history = {}            -- Stack array for back-navigation
local currentTab = 1
local isTransitioning = false -- Guard against tap spamming

local currentSceneName = "views.login"
local currentTitle = "Login"
local currentEndpoint = "/login"

-- UI References
local headerGroup
local headerTitle
local refreshBtn
local backBtn
local tabBar

-- Map root view modules directly to their tab bar indices
local TAB_MAP = {
    ["views.feed"]      = 1,
    ["views.profile"]   = 2,
    ["views.analytics"] = 3,
    ["views.settings"]  = 4
}

-----------------------------------------------------------------------------------------
-- Header Component
-----------------------------------------------------------------------------------------
local function createHeader()
    headerGroup = display.newGroup()

    -- Background
    local headerBg = display.newRect(
        headerGroup,
        leftEdge + screenWidth * 0.5,
        topEdge + headerHeight * 0.5,
        screenWidth,
        headerHeight
    )
    headerBg:setFillColor(0.12, 0.12, 0.15)

    -- Back Button (Hidden on root views)
    backBtn = display.newText({
        parent = headerGroup,
        text = "←",
        x = leftEdge + 25,
        y = topEdge + topInset + (headerHeight - topInset) * 0.5,
        font = native.systemFontBold,
        fontSize = 24
    })
    backBtn:setFillColor(1, 1, 1)
    backBtn.isVisible = false

    backBtn:addEventListener("tap", function()
        Router.goBack()
        return true
    end)

    -- Header Title
    headerTitle = display.newText({
        parent = headerGroup,
        text = "Connecting...",
        x = leftEdge + screenWidth * 0.5,
        y = topEdge + topInset + (headerHeight - topInset) * 0.5,
        font = native.systemFontBold,
        fontSize = 18
    })
    headerTitle:setFillColor(1, 1, 1)

    -- Refresh Button
    refreshBtn = display.newText({
        parent = headerGroup,
        text = "🔄",
        x = leftEdge + screenWidth - 25,
        y = topEdge + topInset + (headerHeight - topInset) * 0.5,
        font = native.systemFont,
        fontSize = 20
    })

    refreshBtn:addEventListener("tap", function()
        if not isTransitioning then
            transition.to(refreshBtn, { time = 300, rotation = refreshBtn.rotation + 360 })
            Router.switchView(currentSceneName, currentTitle, currentEndpoint, true, false)
        end
        return true
    end)

    -- Initially hidden until authenticated
    headerGroup.isVisible = false
end

-----------------------------------------------------------------------------------------
-- TabBar Component
-----------------------------------------------------------------------------------------
local function createTabBar()
    local tabButtons = {
        {
            label = "Feed",
            defaultFile = "buttons/button1.png",
            overFile = "buttons/button1-down.png",
            width = 32, height = 32, selected = true,
            onPress = function() Router.switchView("views.feed", "Feed", "/feed", false, true) end
        },
        {
            label = "Profile",
            defaultFile = "buttons/button2.png",
            overFile = "buttons/button2-down.png",
            width = 32, height = 32,
            onPress = function() Router.switchView("views.profile", "User Profile", "/profile", false, true) end
        },
        {
            label = "Analytics",
            defaultFile = "buttons/button3.png",
            overFile = "buttons/button3-down.png",
            width = 32, height = 32,
            onPress = function() Router.switchView("views.analytics", "Analytics", "/analytics", false, true) end
        },
        {
            label = "Settings",
            defaultFile = "buttons/button4.png",
            overFile = "buttons/button4-down.png",
            width = 32, height = 32,
            onPress = function() Router.switchView("views.settings", "App Settings", "/settings", false, true) end
        }
    }

    tabBar = widget.newTabBar({
        left = leftEdge,
        top = bottomEdge - tabBarHeight,
        width = screenWidth,
        height = tabBarHeight,
        buttons = tabButtons
    })

    -- Initially hidden until authenticated
    tabBar.isVisible = false
end

-----------------------------------------------------------------------------------------
-- Core Navigation Logic
-----------------------------------------------------------------------------------------

--- Checks whether a target scene is an unauthenticated public route
local function isUnprotectedView(sceneName)
    return (sceneName == "views.login" or sceneName == "views.register")
end

--- Updates Header & TabBar visibility based on current scene
local function updateChromeVisibility(sceneName)
    local hideChrome = isUnprotectedView(sceneName)
    
    if headerGroup then
        headerGroup.isVisible = not hideChrome
    end
    if tabBar then
        tabBar.isVisible = not hideChrome
    end
end

--- Navigate to a target scene with route guarding, history tracking, and state safety
function Router.switchView(sceneName, titleText, endpoint, forceRefresh, isRootTab)
    if isTransitioning then return end

    ---------------------------------------------------------------------------------
    -- AUTH GUARD: Redirect to Login if unauthenticated & accessing protected views
    ---------------------------------------------------------------------------------
    local isLoggedIn = preferences.isLoggedIn()
    if not isLoggedIn and not isUnprotectedView(sceneName) then
        print("[Router Guard] Unauthenticated access attempt to " .. tostring(sceneName) .. ". Redirecting to Login.")
        sceneName = "views.login"
        titleText = "Login"
        endpoint = "/login"
        isRootTab = true
    end

    isTransitioning = true

    -- Handle Navigation History Stack
    if isRootTab then
        history = {} -- Reset stack when jumping directly to a main tab or logging out
    elseif currentSceneName and currentSceneName ~= sceneName then
        table.insert(history, {
            scene = currentSceneName,
            title = currentTitle,
            endpoint = currentEndpoint,
            tabIndex = currentTab
        })
    end

    -- Update state tracking
    currentSceneName = sceneName
    currentTitle = titleText
    currentEndpoint = endpoint
    currentTab = TAB_MAP[sceneName] or currentTab

    -- Update Header/TabBar visibility
    updateChromeVisibility(sceneName)

    -- Synchronize UI Tab Selection & Back Button Visibility
    if tabBar and TAB_MAP[sceneName] then
        tabBar:setSelected(TAB_MAP[sceneName], false)
    end
    if backBtn then
        backBtn.isVisible = (#history > 0)
    end

    ---------------------------------------------------------------------------------
    -- DIRECT TRANSITION: Public views (Login/Register) do not pre-fetch API endpoints
    ---------------------------------------------------------------------------------
    if isUnprotectedView(sceneName) then
        composer.gotoScene(sceneName, { effect = "crossFade", time = 150 })
        isTransitioning = false
        return
    end

    ---------------------------------------------------------------------------------
    -- PROTECTED TRANSITION: Fetch data endpoint before changing scene
    ---------------------------------------------------------------------------------
    if headerTitle then
        headerTitle.text = "Updating..."
    end

    networkService.fetch(endpoint, forceRefresh, function(data, err, isFromCache)
        if headerTitle then
            headerTitle.text = titleText
        end

        local options = {
            effect = "crossFade",
            time = 150,
            params = { apiData = data, error = err, isCached = isFromCache }
        }

        composer.gotoScene(sceneName, options)
        isTransitioning = false
    end)
end

--- Pop top scene from history stack and navigate back
function Router.goBack()
    if isTransitioning or #history == 0 then
        return false
    end

    if not preferences.isLoggedIn() then
        Router.switchView("views.login", "Login", "/login", false, true)
        return true
    end

    local previous = table.remove(history)
    
    isTransitioning = true
    currentSceneName = previous.scene
    currentTitle = previous.title
    currentEndpoint = previous.endpoint
    currentTab = previous.tabIndex

    updateChromeVisibility(previous.scene)

    if tabBar and previous.tabIndex then
        tabBar:setSelected(previous.tabIndex, false)
    end
    if backBtn then
        backBtn.isVisible = (#history > 0)
    end

    if headerTitle then
        headerTitle.text = "Loading..."
    end

    networkService.fetch(previous.endpoint, false, function(data, err, isFromCache)
        if headerTitle then
            headerTitle.text = previous.title
        end

        composer.gotoScene(previous.scene, {
            effect = "slideRight",
            time = 200,
            params = { apiData = data, error = err, isCached = isFromCache }
        })
        isTransitioning = false
    end)

    return true
end

-----------------------------------------------------------------------------------------
-- Initialization
-----------------------------------------------------------------------------------------
function Router.init()
    createHeader()
    createTabBar()
end

return Router