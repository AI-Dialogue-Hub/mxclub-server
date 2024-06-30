package vo

type MiniConfigVO struct {
	ID         uint           `json:"id"`
	ConfigName string         `json:"config_name"`
	Content    map[string]any `json:"content"`
}
