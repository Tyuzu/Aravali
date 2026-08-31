-----------------------------------------------------------------------------------------
--
-- views/analytics.lua
-- Analytics Dashboard View
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
    bg:setFillColor(0.94, 0.94, 0.96)

    -------------------------------------------------------------------------------------
    -- Error Handling
    -------------------------------------------------------------------------------------
    if params.error then
        local errorCard = display.newRect(
            sceneGroup,
            leftEdge + screenWidth * 0.5,
            topEdge + headerHeight + contentAreaHeight * 0.5,
            screenWidth - 40,
            120
        )
        errorCard:setFillColor(1, 1, 1)

        local errorText = display.newText({
            parent = sceneGroup,
            text = "⚠️ Connection Error:\n\n" .. tostring(params.error),
            x = errorCard.x,
            y = errorCard.y,
            width = errorCard.width - 20,
            font = native.systemFont,
            fontSize = 15,
            align = "center"
        })
        errorText:setFillColor(0.8, 0.2, 0.2)
        return
    end

    -------------------------------------------------------------------------------------
    -- Metrics Grid
    -------------------------------------------------------------------------------------
    local apiData = params.apiData or {}
    local metrics = {
        { title = "DAILY ACTIVE USERS", value = tostring(apiData.dailyActiveUsers or 0), color = {0.2, 0.5, 0.9} },
        { title = "TOTAL REVENUE", value = "$" .. string.format("%.2f", apiData.totalRevenue or 0), color = {0.2, 0.7, 0.4} },
        { title = "SERVER UPTIME", value = tostring(apiData.serverUptime or "N/A"), color = {0.9, 0.5, 0.1} },
        { title = "CONVERSION RATE", value = tostring(apiData.conversionRate or "N/A"), color = {0.6, 0.3, 0.8} }
    }

    local cardWidth = (screenWidth - 50) * 0.5
    local cardHeight = 110
    local startY = topEdge + headerHeight + 80

    for i, item in ipairs(metrics) do
        local col = (i - 1) % 2
        local row = math.floor((i - 1) / 2)

        local posX = leftEdge + 20 + (col * (cardWidth + 10)) + (cardWidth * 0.5)
        local posY = startY + (row * (cardHeight + 15))

        local card = display.newRect(sceneGroup, posX, posY, cardWidth, cardHeight)
        card:setFillColor(1, 1, 1)
        card:setStrokeColor(0.88, 0.88, 0.9)
        card.strokeWidth = 1

        local label = display.newText({
            parent = sceneGroup,
            text = item.title,
            x = posX,
            y = posY - 25,
            font = native.systemFontBold,
            fontSize = 11
        })
        label:setFillColor(0.5, 0.5, 0.55)

        local val = display.newText({
            parent = sceneGroup,
            text = item.value,
            x = posX,
            y = posY + 10,
            font = native.systemFontBold,
            fontSize = 18
        })
        val:setFillColor(unpack(item.color))
    end
end

scene:addEventListener("create", scene)
return scene