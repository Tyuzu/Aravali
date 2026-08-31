-----------------------------------------------------------------------------------------
--
-- network/network.lua
-- Central Network Client: Data Fetching, JWT Headers, Cache Management & 401 Interception
--
-----------------------------------------------------------------------------------------
local json = require("json")
local preferences = require("services.preferences")

local Network = {}

-- Base URL configuration (adjust IP/port if running on local device/emulator)
-- local BASE_URL = "http://localhost:8080"
local BASE_URL = "http://ec2-16-171-215-85.eu-north-1.compute.amazonaws.com:7000"
-- ec2-16-171-215-85.eu-north-1.compute.amazonaws.com

-- Memory cache store
local memoryCache = {}

-----------------------------------------------------------------------------------------
-- Generic HTTP Request Executor
-----------------------------------------------------------------------------------------
local function httpRequest(endpoint, method, payload, callback)
    local url = BASE_URL .. endpoint

    local headers = {}
    headers["Content-Type"] = "application/json"
    headers["Accept"] = "application/json"

    -- Inject JWT Authorization token if available
    local token = preferences.getAuthToken()
    if token and token ~= "" then
        headers["Authorization"] = "Bearer " .. token
    end

    local params = {
        headers = headers,
        timeout = 10
    }

    if payload then
        params.body = json.encode(payload)
    end

    network.request(url, method, function(event)
        if event.isError then
            print("[Network Error] Connection failed for " .. endpoint)
            callback(nil, "Network connection error", false)
            return
        end

        local responseData = nil
        if event.response then
            responseData = json.decode(event.response)
        end

        -- Handle HTTP 401 Unauthorized globally
        if event.status == 401 then
            print("[Network Auth] 401 Unauthorized returned from " .. endpoint .. ". Logging out.")
            preferences.clearSession()

            -- Lazy load router to avoid circular dependency
            local router = require("router.router")
            router.switchView("views.login", "Login", "/login", false, true)

            callback(nil, "Session expired. Please log in again.", false)
            return
        end

        -- Handle non-200 HTTP statuses
        if event.status < 200 or event.status >= 300 then
            local errMsg = (responseData and responseData.message) or
                               ("Server error (" .. tostring(event.status) .. ")")
            callback(nil, errMsg, false)
            return
        end

        -- Cache valid GET responses
        if method == "GET" and responseData and responseData.success then
            memoryCache[endpoint] = {
                data = responseData.data,
                timestamp = os.time()
            }
        end

        local outputData = responseData and responseData.data or responseData
        callback(outputData, nil, false)
    end, params)
end

-----------------------------------------------------------------------------------------
-- Public API Methods
-----------------------------------------------------------------------------------------

--- Send account payload to /register endpoint
function Network.register(payload, callback)
    httpRequest("/register", "POST", payload, function(data, err)
        if err then
            callback(nil, err)
        else
            callback(data, nil)
        end
    end)
end

--- Send authentication credentials to /login endpoint
function Network.login(email, password, callback)
    local payload = {
        email = email,
        password = password
    }

    httpRequest("/login", "POST", payload, function(data, err)
        if err then
            callback(nil, err)
        else
            callback(data, nil)
        end
    end)
end

--- Fetch data with caching support
function Network.fetch(endpoint, forceRefresh, callback)
    -- Skip caching for authentication routes
    if endpoint == "/login" then
        callback(nil, nil, false)
        return
    end

    -- Return memory cached version if available and not forced
    if not forceRefresh and memoryCache[endpoint] then
        print("[Network] Returning memory cached data for " .. endpoint)
        callback(memoryCache[endpoint].data, nil, true)
        return
    end

    print("[Network] Fetching fresh data from " .. endpoint)
    httpRequest(endpoint, "GET", nil, callback)
end

--- Clear memory cache (useful on logout)
function Network.clearCache()
    memoryCache = {}
end

return Network
