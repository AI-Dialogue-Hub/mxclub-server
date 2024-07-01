package xredis

// Lua 脚本，删除所有模糊匹配的ky, 使用传入的 wildcardKey 进行键模式匹配
var script = `
    local keys = redis.call('KEYS', KEYS[1])
    for _, key in ipairs(keys) do
        redis.call('DEL', key)
    end
`

// DelMatchingKeys 函数，传入 wildcardKey 参数
func DelMatchingKeys(wildcardKey string) error {
	// 传递参数给 Eval 方法
	return cli.Eval(script, []string{wildcardKey})
}
