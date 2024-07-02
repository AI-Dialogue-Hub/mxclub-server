-- 删除模糊匹配的键
-- @param pattern 键的模式
local function deleteKeysByPattern(pattern)
    local cursor = "0"
    local totalDeleted = 0
    repeat
        local result = redis.call("SCAN", cursor, "MATCH", pattern, "COUNT", 100)
        cursor = result[1]
        local keys = result[2]
        if #keys > 0 then
            for i, key in ipairs(keys) do
                redis.call("DEL", key)
                totalDeleted = totalDeleted + 1
            end
        end
    until cursor == "0"
    return totalDeleted
end

local pattern = ARGV[1]
return deleteKeysByPattern(pattern)
