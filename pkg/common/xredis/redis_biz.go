// 常用的redis功能

package xredis

import (
	"errors"
	"fmt"
	"mxclub/pkg/utils"
	"time"
)

func DebounceForOneDay(key string) error {
	duration := 24 * time.Hour
	return Debounce(key, duration)
}

// Debounce 使用redis实现防抖
//
// SET key value [EX seconds] [PX milliseconds] [NX|XX]
//
// @param key 幂等key
//
// @param duration 过期时间
func Debounce(key string, duration time.Duration) error {
	defer utils.RecoverByPrefixNoCtx("Debounce")
	ok, err := cli.SetNX(key, fmt.Sprintf("%v_%v", "debounce", key), duration)
	if err != nil {
		return fmt.Errorf("error setting debounce key: %w", err)
	}
	if !ok {
		// 如果 key 已经存在，SetNX 操作会返回 false，表示这是重复触发
		return errors.New("duplicated")
	}
	return nil
}
