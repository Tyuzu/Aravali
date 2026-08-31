-----------------------------------------------------------------------------------------
--
-- views/profile.lua
-- User Profile View Dashboard
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
    -- Background Frame
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
    -- Error Handling State
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
    -- Profile Data Extraction
    -------------------------------------------------------------------------------------
    local apiData = params.apiData or {}
    local userName = apiData.name or "Unknown User"
    local userEmail = apiData.email or "No email provided"
    local userRole = apiData.role or "Member"

    -------------------------------------------------------------------------------------
    -- Profile Header Card
    -------------------------------------------------------------------------------------
    local cardYCenter = topEdge + headerHeight + 110

    local profileCard = display.newRect(
        sceneGroup,
        leftEdge + screenWidth * 0.5,
        cardYCenter,
        screenWidth - 40,
        160
    )
    profileCard:setFillColor(1, 1, 1)
    profileCard:setStrokeColor(0.88, 0.88, 0.90)
    profileCard.strokeWidth = 1

    local avatarCircle = display.newCircle(
        sceneGroup,
        profileCard.x,
        cardYCenter - 30,
        32
    )
    avatarCircle:setFillColor(0.2, 0.45, 0.85)

    local initials = string.sub(userName, 1, 1)
    local avatarText = display.newText({
        parent = sceneGroup,
        text = initials,
        x = avatarCircle.x,
        y = avatarCircle.y,
        font = native.systemFontBold,
        fontSize = 24
    })
    avatarText:setFillColor(1, 1, 1)

    local nameText = display.newText({
        parent = sceneGroup,
        text = userName,
        x = profileCard.x,
        y = cardYCenter + 18,
        font = native.systemFontBold,
        fontSize = 18
    })
    nameText:setFillColor(0.15, 0.15, 0.2)

    local roleText = display.newText({
        parent = sceneGroup,
        text = userRole,
        x = profileCard.x,
        y = cardYCenter + 42,
        font = native.systemFont,
        fontSize = 14
    })
    roleText:setFillColor(0.45, 0.45, 0.5)

    -------------------------------------------------------------------------------------
    -- Details List Card
    -------------------------------------------------------------------------------------
    local detailsYCenter = cardYCenter + 140

    local detailsCard = display.newRect(
        sceneGroup,
        leftEdge + screenWidth * 0.5,
        detailsYCenter,
        screenWidth - 40,
        90
    )
    detailsCard:setFillColor(1, 1, 1)
    detailsCard:setStrokeColor(0.88, 0.88, 0.90)
    detailsCard.strokeWidth = 1

    local emailLabel = display.newText({
        parent = sceneGroup,
        text = "EMAIL ADDRESS",
        x = detailsCard.x - (detailsCard.width * 0.5) + 20,
        y = detailsYCenter - 20,
        font = native.systemFontBold,
        fontSize = 11
    })
    emailLabel:setFillColor(0.55, 0.55, 0.6)
    emailLabel.anchorX = 0

    local emailValue = display.newText({
        parent = sceneGroup,
        text = userEmail,
        x = detailsCard.x - (detailsCard.width * 0.5) + 20,
        y = detailsYCenter + 5,
        font = native.systemFont,
        fontSize = 15
    })
    emailValue:setFillColor(0.2, 0.2, 0.25)
    emailValue.anchorX = 0
end

scene:addEventListener("create", scene)
return scene