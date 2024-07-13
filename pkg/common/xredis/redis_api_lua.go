package xredis

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
	"time"
)

// Lua 脚本，删除所有模糊匹配的ky, 使用传入的 wildcardKey 进行键模式匹配
var delMatchingKeysScript = `
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
`

var delMatchingKeysScriptHash string

// DelMatchingKeys 函数，传入 wildcardKey 参数
func DelMatchingKeys(ctx jet.Ctx, wildcardKey string) error {
	defer utils.TraceElapsedByName(time.Now(), "DelMatchingKeys")
	if delMatchingKeysScriptHash == "" {
		// 加载脚本
		delMatchingKeysScriptHash = cli.GetSingleClient().ScriptLoad(context.Background(), delMatchingKeysScript).Val()
	}
	// 传递参数给 Eval 方法
	cmd := cli.GetSingleClient().EvalSha(context.Background(), delMatchingKeysScriptHash, []string{}, RealKey(wildcardKey)+"*")
	if err := cmd.Err(); err != nil {
		ctx.Logger().Errorf("deleting wildcard keys failed: %s", err.Error())
		return err
	}
	ctx.Logger().Infof("deleting wildcard keys success, %v", cmd.Val())
	return nil
}

// DelKeys 函数，传入 wildcardKey 参数
func DelKeys(keys ...string) error {
	for _, key := range keys {
		if err := cli.Del(key); err != nil {
			return err
		}
	}
	return nil
}

// DelByPattern 模糊删除缓存
func DelByPattern(ctx jet.Ctx, pattern string) (result int64) {
	keys, _ := cli.GetSingleClient().Keys(context.Background(), RealKey(pattern)+"*").Result()
	for i := 0; i < len(keys); i++ {
		cmd := cli.GetSingleClient().Del(context.Background(), keys[i])
		if cmd.Err() != nil {
			ctx.Logger().Errorf("[DelByPattern] del key %s fail, error: %v", keys[i], cmd.Err().Error())
		} else {
			result += cmd.Val()
		}
	}
	return
}
