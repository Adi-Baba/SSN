-- tools/version.lua
-- Production Version Management for SSN

local function read_file(path)
    local f = io.open(path, "r")
    if not f then return nil end
    local content = f:read("*all")
    f:close()
    return content
end

local function write_file(path, content)
    local f = io.open(path, "w")
    if not f then return false end
    f:write(content)
    f:close()
    return true
end

-- Find ssn.go robustly
local function find_target_file()
    local possible_paths = {
        "pkg/ssn/ssn.go",       -- From root
        "../pkg/ssn/ssn.go",    -- From tools/
        "../../pkg/ssn/ssn.go"  -- Deeply nested
    }
    for _, p in ipairs(possible_paths) do
        local f = io.open(p, "r")
        if f then
            f:close()
            return p
        end
    end
    return nil
end

local function parse_version(v_str)
    local major, minor, patch = v_str:match("(%d+)%.(%d+)%.(%d+)")
    if not major then return nil end
    return {
        major = tonumber(major),
        minor = tonumber(minor),
        patch = tonumber(patch),
        original = v_str
    }
end

local function format_version(v_obj)
    return string.format("%d.%d.%d", v_obj.major, v_obj.minor, v_obj.patch)
end

local function get_current_version(path)
    local content = read_file(path)
    if not content then return nil, "Read Error" end
    local v_str = content:match('const%s+Version%s+=%s+"([^"]+)"')
    if not v_str then return nil, "Pattern Not Found" end
    return parse_version(v_str), content
end

local function bump_version(v_obj, type)
    if type == "major" then
        v_obj.major = v_obj.major + 1
        v_obj.minor = 0
        v_obj.patch = 0
    elseif type == "minor" then
        v_obj.minor = v_obj.minor + 1
        v_obj.patch = 0
    elseif type == "patch" then
        v_obj.patch = v_obj.patch + 1
    else
        return false
    end
    return true
end

-- CLI Logic
local cmd = arg[1] or "status"
local subcmd = arg[2]

local file_path = find_target_file()
if not file_path then
    print("\27[31mError: Could not locate pkg/ssn/ssn.go\27[0m")
    os.exit(1)
end

local current_v, raw_content = get_current_version(file_path)
if not current_v then
    print("\27[31mError processing version file: " .. (raw_content or "Unknown") .. "\27[0m")
    os.exit(1)
end

if cmd == "status" then
    print("File:    " .. file_path)
    print("Version: \27[32m" .. format_version(current_v) .. "\27[0m")

elseif cmd == "bump" then
    if not subcmd then
        print("Usage: lua tools/version.lua bump [major|minor|patch]")
        os.exit(1)
    end
    
    local old_v_str = format_version(current_v)
    if bump_version(current_v, subcmd) then
        local new_v_str = format_version(current_v)
        local new_content = raw_content:gsub('const%s+Version%s+=%s+"[^"]+"', 'const Version = "' .. new_v_str .. '"')
        
        if write_file(file_path, new_content) then
            print("Bumped " .. subcmd .. ": \27[33m" .. old_v_str .. "\27[0m -> \27[32m" .. new_v_str .. "\27[0m")
        else
            print("\27[31mError: Failed to write to file.\27[0m")
            os.exit(1)
        end
    else
        print("\27[31mError: Invalid bump type. Use major, minor, or patch.\27[0m")
        os.exit(1)
    end

elseif cmd == "set" then
    if not subcmd then
        print("Usage: lua tools/version.lua set X.Y.Z")
        os.exit(1)
    end
    local new_v = parse_version(subcmd)
    if not new_v then
        print("Error: Invalid version format (must be x.y.z)")
        os.exit(1)
    end
    
    local new_v_str = format_version(new_v)
    local new_content = raw_content:gsub('const%s+Version%s+=%s+"[^"]+"', 'const Version = "' .. new_v_str .. '"')
    if write_file(file_path, new_content) then
        print("Set version: \27[32m" .. new_v_str .. "\27[0m")
    else
        print("Error writing file.")
    end

else
    print("Unknown command: " .. cmd)
    print("Available: status, bump [major|minor|patch], set [version]")
    os.exit(1)
end
