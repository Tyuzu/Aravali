-----------------------------------------------------------------------------------------
--
-- main.lua
-- Application Entry Point
--
-----------------------------------------------------------------------------------------

local composer = require("composer")
local router = require("router.router")
local preferences = require("services.preferences")

-- Hide status bar on devices for full-screen layout
display.setStatusBar(display.HiddenStatusBar)

-----------------------------------------------------------------------------------------
-- Custom Splash Screen Execution
-----------------------------------------------------------------------------------------
-- 1. Create a black background to cover sizing differences
local splashBg = display.newRect(display.contentCenterX, display.contentCenterY, display.actualContentWidth, display.actualContentHeight)
splashBg:setFillColor(0, 0, 0) -- Change to your brand color hex if desired

-- 2. Load your custom splash logo centered on screen
-- Replace "my_custom_splash.png" with your actual file name
local splashLogo = display.newImageRect("images/my_custom_splash.png", 720, 720) 
splashLogo.x = display.contentCenterX
splashLogo.y = display.contentCenterY

-- Automatically purge unused display objects on scene changes to keep memory footprint low
composer.recycleOnSceneChange = true

-----------------------------------------------------------------------------------------
-- Hardware Key Listener (Android Physical Back Button Handling)
-----------------------------------------------------------------------------------------
local function onKeyEvent(event)
    -- Intercept the physical back button release ("up" phase)
    if event.keyName == "back" and event.phase == "up" then
        
        -- Step 1: Attempt to pop the top scene from the router history stack
        local handled = router.goBack()

        -- Step 2: If router.goBack() returned false, the user is on a root tab or login screen
        if not handled then
            native.showAlert(
                "Exit Application", 
                "Are you sure you want to exit?", 
                { "Cancel", "Exit" }, 
                function(e)
                    if e.action == "clicked" and e.index == 2 then
                        native.requestExit()
                    end
                end
            )
        end

        -- Return true to notify Solar2D that the event has been handled
        return true
    end

    -- Pass-through for other keys (volume buttons, etc.)
    return false
end

-- Register global runtime event listener for physical keys
Runtime:addEventListener("key", onKeyEvent)

-----------------------------------------------------------------------------------------
-- App Startup Chain (Deferred until Splash ends)
-----------------------------------------------------------------------------------------

local function startApp()
    -- 1. Initialize encrypted preferences store & load disk cache
    preferences.init()

    -- 2. Build Header & TabBar chrome UI instances
    router.init()

    -- 3. Evaluate auth state on boot to route user to correct screen
    if preferences.isLoggedIn() then
        -- Authenticated: Direct to Feed view
        router.switchView("views.feed", "Feed", "/feed", false, true)
    else
        -- Unauthenticated: Direct to Login screen
        router.switchView("views.login", "Login", "/login", false, true)
    end

    -- 4. Smoothly fade out and destroy the splash elements over 400ms
    transition.to(splashBg, { time = 400, alpha = 0, onComplete = display.remove })
    transition.to(splashLogo, { time = 400, alpha = 0, onComplete = display.remove })
end

-- Hold the splash screen visible for 2000 milliseconds (2 seconds) before executing startApp
timer.performWithDelay(2000, startApp)
