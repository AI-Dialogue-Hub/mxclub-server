package vo

import "mxclub/pkg/common/xmysql"

type MiniConfigVO struct {
	ID          uint             `json:"id"`
	ConfigName  string           `json:"config_name"`
	DisPlayName string           `json:"dis_play_name"`
	Content     xmysql.JSONArray `json:"content"`
}
