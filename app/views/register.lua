-----------------------------------------------------------------------------------------
--
-- views/register.lua
-- Unprotected Registration View: Account Creation & Navigation to Login
--
-----------------------------------------------------------------------------------------

local composer = require("composer")
local widget = require("widget")
local networkService = require("network.network")
local preferences = require("services.preferences")

local scene = composer.newScene()

local nameField
local emailField
local passwordField
local confirmPasswordField
local statusText
local registerBtn
local loginLinkText

-----------------------------------------------------------------------------------------
-- Helper: Dismiss Keyboard
-----------------------------------------------------------------------------------------
local function dismissKeyboard()
    native.setKeyboardFocus(nil)
end

-----------------------------------------------------------------------------------------
-- Helper: Navigate to Login View
-----------------------------------------------------------------------------------------
local function goToLogin()
    dismissKeyboard()
    local router = require("router.router")
    router.switchView("views.login", "Login", "/login", false, true)
end

-----------------------------------------------------------------------------------------
-- Registration Handler
-----------------------------------------------------------------------------------------
local function handleRegister()
    dismissKeyboard()

    local name = nameField and nameField.text or ""
    local email = emailField and emailField.text or ""
    local password = passwordField and passwordField.text or ""
    local confirmPassword = confirmPasswordField and confirmPasswordField.text or ""

    -- Client-side Validations
    if name == "" or email == "" or password == "" or confirmPassword == "" then
        statusText.text = "Please fill in all fields."
        statusText:setFillColor(1, 0.3, 0.3)
        return
    end

    if password ~= confirmPassword then
        statusText.text = "Passwords do not match."
        statusText:setFillColor(1, 0.3, 0.3)
        return
    end

    if #password < 6 then
        statusText.text = "Password must be at least 6 characters."
        statusText:setFillColor(1, 0.3, 0.3)
        return
    end

    statusText.text = "Creating account..."
    statusText:setFillColor(0.8, 0.8, 0.8)
    registerBtn:setEnabled(false)

    -- Send payload to registration endpoint
    local payload = {
        name = name,
        email = email,
        password = password
    }

    networkService.register(payload, function(data, err)
        if not statusText or not statusText.parent then return end

        if err then
            statusText.text = "Error: " .. tostring(err)
            statusText:setFillColor(1, 0.3, 0.3)
            registerBtn:setEnabled(true)
        else
            statusText.text = "Account created! Redirecting..."
            statusText:setFillColor(0.3, 1, 0.3)

            -- If backend returns JWT on register, save session & log user in automatically
            if data and data.token then
                preferences.saveSession(data.token, data.user or email)
                local router = require("router.router")
                router.switchView("views.feed", "Feed", "/feed", true, true)
            else
                -- Fallback: Send to login screen after 1.5s delay if manual login required
                timer.performWithDelay(1500, goToLogin)
            end
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

    -- Header / Title
    local titleText = display.newText({
        parent = sceneGroup,
        text = "Create Account",
        x = display.contentCenterX,
        y = display.contentCenterY - 210,
        font = native.systemFontBold,
        fontSize = 28
    })
    titleText:setFillColor(1, 1, 1)

    local subtitleText = display.newText({
        parent = sceneGroup,
        text = "Join us to get started",
        x = display.contentCenterX,
        y = display.contentCenterY - 175,
        font = native.systemFont,
        fontSize = 14
    })
    subtitleText:setFillColor(0.6, 0.6, 0.6)

    -- Native Name Field
    nameField = native.newTextField(display.contentCenterX, display.contentCenterY - 110, 260, 42)
    nameField.placeholder = "Full Name"
    nameField.hasBackground = true
    sceneGroup:insert(nameField)

    -- Native Email Field
    emailField = native.newTextField(display.contentCenterX, display.contentCenterY - 55, 260, 42)
    emailField.placeholder = "Email Address"
    emailField.inputType = "email"
    emailField.hasBackground = true
    sceneGroup:insert(emailField)

    -- Native Password Field
    passwordField = native.newTextField(display.contentCenterX, display.contentCenterY, 260, 42)
    passwordField.placeholder = "Password"
    passwordField.isSecure = true
    passwordField.hasBackground = true
    sceneGroup:insert(passwordField)

    -- Native Confirm Password Field
    confirmPasswordField = native.newTextField(display.contentCenterX, display.contentCenterY + 55, 260, 42)
    confirmPasswordField.placeholder = "Confirm Password"
    confirmPasswordField.isSecure = true
    confirmPasswordField.hasBackground = true
    sceneGroup:insert(confirmPasswordField)

    -- Status/Error Message Display
    statusText = display.newText({
        parent = sceneGroup,
        text = "",
        x = display.contentCenterX,
        y = display.contentCenterY + 105,
        width = 280,
        font = native.systemFont,
        fontSize = 13,
        align = "center"
    })

    -- Submit Register Button
    registerBtn = widget.newButton({
        label = "Sign Up",
        x = display.contentCenterX,
        y = display.contentCenterY + 155,
        width = 260,
        height = 46,
        shape = "roundedRect",
        cornerRadius = 8,
        fillColor = { default = { 0.2, 0.45, 0.9 }, over = { 0.15, 0.35, 0.75 } },
        labelColor = { default = { 1, 1, 1 }, over = { 0.9, 0.9, 0.9 } },
        fontSize = 16,
        onRelease = handleRegister
    })
    sceneGroup:insert(registerBtn)

    -- Direct to Login Screen Toggle
    loginLinkText = display.newText({
        parent = sceneGroup,
        text = "Already have an account? Sign In",
        x = display.contentCenterX,
        y = display.contentCenterY + 215,
        font = native.systemFont,
        fontSize = 14
    })
    loginLinkText:setFillColor(0.4, 0.7, 1)
    loginLinkText:addEventListener("tap", function()
        goToLogin()
        return true
    end)
end

function scene:show(event)
    if event.phase == "will" then
        if nameField then nameField.isVisible = true end
        if emailField then emailField.isVisible = true end
        if passwordField then passwordField.isVisible = true end
        if confirmPasswordField then confirmPasswordField.isVisible = true end
        if statusText then statusText.text = "" end
        if registerBtn then registerBtn:setEnabled(true) end
    end
end

function scene:hide(event)
    if event.phase == "will" then
        dismissKeyboard()
        if nameField then nameField.isVisible = false end
        if emailField then emailField.isVisible = false end
        if passwordField then passwordField.isVisible = false end
        if confirmPasswordField then confirmPasswordField.isVisible = false end
    end
end

function scene:destroy(event)
    if nameField then nameField:removeSelf(); nameField = nil end
    if emailField then emailField:removeSelf(); emailField = nil end
    if passwordField then passwordField:removeSelf(); passwordField = nil end
    if confirmPasswordField then confirmPasswordField:removeSelf(); confirmPasswordField = nil end
end

scene:addEventListener("create", scene)
scene:addEventListener("show", scene)
scene:addEventListener("hide", scene)
scene:addEventListener("destroy", scene)

return scene