local App = {}

function App.start()
    display.setStatusBar(display.HiddenStatusBar)

    local bg = display.newRect(
        display.contentCenterX,
        display.contentCenterY,
        display.actualContentWidth,
        display.actualContentHeight
    )
    bg:setFillColor(0.95, 0.95, 0.97)

    local title = display.newText({
        text = "Hello World",
        x = display.contentCenterX,
        y = display.contentCenterY,
        font = native.systemFontBold,
        fontSize = 32
    })
    title:setFillColor(0.15, 0.18, 0.22)
end

return App
