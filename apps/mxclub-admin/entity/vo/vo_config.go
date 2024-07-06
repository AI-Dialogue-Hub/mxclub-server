package vo

import "mxclub/pkg/common/xmysql"

type MiniConfigVO struct {
	ID          uint             `json:"id"`
	ConfigName  string           `json:"config_name"`
	DisPlayName string           `json:"display_name"`
	Content     xmysql.JSONArray `json:"content"`
}

type NotificationsVO struct {
	ID      uint   `json:"id"`
	Image   string `json:"image"`
	Message string `json:"message"`
}
