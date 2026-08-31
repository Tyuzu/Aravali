-----------------------------------------------------------------------------------------
--
-- views/feed.lua
-- Dynamic Feed View Area
--
-----------------------------------------------------------------------------------------

local composer = require("composer")
local scene = composer.newScene()

function scene:create(event)
    local sceneGroup = self.view
    local params = event.params or {}

    local topEdge = display.screenOriginY
    local leftEdge = display.screenOriginX
    local screenWidth = display.actualContentWidth
    local screenHeight = display.actualContentHeight

    local headerHeight = 60
    local tabBarHeight = 50
    local contentAreaHeight = screenHeight - (headerHeight + tabBarHeight)

    -------------------------------------------------------------------------------------
    -- Background
    -------------------------------------------------------------------------------------
    local bg = display.newRect(
        sceneGroup,
        leftEdge + screenWidth * 0.5,
        topEdge + headerHeight + contentAreaHeight * 0.5,
        screenWidth,
        contentAreaHeight
    )
    bg:setFillColor(0.95, 0.95, 0.97)

    -------------------------------------------------------------------------------------
    -- Content Card Box
    -------------------------------------------------------------------------------------
    local card = display.newRect(
        sceneGroup,
        leftEdge + screenWidth * 0.5,
        topEdge + headerHeight + contentAreaHeight * 0.5,
        screenWidth - 40,
        contentAreaHeight - 40
    )
    card:setFillColor(1, 1, 1)
    card:setStrokeColor(0.85, 0.85, 0.88)
    card.strokeWidth = 1

    -------------------------------------------------------------------------------------
    -- Render Data / Error
    -------------------------------------------------------------------------------------
    local displayText = ""

    if params.error then
        displayText = "⚠️ Connection Error:\n\n" .. tostring(params.error)
    elseif params.apiData then
        if params.apiData.title then
            displayText = "📌 " .. tostring(params.apiData.title) .. "\n\n" .. tostring(params.apiData.body)
        elseif params.apiData.name then
            displayText = "👤 " .. tostring(params.apiData.name) .. "\n" ..
                          "📧 " .. tostring(params.apiData.email) .. "\n" ..
                          "💼 Role: " .. tostring(params.apiData.role)
        else
            displayText = "Data received, but unmapped format."
        end
    else
        displayText = "No content available."
    end

    local contentLabel = display.newText({
        parent = sceneGroup,
        text = displayText,
        x = card.x,
        y = card.y,
        width = card.width - 40,
        font = native.systemFont,
        fontSize = 16,
        align = "center"
    })
    contentLabel:setFillColor(0.2, 0.2, 0.25)
end

scene:addEventListener("create", scene)
return scene