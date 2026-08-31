-----------------------------------------------------------------------------------------
--
-- services/preferences.lua
-- Data Persistence & Encrypted Preference Storage
--
-----------------------------------------------------------------------------------------

local json = require("json")

local M = {}

-- File where data will be saved
local FILE_NAME = "app_preferences.json"
local FILE_PATH = system.pathForFile(FILE_NAME, system.DocumentsDirectory)

-- Hardcoded default application state
local DEFAULT_SETTINGS = {
    authToken = "",
    userEmail = "",
    isLoggedIn = false,
    theme = "dark",
    notificationsEnabled = true,
    volume = 0.8
}

-- In-memory cache table
local cache = nil

--------------------------------------------------------------------------------
-- Utilities
--------------------------------------------------------------------------------

local function copyTable(src)
    local dest = {}
    for k, v in pairs(src) do
        if type(v) == "table" then
            dest[k] = copyTable(v)
        else
            dest[k] = v
        end
    end
    return dest
end

--------------------------------------------------------------------------------
-- Obfuscation Helpers (Protects sensitive data from plain-text inspection)
--------------------------------------------------------------------------------
local SECRET_KEY = "YourAppSecretKeyHere"

-- Universal bitwise XOR helper safe across all Lua environments
local function xorByte(a, b)
    local r = 0
    for i = 0, 7 do
        local p = 2^i
        local a_bit = (a / p) % 2 >= 1
        local b_bit = (b / p) % 2 >= 1
        if a_bit ~= b_bit then
            r = r + p
        end
    end
    return r
end

local function obfuscate(inputStr)
    if not inputStr or inputStr == "" then return "" end
    local output = {}
    for i = 1, #inputStr do
        local charCode = string.byte(inputStr, i)
        local keyByte = string.byte(SECRET_KEY, ((i - 1) % #SECRET_KEY) + 1)
        output[i] = string.char(xorByte(charCode, keyByte))
    end
    return table.concat(output)
end

--------------------------------------------------------------------------------
-- Core Disk I/O
--------------------------------------------------------------------------------
local function saveToDisk()
    local file, err = io.open(FILE_PATH, "w")
    if file then
        local encodedData = json.encode(cache)
        file:write(encodedData)
        io.close(file)
        return true
    else
        print("Preferences Error: Failed to write to disk - " .. tostring(err))
        return false
    end
end

--------------------------------------------------------------------------------
-- Public API
--------------------------------------------------------------------------------

function M.init()
    local file = io.open(FILE_PATH, "r")
    
    if file then
        local contents = file:read("*a")
        io.close(file)
        
        local decoded, pos, msg = json.decode(contents)
        if decoded then
            cache = decoded
        else
            print("Preferences Error: JSON decode failed (" .. tostring(msg) .. "). Resetting defaults.")
            cache = copyTable(DEFAULT_SETTINGS)
            saveToDisk()
        end
    else
        cache = copyTable(DEFAULT_SETTINGS)
        saveToDisk()
    end
end

-- Ensures cache is loaded before any read/write operations
local function ensureInitialized()
    if not cache then
        M.init()
    end
end

function M.get(key)
    ensureInitialized()
    if cache[key] ~= nil then
        return cache[key]
    end
    return DEFAULT_SETTINGS[key]
end

function M.set(key, value)
    ensureInitialized()
    cache[key] = value
    saveToDisk()
end

--------------------------------------------------------------------------------
-- Secure Auth & Session Helpers
--------------------------------------------------------------------------------

--- Save active session parameters
function M.saveSession(token, email)
    ensureInitialized()
    if token and token ~= "" then
        cache.authToken = obfuscate(token)
        cache.isLoggedIn = true
    else
        cache.authToken = ""
        cache.isLoggedIn = false
    end

    if email and email ~= "" then
        cache.userEmail = obfuscate(email)
    end

    saveToDisk()
end

--- Wipe current active session parameters
function M.clearSession()
    ensureInitialized()
    cache.authToken = ""
    cache.userEmail = ""
    cache.isLoggedIn = false
    saveToDisk()
end

--- Legacy wrapper for clearSession
function M.logout()
    M.clearSession()
end

--- Get de-obfuscated auth token
function M.getAuthToken()
    ensureInitialized()
    if cache.authToken and cache.authToken ~= "" then
        return obfuscate(cache.authToken)
    end
    return nil
end

--- Get de-obfuscated user email
function M.getUserEmail()
    ensureInitialized()
    if cache.userEmail and cache.userEmail ~= "" then
        return obfuscate(cache.userEmail)
    end
    return nil
end

--- Check authentication state
function M.isLoggedIn()
    ensureInitialized()
    return cache.isLoggedIn == true and (cache.authToken and cache.authToken ~= "")
end

return M