-----------------------------------------------------------------------------------------
--
-- views/login.lua
-- Unprotected Login View: Native Inputs, Validation & Navigation to Register
--
-----------------------------------------------------------------------------------------

local composer = require("composer")
local widget = require("widget")
local networkService = require("network.network")
local preferences = require("services.preferences")

local scene = composer.newScene()

local emailField
local passwordField
local statusText
local loginBtn
local registerLinkText

-----------------------------------------------------------------------------------------
-- Helper: Dismiss Keyboard
-----------------------------------------------------------------------------------------
local function dismissKeyboard()
    native.setKeyboardFocus(nil)
end

-----------------------------------------------------------------------------------------
-- Helper: Navigate to Register View
-----------------------------------------------------------------------------------------
local function goToRegister()
    dismissKeyboard()
    local router = require("router.router")
    router.switchView("views.register", "Register", "/register", false, true)
end

-----------------------------------------------------------------------------------------
-- Login Request Handler
-----------------------------------------------------------------------------------------
local function handleLogin()
    dismissKeyboard()

    local email = emailField and emailField.text or ""
    local password = passwordField and passwordField.text or ""

    -- Basic Validation
    if email == "" or password == "" then
        statusText.text = "Please enter both email and password."
        statusText:setFillColor(1, 0.3, 0.3)
        return
    end

    statusText.text = "Authenticating..."
    statusText:setFillColor(0.8, 0.8, 0.8)
    loginBtn:setEnabled(false)

    networkService.login(email, password, function(data, err)
        if not statusText or not statusText.parent then return end

        if err then
            statusText.text = "Error: " .. tostring(err)
            statusText:setFillColor(1, 0.3, 0.3)
            loginBtn:setEnabled(true)
        else
            statusText.text = "Success! Redirecting..."
            statusText:setFillColor(0.3, 1, 0.3)

            -- Save auth session locally
            preferences.saveSession(data.token, data.user)

            -- Navigate directly to the main feed view
            local router = require("router.router")
            router.switchView("views.feed", "Feed", "/feed", true, true)
        end
    end)
end

-----------------------------------------------------------------------------------------
-- Composer Lifecycle
-----------------------------------------------------------------------------------------
function scene:create(event)
    local sceneGroup = self.view

    -- Background touch listener to dismiss native keyboard
    local bg = display.newRect(sceneGroup, display.contentCenterX, display.contentCenterY, display.actualContentWidth, display.actualContentHeight)
    bg:setFillColor(0.08, 0.08, 0.1)
    bg:addEventListener("tap", dismissKeyboard)

    -- App Logo / Title
    local titleText = display.newText({
        parent = sceneGroup,
        text = "Welcome Back",
        x = display.contentCenterX,
        y = display.contentCenterY - 180,
        font = native.systemFontBold,
        fontSize = 28
    })
    titleText:setFillColor(1, 1, 1)

    local subtitleText = display.newText({
        parent = sceneGroup,
        text = "Sign in to access your dashboard",
        x = display.contentCenterX,
        y = display.contentCenterY - 145,
        font = native.systemFont,
        fontSize = 14
    })
    subtitleText:setFillColor(0.6, 0.6, 0.6)

    -- Native Email Field
    emailField = native.newTextField(display.contentCenterX, display.contentCenterY - 60, 260, 44)
    emailField.placeholder = "Email Address"
    emailField.inputType = "email"
    emailField.hasBackground = true
    sceneGroup:insert(emailField)

    -- Native Password Field
    passwordField = native.newTextField(display.contentCenterX, display.contentCenterY + 5, 260, 44)
    passwordField.placeholder = "Password"
    passwordField.isSecure = true
    passwordField.hasBackground = true
    sceneGroup:insert(passwordField)

    -- Status/Error Message Display
    statusText = display.newText({
        parent = sceneGroup,
        text = "",
        x = display.contentCenterX,
        y = display.contentCenterY + 60,
        width = 280,
        font = native.systemFont,
        fontSize = 13,
        align = "center"
    })

    -- Submit Button
    loginBtn = widget.newButton({
        label = "Sign In",
        x = display.contentCenterX,
        y = display.contentCenterY + 115,
        width = 260,
        height = 46,
        shape = "roundedRect",
        cornerRadius = 8,
        fillColor = { default = { 0.2, 0.45, 0.9 }, over = { 0.15, 0.35, 0.75 } },
        labelColor = { default = { 1, 1, 1 }, over = { 0.9, 0.9, 0.9 } },
        fontSize = 16,
        onRelease = handleLogin
    })
    sceneGroup:insert(loginBtn)

    -- Navigation Link to Register View
    registerLinkText = display.newText({
        parent = sceneGroup,
        text = "Don't have an account? Sign Up",
        x = display.contentCenterX,
        y = display.contentCenterY + 175,
        font = native.systemFont,
        fontSize = 14
    })
    registerLinkText:setFillColor(0.4, 0.7, 1)
    registerLinkText:addEventListener("tap", function()
        goToRegister()
        return true
    end)
end

function scene:show(event)
    if event.phase == "will" then
        if emailField then emailField.isVisible = true end
        if passwordField then passwordField.isVisible = true end
        if statusText then statusText.text = "" end
        if loginBtn then loginBtn:setEnabled(true) end
    end
end

function scene:hide(event)
    if event.phase == "will" then
        dismissKeyboard()
        if emailField then emailField.isVisible = false end
        if passwordField then passwordField.isVisible = false end
    end
end

function scene:destroy(event)
    if emailField then
        emailField:removeSelf()
        emailField = nil
    end
    if passwordField then
        passwordField:removeSelf()
        passwordField = nil
    end
end

scene:addEventListener("create", scene)
scene:addEventListener("show", scene)
scene:addEventListener("hide", scene)
scene:addEventListener("destroy", scene)

return scene